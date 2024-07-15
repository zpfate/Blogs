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

* 每一个线程都有唯一的一个与之对应的`RunLoop`对象
* `RunLoop`保存在一个全局的`Dictionary`里面，线程作为`key`，`RunLoop`作为`value`
* 线程刚创建的时候没有`RunLoop`，在第一次获取的时候创建
* `RunLoop`在线程结束时销毁
* 主线程的`RunLoop`已经自动获取创建，子线程默认没有开启`RunLoop`

### 应用范畴

* 定时器（Timer)、PerformSelector
* GCD Async Main Queue
* 事件响应、手势识别、界面刷新
* AutoreleasePool

### RunLoop相关（Core Foundation）

![image-20220329105042930](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648522243.png)

#### CFRunLoopModeRef

* `CFRunLoopModeRef`代表`RunLoop`的运行模式
* 一个`RunLoop`包含若干个`Mode`，每个`Mode`又包含若干`Source0/Source1/Timer/Observer`
* `RunLoop`启动只能选择其中一个`Mode`，作为`currentMode`
* 如果需要切换`Mode`，只能退出当前`Loop`，再重新选择一个Mode进入
  - 不同组的`Source0/Source1/Timer/Observer`能分隔开来，互不影响
* 如果`Mode`里没有任何`Source0/Source1/Timer/Observer`，RunLoop会立马退出
* 常见的两种`Mode`
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

![RunLoop_1](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1658471907.png)

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