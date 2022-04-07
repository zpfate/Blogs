Block相关

block本质是一个OC对象，内部也有isa指针，封装了函数调用以及函数调用环境

![img](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_16_09_53_59.jpg)

1. isa 指针，所有对象都有该指针，用于实现对象相关的功能。
2. flags，用于按 bit 位表示一些 block 的附加信息，本文后面介绍 block copy 的实现代码可以看到对该变量的使用。
3. reserved，保留变量。
4. invoke，函数指针，指向具体的 block 实现的函数调用地址。
5. descriptor， 表示该 block 的附加描述信息，主要是 size 大小，以及 copy 和 dispose 函数的指针。
6. variables，capture捕获 过来的变量，block 能够访问它外部的局部变量，就是因为将这些变量（或变量的地址）复制到了结构体中。

### 变量捕获

![image-20220302170846862](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_02_17_08_47.png)

<font color=red>*auto 默认的变量申明关键字，与static相对*</font>



![image-20220303092726423](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_03_09_27_26.png

![image-20220307135319690](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_07_13_53_20.png)

### 三种block

1. NSGlobalBlock 全局的静态 block，不会访问任何auto变量。
2. NSStackBlock 保存在栈中的 block，访问了auto变量，当函数返回时会被销毁。
3. NSMallocBlock 保存在堆中的 block，NSStackBlock调用了copy变成NSMallocBlock， 当引用计数为 0 时会被销毁。

![image-20220303101903412](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_03_10_19_03.png)

### __block

* \_\_block可以用于解决block内部无法修改auto变量值问题__
* block不能修饰全局变量、静态变量（static）
* 编译器会将__block变量包装成一个对象

### 修饰词用copy

没有copy操作，就不会在堆上，无法控制生命周期。

### block内存管理

* 当block在栈上时，不会对__block变量， 对象类型的auto变量产生强引用
* 当block被copy到堆时 

   * 会调用到block内部的copy函数
   * copy函数会调用_Block_object_assign函数
   * _Block_object_assign函数会对__block变量形成强引用(retain）
* 当block从堆中移除时

  * 会调用block内部的dispose函数
  * dispose函数内部会调用_Block_object_dispose函数
  * _Block_object_dispose函数会自动释放引用的__block变量(release)


##### __block修饰对象

* 当__block变量在栈上时，不会对指向的对象产生强引用

* 当__block变量被copy到堆时

  * 会调用__block内部的copy函数

   * copy函数内部会调用_Block_object_assign函数

   * _Block_object_assign函数会根据所指向对象的修饰符（\_\_Strong、\_\_weak、__unsafe_unretained）做出相应操作，形成强引用（retain)或者弱引用 

     **注意 ** 这里仅限于ARC时会retain，MRC时不会retain

* 当__block变量从堆上移除时

   * 会调用__block内部的dispose函数
   * dispose函数内部会调用_Block_object_dipose函数
   * _Block_object_dipose函数会自动释放指向的对象（release）

#### 解决循环引用

##### MRC

* 用__unsafe_unretained
* 用__block解决

##### ARC

* __unsafe_unretained (不会产生强引用，不安全，指向的对象销毁时，指针存储的地址值不变)

  ```objective-c
  __unsafe_unretained typeof(self) weakSelf = self;
  ```

* __weak (不会产生强引用，指向的对象销毁时会自动将指针置为nil)

  ```objective-c
  __weak typeof(self) weakSelf = self;
  ```

* __block （必须调用block，并将\_\_block修饰的变量置为nil）

  ```objective-c
    __block Person *p = [[Person alloc] init];
    p.age = 20;
    p.block = ^{
       NSLog(@"age is %zd", p.age);
       p = nil;
    };
    p.block();
  ```

  

  ![image-20220316165156112](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1647850786.png)



## KVO相关

Key-Value Observing 键值监听

### 实现原理

![image-20220314222151383](https://github.com/zpfate/uPic/2022%2003%20141647269033.png)

在程序运行中，动态创建一个NSKVONotifiying_XXX类



## 内存管理

内存管理技术：GC垃圾回收， 引用计数法

引用计数法可以及时的回收引用计数为0的对象，减少查找次数。但是引用计数会带来循环引用的问题。



## OC对象

### NSObject



### class、metaclass



### isa、superclass



![image-20220222100823130](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_22_10_08_24.png)

## Category和Extension

category的底层结构是_category_t的结构体，所含的方法在运行中通过runtime动态的将方法合并到类对象和元类对象中

![image-20220223111528215](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_23_11_15_29.png)

### 区别

1. extension扩展就是类的一部分，在编译期和头文件声明以及类实现一起形成一个完整的类，一般用来隐藏类的私有信息，无法为系统类添加extension
2. category类目则不一样，它是在运行期决定的，无法添加实例变量。
3. category不能直接添加成员变量，可以通过runtime关联对象间接添加

### 关联对象原理

实现关联对象技术的核心对象：

* AssociationsManager

* AssociationHashMap

* ObjectAssociationMap

* ObjcAssociation

  ![image-20220224154445087](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_24_15_44_46.png)

  ![image-20220224154855636](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_24_15_48_55.png)

## load和initialize

### 区别

#### 调用方式

1. load方式是根据函数地址直接调用
2. initialize是通过objc_msgSend调用

#### 调用时刻

1. load是runtime加载类、分类的时候调用（只会调用一次）
2. initialize是类第一次收到消息的时候调用，每一个类只会初始化一次（父类的initialize方法可能被调用多次）

#### 调用顺序

##### load

1. 先调用类的load，先编译的类先调用
2. 调用子类的load之前，会先调用父类的load
3. 再调用分类的load，先编译的分类，先调用

##### initialize

1. 先初始化父类
2. 再初始化子类（子类没有则调用父类的initialize方法）



## Runtime

Objective-C是一门动态性编程语言，动态性由Runtime API支撑

Runtime是一套c语言的api，封装了很多动态性相关的函数

平时编写的Objc代码，底层都是转成了Runtime API进行调用

### isa

在arm64位架构之前，isa就是一个普通的指针，存储着Class、Meta-Class对象的内存地址

从arm64架构开始，对isa进行了优化，变成了一个共用体（union）结构，还使用位域存储更多信息

![image-20220321154521223](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1647848721.png)

### Class的结构

![image-20220324093942864](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648085983.png)

### objc_msgSend

* OC中的方法调用，其实都是转换为objc_msgSend函数的调用
* objc_msgSend的执行流程可分为三大阶段
  * 消息发送
  * 动态方法解析
  * 消息转发

![image-20220324110152676](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648090913.png)

#### 动态解析

![image-20220324140711722](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648102031.png "动态解析流程")

```objc
// 动态方法解析
+ (BOOL)resolveInstanceMethod:(SEL)sel {
    if (sel == @selector(test)) {
        Method method = class_getInstanceMethod(self, @selector(instanceTest));
        class_addMethod(
                        self,
                        sel,
                        method_getImplementation(method),
                        method_getTypeEncoding(method)
                        );
        return YES;
       
    }
    return [super resolveInstanceMethod:sel];
}

+ (BOOL)resolveClassMethod:(SEL)sel {
    
    if (sel == @selector(test)) {
        Method method = class_getClassMethod(self, @selector(classTest));
        class_addMethod(
                        object_getClass(self),
                        sel,
                        method_getImplementation(method),
                        method_getTypeEncoding(method)
                        );
        return YES;
    }
    return [super resolveClassMethod:sel];
}
```

![image-20220325101031211](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648174231.png "Objective-C type encodings")



#### 消息转发

NSInvocation封装了一个方法调用，包括了 方法调用者、方法、方法参数

![image-20220324172139395](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648113699.png)

```objective-c
// 消息转发
- (id)forwardingTargetForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
//        return [[Cat alloc] init];
        return nil;
    }
    return [super forwardingTargetForSelector: aSelector];
}

// 类方法的消息转发
+ (id)forwardingTargetForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
//        return [Cat class];
    }
    return [super forwardingTargetForSelector:aSelector];
}
// 方法签名: 返回值类型 参数类型
- (NSMethodSignature *)methodSignatureForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
        Method method = class_getInstanceMethod(object_getClass(self), @selector(instanceTest));
        return [NSMethodSignature signatureWithObjCTypes:method_getTypeEncoding(method)];
        // 也可以这么生成方法签名
        return [[[Cat alloc] init] methodSignatureForSelector:aSelector];
    }
    return [super methodSignatureForSelector:aSelector];
}


+ (NSMethodSignature *)methodSignatureForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
        Method method = class_getClassMethod([Cat class], @selector(test));
        return [NSMethodSignature signatureWithObjCTypes:method_getTypeEncoding(method)];
        // 也可以这么生成方法签名
//        return [[[Cat alloc] init] methodSignatureForSelector:aSelector];
    }
    return [super methodSignatureForSelector:aSelector];
}
// NSInvocation封装了一个方法调用
// anInvocation.target 方法调用者
// anInvocation.selector 方法名
// [anInvocation getArgument:NULL atIndex:0] 方法参数 参数顺序receiver,selector,other
// [anInvocation getReturnValue:&value]; 获取返回值

- (void)forwardInvocation:(NSInvocation *)anInvocation {
//    anInvocation.target = [[Cat alloc] init];
//    [anInvocation invoke];
    [anInvocation invokeWithTarget:[[Cat alloc] init]];
}

+ (void)forwardInvocation:(NSInvocation *)anInvocation {

    [anInvocation invokeWithTarget:[Cat class]];
}
```

#### super

![image-20220325151039767](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648192240.png)

[super message]底层实现是消息发送的时候，从父类开始寻找方法实现，消息接收者仍然是子类对象。

super底层调用的是objc_msgSendSuper函数

![image-20220326213958881](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648302138.png "super底层调用的是objc_msgSendSuper")

该函数需要传入一个objc_super的结构体

![image-20220326214837125](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648302519.png "objc_super")

最后还是需要看class与superClass方法的内部实现:

![image-20220326220001204](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648303204.png "class和superClass的实现")

#### isKindOfClass、isMemberOfClass

两个方法的源码实现如下：

![image-20220326220222191](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648303343.png)

可以明显的看出，类方法（+方法）会获取元类对象，所以后面比较的也应该是元类对象。

![image-20220328100112973](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648432873.png)

### runtime应用

* 通过runtime找出控件私有成员变量 kvc修改属性
* 字典转模型
* 替换方法， method_swizzling方法交换（hook系统原生实现）
* 自动归档解档
* 关联对象添加属性
* 利用消息转发机制解决方法找不到异常管理

<font color = red> hook类簇:</font>

![image-20220329093504754](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648517705.png)

## RunLoop

运行循环，保持程序持续运行，在程序运行过程中循环做一些事情，处理App中各种事件（触摸、定时器），节省CPU资源，提高程序性能

### API

* Foundation ： NSRunLoop
* Core Foundation ： CFRunLoopRef

```objc
  // 当前线程的Runloop
    NSRunLoop *currentRunloop = [NSRunLoop currentRunLoop];
    CFRunLoopRef currentRunloopRef = CFRunLoopGetCurrent();

    // 主线程的Runloop
    NSRunLoop *mainRunloop = [NSRunLoop currentRunLoop];
    CFRunLoopRef mainloopRef = CFRunLoopGetMain();
```

![image-20220329135324877](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648533205.png "两者关联")

### RunLoop与线程

* 每一个线程都有唯一的一个与之对应的RunLoop对象
* RunLoop保存在一个全局的Dictionary里面，线程作为key，RunLoop作为value
* 线程刚创建的时候没有RunLoop，在第一次获取的时候创建
* RunLoop在线程结束时销毁
* 主线程的RunLoop已经自动获取创建，子线程默认没有开启RunLoop

### 应用范畴

* 定时器（Timer)、PerformSelector
* GCD Async Main Queue
* 事件响应、手势识别、界面刷新
* AutoreleasePool

### RunLoop相关（Core Foundation）

![image-20220329105042930](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648522243.png)

#### CFRunLoopModeRef

* CFRunLoopModeRef代表RunLoop的运行模式
* 一个RunLoop包含若干个Mode，每个Mode又包含若干Source0/Source1/Timer/Observer
* RunLoop启动只能选择其中一个Mode，作为currentMode
* 如果需要切换Mode，只能退出当前Loop，再重新选择一个Mode进入
  - 不同组的Source0/Source1/Timer/Observer能分隔开来，互不影响
* 如果Mode里没有任何Source0/Source1/Timer/Observer，RunLoop会立马退出
* 常见的两种Mode
  * KCFRunLoopDefaultMode（NSDefaultRunLoopMode）：App的默认Mode，通常主线程在该Mode下运行
  * UITrackingRunLoopMode:界面跟踪Mode，用于ScrollView追踪触摸滑动，保证页面滑动时不受其他Mode影响

#### RunLoop处理事件

* Source0
  * 触摸事件处理
  * performSelector:onThread:
* Source1 
  * 基于Port的线程间通信
  * 系统事件捕捉（Source1捕捉包装成Source1处理）
* Timers
  * NSTimer
  * performSelector: withObject: afterDelay:
* Observers
  * 用于监听RunLoop的状态
  * UI刷新（BerforeWaiting）
  * AutoRelease Pool（BerforeWaiting）

#### RunLoop运行逻辑

1. 通知Observers：进入Loop

2. 通知Observers：即将处理Timers

3. 通知Observers：即将处理Sources

4. 处理Blocks

5. 处理Source0（可能再次处理Blocks）

6. 如果存在Source1，就跳转到8-3

7. 通知Observers：开始休眠（等待消息唤醒）

8. 通知Observers：结束休眠（被某个消息唤醒）

   > 1. 处理Timer
   > 2. 处理GCD Async To Main Queue
   > 3. 处理Source1

9. 处理Blocks

10. 根据前面的执行结果，决定如何操作

    > 1. 回到 2
    > 2. 退出RunLoop

11. 通知Observers：退出Loop

#### RunLoop休眠的实现原理

![image-20220330163219333](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648629140.png)

#### RunLoop应用

* 控制线程生命周期（线程保活）
* 解决NSTimer在滑动时停止工作
* 监控应用卡顿
* 性能优化

## 多线程

### 多线程方案

| 技术方案    | 简介                                                         | 语言 | 线程生命周期 | 使用频率 |
| ----------- | ------------------------------------------------------------ | ---- | ------------ | -------- |
| pthread     | 1. 一套通用的多线程API<br />2. 适用于Unix/Linux/Windows等系统<br />3.  跨平台/可移植 | C    | 程序员管理   | 几乎不用 |
| NSThread    | 1. 使用更加面向对象<br />2. 简单易用，可直接操作线程对象     | OC   | 程序员管理   | 偶尔使用 |
| GCD         | 1. 旨在替代NSThread等线程技术<br />2. 充分利用设备的多核     | C    | 自动管理     | 经常使用 |
| NSOperation | 1. 基于GCD（底层是GCD）<br />2. 比GCD多了一些简单实用的功能<br />3. 使用更加面向对象 | OC   | 自动管理     | 经常使用 |

### GCD

#### GCD的队列

- <font color = red>并发</font>队列（Concurrent Dispatch Queue）

  > * 可以让多个任务<font color = blue>并发</font>（<font color = blue>同时</font>）执行（自动开启多个线程同时执行任务）
  > * 并发功能只有在异步（dispatch_async）函数下有效

- <font color = red>串行</font>队列（Serial Dispatch Queue）

  >* 让任务一个接着一个地执行
  >
  >

<font color = blue>同步</font>： 在当前线程中执行任务，不具备开启线程的能力

<font color = blue>异步</font>：在新的线程中执行任务，具备开启线程的能力（主队列不开启）

<font color = red>并发</font>：多个任务并发（同时）执行

<font color = red>串行</font>：一个任务执行完毕，再执行下一个任务

**同步和异步主要影响能不能开启新的线程，并发和串行主要影响执行任务的方式**

| --            | 并发队列                                                     | 手动创建的串行队列                                           | 主队列                                                       |
| ------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 同步（sync）  | <font color = yellow>没有</font>开启新线程<br /><font color = green>串行</font>执行任务 | <font color = yellow>没有</font>开启新线程<br /><font color = green>串行</font>执行任务 | <font color = yellow>没有</font>开启新线程<br /><font color = green>串行</font>执行任务 |
| 异步（async） | <font color = blue>有</font>开启新线程<br /><font color = red>并发</font>执行任务 | <font color = blue>有</font>开启新线程<br /><font color = green>串行</font>执行任务 | <font color = yellow>没有</font>开启新线程<br /><font color = green>串行</font>执行任务 |

<font color = red>**使用sync函数往当前串行队列中添加任务，会卡住当前的串行**

**队列（产生死锁）**</font>

![image-20220401105327554](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648781608.png)

#### dispatch_group

```objective-c
    dispatch_group_t group = dispatch_group_create();
    // 创建并发队列
    dispatch_queue_t queue = dispatch_queue_create("group", DISPATCH_QUEUE_CONCURRENT);
    
    dispatch_group_async(group, queue, ^{
        for (int i = 0; i < 5; i++) {
            NSLog(@"任务1---%@", [NSThread currentThread]);
        }
    });
    dispatch_group_async(group, queue, ^{
        for (int i = 0; i < 5; i++) {
            NSLog(@"任务2---%@", [NSThread currentThread]);
        }
    });
    
    // 等前面的任务执行完毕后,自动执行
//    dispatch_group_notify(group, queue, ^{
//        dispatch_async(dispatch_get_main_queue(), ^{
//            for (int i = 0; i < 5; i++) {
//                NSLog(@"任务2---%@", [NSThread currentThread]);
//            }
//        });
//    });
    
    dispatch_group_notify(group, dispatch_get_main_queue(), ^{
            for (int i = 0; i < 5; i++) {
                NSLog(@"任务3---%@", [NSThread currentThread]);
            }
    });
```

### iOS线程同步方案

* OSSpinLock
* os_unfair_lock
* Pthread_mutex
* dispatch_semaphore
* dispatch_quequ(DISPATCH_QUEUE_SERIAL)
* NSLock
* NSCondition
* NSConditionLock
* @synchronized

#### OSSpinLock 自旋锁

等待锁的线程会出于忙等（busy-wait)状态，一直占用着CPU资源。

