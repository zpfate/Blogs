

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

### 加载

1. 通过Runtime加载某个类所有的Category
2. 把所有的Category的方法、属性、协议数据，合并到一个大数组中（后面参与编译的Category数据，会在数组前部）
3. 将合并后的分类数据（方法、属性、协议），插入到类原来数据的前面

### 和Extension区别

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
  * UITrackingRunLoopMode：界面跟踪Mode，用于ScrollView追踪触摸滑动，保证页面滑动时不受其他Mode影响

#### RunLoop处理事件

* Source0
  * 触摸事件处理
  * performSelector:onThread:
* Source1 
  * 基于Port的线程间通信
  * 系统事件捕捉（Source1捕捉包装成Source0处理）
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

### CADisplayLink

```objective-c
self.link = [CADisplayLink displayLinkWithTarget:[TFProxy proxyWithTarget:self] selector:@selector(linkTest)];
    [self.link addToRunLoop:[NSRunLoop mainRunLoop] forMode:NSDefaultRunLoopMode];
```

### NSTimer

NSTimer依赖于RunLoop，如果RunLoop任务过于繁重，定时器可能不准时

```objective-c
self.timer = [NSTimer scheduledTimerWithTimeInterval:1.0 target:[TFProxy proxyWithTarget:self] selector:@selector(timerTest) userInfo:nil repeats:YES];
```

### NSProxy

* 继承自NSProxy的类大部分方法都会走消息转发
* 消息转发的效率更高

头文件

```objective-c
@interface TFProxy : NSProxy
@property (nonatomic, weak) id target;

+ (instancetype)proxyWithTarget:(id)target;

@end
```

实现

```objective-c
@implementation TFProxy

+ (instancetype)proxyWithTarget:(id)target {
    TFProxy *proxy = [TFProxy alloc];
    proxy.target = target;
    return proxy;
}

- (id)forwardingTargetForSelector:(SEL)aSelector {
    return self.target;
}

@end
```

### GCD定时器

GCD定时器更加准确，依赖于系统内核

```objc
   // 队列
    dispatch_queue_t queue = dispatch_get_main_queue();
    // 创建GCD定时器
    dispatch_source_t timer = dispatch_source_create(DISPATCH_SOURCE_TYPE_TIMER, 0, 0, queue);
    
    // 设置时间
    NSTimeInterval start = 2.0; // 2s后开始执行
    NSTimeInterval interval = 1.0; // 每隔1s执行
    // 纳秒
    dispatch_source_set_timer(timer, dispatch_time(DISPATCH_TIME_NOW, start * NSEC_PER_SEC), interval * NSEC_PER_SEC, 0);
    
    // 设置回调
    dispatch_source_set_event_handler(timer, ^{
        NSLog(@"-----");
    });
		// 销毁
    dispatch_resume(timer);
```

## 内存管理

### iOS内存布局

![image-20220408101744432](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649384265.png)

### Tagged Pointer

* 从64bit开始，iOS引入了Tagged Pointer技术，用于优化NSNumber，NSDate，NSString等小对象的存储
* 在没有使用Tagged Pointer之前，NSNumber等对象需要动态分配内存，维护引用计数等，NSNumber指针存储的是堆中NSNumber对象的地址值。
* 使用Tagged Pointer之后，NSNumber指针里面存储的数据变成了：Tag + Data，也就是将数据直接存储在了指针中
* 当指针不够存储数据时，才会使用动态分配内存的方式来存储数据
* objc_msgSend能识别到Tagged Pointer，比如NSNumber的intValue方法，直接从指针提取数据，节省了以前的调用开销。
* 在mac平台最低有效位是1就是Tagged Pointer，iOS平台为1UL<<63，最高有效位是1

### 内存管理

* iOS中使用引用计数来管理OC对象的内存，引用计数储存在isa指针中或者SideTable中
* 一个新建的OC对象引用计数默认是1，当引用计数减为0，OC对象就会销毁，释放其占用的内存空间
* 调用retain会让OC对象的引用计数+1，调用release会让OC对象的引用计数-1

### copy和mutableCopy

1. copy 不可变拷贝， 产生不可变副本
2. mutableCopy可变拷贝，产生可变副本

