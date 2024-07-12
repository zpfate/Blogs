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

<font color = red>**使用sync函数往当前串行队列中添加任务，会卡住当前的串行**</font>

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

## 

[Demo](https://github.com/zpfate/Lesson-Objc/tree/main/Multi-ThreadDemo)