目前已经不安全，不建议使用，会造成优先级反转

<font color = red>**OSSpinLock已经不建议使用，iOS10.0之后建议使用os_unfair_lock**</font>

> 如果等待锁的线程优先级比较高，它会一直占用CPU资源，优先级低的线程就无法释放锁

##### 使用

需要\#import <libkern/OSAtomic.h>

```objc
 // 初始化
 OSSpinLock _lock = OS_SPINLOCK_INIT;
 // 尝试加锁(如果需要等待就不加锁,直接返回false;如果不需要等待就加锁,返回true)
 bool result = OSSpinLockTry(&_lock);
 // 加锁
 OSSpinLockLock(&_lock);
 // 解锁
 OSSpinLockUnlock(&_lock);
```

#### os_unfair_lock 互斥锁

* os_unfair_lock用于取代不安全的OSSpinLock，从iOS10开始才支持

* 从底层调用看，等待os_unfair_lock锁的线程处于休眠状态，并非忙等

##### 使用

需要\#import <os/lock.h>

```objective-c
 // 初始化
  os_unfair_lock lock = OS_UNFAIR_LOCK_INIT;
  // 尝试加锁
  bool result = os_unfair_lock_trylock(&lock);
  // 加锁
  os_unfair_lock_lock(&lock);
  // 解锁
  os_unfair_lock_unlock(&lock);
```

