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