|             | NSString                                             | NSMutableString                                      | NSArray                                             | NSMutableArray                                      | NSDictionary                                             | NSMutableDictionary                                      |
| ----------- | ---------------------------------------------------- | ---------------------------------------------------- | --------------------------------------------------- | --------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| copy        | NSString<br /><font color = green>浅拷贝</font>      | NSString<br /><font color = red>深拷贝</font>        | NSArray<br /><font color = green>浅拷贝</font>      | NSArray<br /><font color = red>深拷贝</font>        | NSDictionary<br /><font color = green>浅拷贝</font>      | NSDictionary<br /><font color = red>深拷贝</font>        |
| mutableCopy | NSMutableString<br /><font color = red>深拷贝</font> | NSMutableString<br /><font color = red>深拷贝</font> | NSMutableArray<br /><font color = red>深拷贝</font> | NSMutableArray<br /><font color = red>深拷贝</font> | NSMutableDictionary<br /><font color = red>深拷贝</font> | NSMutableDictionary<br /><font color = red>深拷贝</font> |

### weak

weak是弱引用，用weak来修饰的引用对象的计数器不会增加，而且weak会在引用对象释放的时候，自动置为nil

#### 当 weak 指向的对象被释放时，如何让 weak 指针置为 nil 的呢

* 调用 objc_release
* 因为对象的引用计数为0，所以执行 dealloc
* 在 dealloc 中，调用了 _objc_rootDealloc 函数
* 在 _objc_rootDealloc 中，调用了 object_dispose 函数
* 调用 objc_destructInstance
* 最后调用 objc_clear_deallocating,详细过程如下：
  a. 从 weak 表中获取废弃对象的地址为键值的记录
  b. 将包含在记录中的所有附有 weak 修饰符变量的地址，赋值为 nil
  c. 将 weak 表中该记录删除
  d. 从引用计数表中删除废弃对象的地址为键值的记录



## 性能优化

### CPU和GPU

![image-20220413153848554](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841056.png)

### 卡顿优化-CPU

![image-20220413154742051](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841048.png)

### 卡顿优化-GPU

![image-20220413160647442](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841046.png)

### 离屏渲染

![image-20220413161034602](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841042.png)

### 卡顿检测

* 平时所说的卡顿 主要是因为主线程执行了比较耗时的操作
* 可以添加Observer到主线程RunLoop中，通过监听RunLoop状态切换的耗时，以达到监控卡顿的目的

### 耗电优化

![image-20220413161634558](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841040.png)

![image-20220413164049114](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841037.png)

![image-20220413164351339](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841032.png)

## 启动优化

### 启动方式

* 冷启动：从零开始启动APP
* 热启动：APP已经存在内存中，在后台存活着，再次启动APP

### APP的启动

APP启动可以分成三大阶段

* dyld
* runtime
* main

#### dyld

* Apple的动态链接器，可以用来装载Mach-O文件（可执行文件，动态库等）
* 启动APP时，dyld所做的事情有
  * 装载APP的可执行文件，同时会递归加载所有依赖的动态库
  * 当dyld把可执行文件、动态库都装载完毕后，会通知Runtime进行下一步处理

#### Runtime

启动APP时，Runtime所做的事情：

> * 使用map_images进行可执行文件内容的解析和处理
> * 在load_images中调用call_load_methods，调用所有Class和Category的+load方法
> * 进行各种objc结构的初始化（注册Objc类、初始化类对象等等）
> * 调用c++静态初始化器和\_\_attribute_((constructor))修饰的函数

到此为止，可执行文件和动态库中的所有符号（Class、Protocol、Selector、IMP、...）都已经按格式成功加载到内存中，被runtime所管理

#### main

总结一下：

>* APP的启动由dyld主导，将可执行文件加载到内存，顺便加载所有的依赖库
>* 并由runtime负责加载成objc定义的结构
>* 所有初始化工作完成结束后，dyld就会调用main函数
>* 接下来就是UIApplication函数，didFinishLaunchingWithOptions:

### 启动优化

![image-20220413170106569](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841021.png)

## 应用瘦身

* 资源

  对资源进行无损的压缩，比如使用webp格式图片

  去除没有用到的资源： https://github.com/tinymind/LSUnusedResources 

* 可执行文件

  编译器优化：Strip Linked Product、Make Strings Read-Only、Symbols Hidden by Default 设置为 YES，去掉异常支持，Enable C++ Exceptions、Enable Objective-C Exceptions 设置为 NO， Other C Flags 添加 -fno-exceptions 

  利用 AppCode 检测未使用的代码：菜单栏 -> Code -> Inspect Code

![image-20220413170702762](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841015.png)

![image-20220413170916344](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1649841004.png)

## 设计模式

### MVC

* 优点是View、Model可以复用
* Controller随着版本的迭代过于臃肿

## 推荐

严蔚敏《数据结构》

《大话数据结构与算法》

《HTTP权威指南》

《TCP/IP详解卷1：协议》