#### pthread_mutex（跨平台）

* mutex叫做"互斥锁"，等待锁的线程会处于休眠状态

##### 使用

需要#import <pthread.h>

![image-20220402101501432](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648865702.png)

```objective-c
    // 静态初始化
//        pthread_mutex_t _mutex = PTHREAD_MUTEX_INITIALIZER;
    /**
     #define PTHREAD_MUTEX_NORMAL        0
     #define PTHREAD_MUTEX_ERRORCHECK    1
     #define PTHREAD_MUTEX_RECURSIVE        2 递归锁
     #define PTHREAD_MUTEX_DEFAULT        PTHREAD_MUTEX_NORMAL 普通的锁
     */
    
    // 初始化锁的属性
    pthread_mutexattr_t attr;
    pthread_mutexattr_init(&attr);
    // 锁的类型
    pthread_mutexattr_settype(&attr, PTHREAD_MUTEX_DEFAULT);
    // 初始化锁
    pthread_mutex_t _mutex;
    pthread_mutex_init(&_mutex, &attr);
    // 销毁属性
    pthread_mutexattr_destroy(&attr);
   // 加锁
    pthread_mutex_lock(&_mutex);
    // 解锁
    pthread_mutex_unlock(&_mutex);
    // 销毁锁
    pthread_mutex_destroy(&_mutex);
```

