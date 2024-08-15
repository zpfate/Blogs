做iOS开发的同学对`UIApplication`一定不陌生，当我们想要弹框，蒙层等功能时，都会使用类似如下的操作：

```objc
    UIWindow *delegateWindow = [UIApplication sharedApplication].delegate.window;
```

```objc
    UIWindow *keyWindow = [UIApplication sharedApplication].keyWindow;
```
这里面`delegateWindow`和`keyWindow`的区别你是否了解？还有`UIApplication`是不是也深入了解过？

## UIApplication

当我们想了解一个类时，基本操作首先打开官方文档的介绍。

![UIApplication官方介绍](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/ZVbDa2jIxeyTHP620200407155049.png)

翻译一下：iOS中程序运行的集中控制和协调点。

每个iOS应用都只有一个UIApplication实例（少部分是UIApplication的子类），我们经常使用`[UIApplication sharedApplication]`来获取它，也可知这是一个单例类。

application对象主要作用是处理初始化传入的用户事件进程，它调度由控制对象 (UIControl 的实例) 转发给它的动作消息(action message) 到适当的目标对象。

同时，application 对象维护一个打开的窗口列表 (UIWindow 对象),通过这些可以检索任何应用程序的 UIView 对象

UIApplication类还定义了一个delegate（AppDelegate），该delegate遵循UIApplicationDelegate协议，必须实现的一些方法，均是application向delegate通知的重要的运行时事件，比如说程序的启动，进入活跃状态，退出等等。

下面列举一些UIApplication类常用的方法

#### 获取app实例

```Objc
[UIApplication sharedApplication]
```

#### 获取delegate（AppDelegate）

```objc
[UIApplication sharedApplication].delegate
```

#### 获取app windows

* delegateWindow

  ```objc
  [UIApplication sharedApplication].delegate.window
  ```

* keyWindow

  ```objective-c
  [UIApplication sharedApplication].keyWindow
  ```

  ***注意：***

  正常情况下我们获取keyWindow与delegateWindow,会发现两者一致，如下：

  ![keyWindow与delegateWindow](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/image-2020040720405172020200407204051.png)

  

  事实上，当你在`Appdelegate`中的`window`调用`makeKeyAndVisible`方法之前打印`keyWindow`，会发现为`null`，一个window只需要调用`makeKeyAndVisible`方法就可以成为`keyWindow`。已经废弃的`UIAlertView`在展示的时候`keyWindow`也会发生变化，如下所示：

  ![AlertView的keyWindow](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/image-2020040809112985320200408093931.png)

  所以`keyWindow`有可能是会改变的，但是同一时间只会有一个`UIWindow`是`keyWindow`

  `UIWindow`是根据`UIWindowLevel`来进行排序的，`level`高的将排在所有比他低的层级前面，`keyWindow`显示在相同级别的最上层。

  

  * windows

    ```objc
    [UIApplication sharedApplication].windows
    ```

#### 打开 URL 资源 

```objc
// iOS 10废弃
[UIApplication sharedApplication] openURL 
[UIApplication sharedApplication] openURL:options:completionHandler:
```



#### 控制和处理事件

- sendEvent: 将事件调度到 app 中的相应响应者对象
- sendAction:to:from:forEvent: : 将选择器识别的动作消息发送到指定目标
- beginIgnoringInteractionEvents : 告知接收者暂停处理触摸事件
- endIgnoringInteractionEvents : 告知接收者恢复相关事件的处理
- ignoringInteractionEvents : 一个布尔值, 指示接收者是否忽略由触摸在屏幕上启动的事件
- applicationSupportsShakeToEdit : 一个布尔值,用于确定摇动设备是否显示撤销重做用户界面

#### 注册远程通知 

* registerForRemoteNotifications : 注册通过 Apple Push Notification 服务接收远程通知.

* unregisterForRemoteNotifications : 取消注册通过 Apple Push Notification 服务接收远程通知.

* registeredForRemoteNotifications : 一个布尔值, 指示 app 当前是否注册了远程通知.

#### 管理后台执行

* applicationState :  app 的 runtime 状态

* backgroundTimeRemaining : app 在后台运行的时间量

* backgroundRefreshStatus : app 进入到后台,使其可以执行 background behaviors 的能力.

* setMinimumBackgroundFetchInterval: 指定后台提取操作之间必须经过的最短时间.

* beginBackgroundTaskWithName:expirationHandler: 标记具有指定名称的新的长时间运行的后台任务开始.

* beginBackgroundTaskWithExpirationHandler : 标记着一个新的长期运行的后台任务的开始.

* endBackgroundTask : 标记特定长时间运行的后台任务结束.



