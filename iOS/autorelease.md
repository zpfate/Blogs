# autorelease

## MRC下使用

在MRC环境中使用自动释放池需要用到 `NSAutoreleasePool`对象，其生命周期就相当于C语言变量的作用域。对于所有调用过`autorelease`方法的对象，在废弃 `NSAutoreleasePool`对象时，都将调用 `release`实例方法。

```objective-c
//MRC环境下的测试：
//第一步：生成并持有释放池NSAutoreleasePool对象;
NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];

//第二步：调用对象的autorelease实例方法;
id obj = [[NSObject alloc] init];
[obj autorelease];

//第三步：废弃NSAutoreleasePool对象;
[pool drain];   //向pool管理的所有对象发送消息，相当于[obj release]

//obi已经释放，再次调用会崩溃(Thread 1: EXC_BAD_ACCESS (code=EXC_I386_GPFLT))
NSLog(@"打印obj：%@", obj);
```

## ARC下使用

ARC环境不能使用 `NSAutoreleasePool`类也不能调用 `autorelease`方法，代替它们实现对象自动释放的是 `@autoreleasepool`块和 `__autoreleasing`修饰符。

```objective-c
//ARC环境下的测试：
@autoreleasepool {
    id obj = [[NSObject alloc] init];
    NSLog(@"打印obj：%@", obj); 
}
```

为了探究释放池的底层实现，我们在终端使用 clang-rewrite-objc+文件名命令将上述OC代码转化为C++源码：

```c++
int main(int argc, const char * argv[]) {
    /* @autoreleasepool */
    {
        __AtAutoreleasePool __autoreleasepool;
        NSLog((NSString *)&__NSConstantStringImpl__var_folders_3f_crl5bnj956d806cp7d3ctqhm0000gn_T_main_d37e0d_mi_0);
     }//大括号对应释放池的作用域

     return 0;
}
```

### __AtAutoreleasePool结构体的实现代码

```c++
extern "C" __declspec(dllimport) void * objc_autoreleasePoolPush(void);
extern "C" __declspec(dllimport) void objc_autoreleasePoolPop(void *);

struct __AtAutoreleasePool {
  __AtAutoreleasePool() {atautoreleasepoolobj = objc_autoreleasePoolPush();}
  ~__AtAutoreleasePool() {objc_autoreleasePoolPop(atautoreleasepoolobj);}
  void * atautoreleasepoolobj;
};
```



* 自动释放池的主要底层数据结构是：``__AutoReleasePool`，`AutoReleasePoolPage`
* 调用了``autorelease`的对象最终都是通过`AutoReleasePoolPage`
* 每个`AutoReleasePoolPage`对象占用4096字节内存，除了用来存放他内部的成员变量，剩下的空间用来存放`autorelease`对象的地址
* 所有的`AutoReleasePoolPage`对象是通过双向链表的形式连接在一起



## AutoreleasePoolPage的结构

```c++
//大致在641行代码开始
class AutoreleasePoolPage {
#   define EMPTY_POOL_PLACEHOLDER ((id*)1)  //空池占位
#   define POOL_BOUNDARY nil                //边界对象(即哨兵对象）
    static pthread_key_t const key = AUTORELEASE_POOL_KEY;
    static uint8_t const SCRIBBLE = 0xA3;  // 0xA3A3A3A3 after releasing
    static size_t const SIZE = 
#if PROTECT_AUTORELEASEPOOL
        PAGE_MAX_SIZE;  // must be multiple of vm page size
#else
        PAGE_MAX_SIZE;  // size and alignment, power of 2
#endif
    static size_t const COUNT = SIZE / sizeof(id);
    magic_t const magic;                  //校验AutoreleasePagePoolPage结构是否完整
    id *next;                             //指向新加入的autorelease对象的下一个位置，初始化时指向begin()
    pthread_t const thread;               //当前所在线程，AutoreleasePool是和线程一一对应的
    AutoreleasePoolPage * const parent;   //指向父节点page，第一个结点的parent值为nil
    AutoreleasePoolPage *child;           //指向子节点page，最后一个结点的child值为nil
    uint32_t const depth;                 //链表深度，节点个数
    uint32_t hiwat;                       //数据容纳的一个上限
    //......
};
```



* `magic` 用来校验 `AutoreleasePoolPage` 的结构是否完整；
* `next` 指向最新添加的 `autoreleased `对象的下一个位置，初始化时指向 begin() ；
* `thread` 指向当前线程；
* `parent` 指向父结点，第一个结点的 `parent` 值为 `nil` ；
* `child` 指向子结点，最后一个结点的 `child` 值为` nil` ；
* `depth` 代表深度，从 0 开始，往后递增 1；
* `hiwat` 代表 `high water mark `。

![img](https://upload-images.jianshu.io/upload_images/5835116-170e421948bfb845.png?imageMogr2/auto-orient/strip|imageView2/2/w/836/format/webp)

* 调用`push`方法会将一个`POOL_BOUNDARY`入栈，并且会返回其存放的内存地址
* 调用`pop`放时传入一个`POOL_BOUNDARY`的内存地址，会从最后一个入栈的对象开始发送`release`消息，直到遇到这个`POOL_BOUNDARY`
* `id *next`指向了下一个能存放`autorelease`对象地址的区域

## RunLoop和AutoRelease

![ios 什么是自动释放池 ios自动释放池原理_自动释放池_05](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1719821549039)

iOS在主线程的`RunLoop`注册了2个`Observer`

* 第一个`Observer`监听了`KCFRunLoopEntry`事件，会调用`objc_autoreleasePoolPush()`

* 第二个`Observer`监听了`KCFRunLoopBeforeWaiting`事件，会调用`objc_autoreleasePoolPop()`、`objc_autoreleasePoolPush()`，监听了`KCFRunLoopBeforeExit`事件，会调用`objc_autoreleasePoolPop()`

每一个线程都会维护自己的 `Autoreleasepool`栈，所以子线程虽然默认没有开启 `RunLoop`，但是依然存在 `AutoreleasePool`，在子线程退出的时候会去释放` autorelease`对象。

`ARC`会根据一些情况进行优化，添加 `__autoreleasing`修饰符，其实这就相当于对需要延时释放的对象调用了 `autorelease`方法。从源码分析的角度来看，如果子线程中没有创建 `AutoreleasePool` ，而一旦产生了 `Autorelease`对象，就会调用 `autoreleaseNoPage`方法自动创建 `hotpage`，并将对象加入到其栈中。所以，一般情况下，子线程中即使我们不手动添加自动释放池，也不会产生内存泄漏。

### AutoreleasePool需要手动添加的情况

1. 编写的不是基于UI框架的程序，例如命令行工具；
2. 通过循环方式创建大量临时对象；
3. 使用非Cocoa程序创建的子线程；