**递归锁** 允许<font color = red>同一个线程</font>对一把锁重复加锁

#### NSLock、NSRecursiveLock

* NSLock就是对mutex普通锁的封装
* NSRecursiveLock是对mutex递归锁的封装

```objective-c
  // 初始化
    NSLock *lock = [[NSLock alloc] init];
    // 尝试加锁
    [lock tryLock];
    // 加锁
    [lock lock];
    // 解锁
    [lock unlock];
    // 在传入时间到来前加锁,时间没到就休眠,时间到了加锁,
    // 加锁成功就返回YES,加锁失败或超出时间就返回NO
    lock lockBeforeDate:<#(nonnull NSDate *)#>
```

#### NSCondition

* NSConditionLock是对mutex和cond条件锁的封装

```objective-c
    // 初始化
    NSCondition *condition = [[NSCondition alloc] init];
    // 加锁
    [condition lock];
    // 等待
    [condition wait];
    // 通知
    [condition signal];
    // 广播
    [condition broadcast];
    // 解锁
    [condition unlock];
```



#### NSConditionLock

* 是对NSCondition的进一步封装
* 可以设置条件的值

```objective-c
   // 初始化 默认condition为0
    NSConditionLock *conditionLock = [[NSConditionLock alloc] initWithCondition:1];
    // 加锁 直接加锁 不管condition的值
    [conditionLock lock];
    // 符合条件加锁
    [conditionLock lockWhenCondition:1];
    // 解锁
    [conditionLock unlock];
    // 符合条件解锁
    [conditionLock unlockWithCondition:1];
```

