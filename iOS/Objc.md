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

#### RunLoop运行的逻辑

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

## 启动优化

### 启动过程

#### 1. pre-main阶段

* 加载应用的可执行文件
* 加载动态连接器dyld