#### dispatch_queue

* 直接使用GCD的串行队列，也是可以实现线程同步的

```objective-c
// 创建队列
dispatch_quue_t queue = dispatch_queue_create("ticket", 
DISPATCH_QUEUE_SERIAL);
// 同步执行
dispatch_sync(queue, ^{
    [super _saveMoney];
});
```

#### dispatch_semaphore

* semaphore叫做信号量
* 信号量的初始值 可以用来控制线程并发访问的最大数量

```objective-c
// 信号量初始值
int value = 1;
// 初始化信号量
dispatch_semaphore_t semaphore = dispatch_semaphore_create(value);
// 如果信号量的值<=0,当前线程就会进入休眠等待(直到信号量的值>0)
// 如果信号的值>0,就减一,然后往下执行代码
dispatch_semaphore_wait(semaphore, DISPATCH_TIME_FOREVER);
// 让信号量的值-1
dispatch_semaphore_signal(semaphore);
```

#### @synchronized

* <font color = red>@synchronized</font>是对mutex递归锁的封装

* 源码查看objc-sync.mm

### iOS线程同步方案性能比较

1. os_unfair_lock
2. OSSpinLock
3. dispatch_semaphore
4. pthread_mutex
5. dispatch_queue
6. NSLock
7. NSCondition
8. pthread_mutex（recursive）
9. NSRecursiveLock
10. NSConditionLock
11. @synchronized

#### 自旋锁、互斥锁比较

什么情况下使用自旋锁比较划算？

* 预计线程等待锁的时间很短
* 加锁的代码（临界区）经常被调用，但竞争情况很少发生
* CPU资源不紧张
* 多核处理器

什么情况使用互斥锁比较划算？

* 预计线程等待锁的时间比较长
* 单核处理器
* 临界区有IO操作
* 代码比较复杂，或者循环量大
* 临界区竞争非常激烈

### atomic

* <font color = red>atomic</font>用于保证属性setter、getter的原子性操作，相当于在getter和setter内部加了线程同步的锁
* 可以参考源码objc4的objc-accessors.m
* 它并不能保证使用属性的过程是线程安全的（比如集合的操作）

### 读写安全

* 同一时间单写多读
* 同一时间不允许既有写的操作，又有读的操作

#### 实现方案

* pthread_rwlock: 读写锁
* dispatch_barrier_async:异步栅栏调用

#### pthread_rwlock

使用需要#import <pthread.h>

```objective-c
- (void)initLock {
  // 初始化
  pthread_rwlock_init(&_lock, NULL);
}

- (void)read {
    // 加锁
    pthread_rwlock_rdlock(&_lock);

    NSLog(@"%s", __func__);
    pthread_rwlock_unlock(&_lock);
}

- (void)write {
    // 加锁
    pthread_rwlock_wrlock(&_lock);
    NSLog(@"%s", __func__);
    // 解锁
    pthread_rwlock_unlock(&_lock);
}
```

#### dispatch_barrier_async

* 这个函数传入的并发队列必须是自己通过dispatch_queue_create创建的
* 如果传入的是一个串行或是一个全局并发队列，那这个函数变等同于dispatch_async的效果

```objective-c
    dispatch_queue_t queue = dispatch_queue_create("rw", DISPATCH_QUEUE_CONCURRENT);
    dispatch_async(queue, ^{
       // 读
        [self read];
    });
    dispatch_barrier_async(queue, ^{
        [self write];  // 写
    });
```

## 内存管理





## 启动优化

### 启动过程

#### 1. pre-main阶段

* 加载应用的可执行文件
* 加载动态连接器dyld