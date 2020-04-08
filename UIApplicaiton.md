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

  ```
  [UIApplication sharedApplication].delegate.window
  ```

* keyWindow

  ```
  [UIApplication sharedApplication].keyWindow
  ```

  ***注意：***

  正常情况下我们获取keyWindow与delegateWindow,会发现两者一致，如下：

  ![keyWindow与delegateWindow](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/image-2020040720405172020200407204051.png)

  

  事实上，当你在Appdelegate中的window调用`makeKeyAndVisible`方法之前打印keyWindow，会发现为`null`，一个window只需要调用`makeKeyAndVisible`方法就可以成为keyWindow。已经废弃的UIAlertView在展示的时候keyWindow也会发生变化，如下所示：

  ![AlertView的keyWindow](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/image-2020040809112985320200408093931.png)

  所以keyWindow有可能是会改变的，但是同一时间只会有一个UIWindow是keyWindow

  UIWindow是根据UIWindowLevel来进行排序的，level高的将排在所有比他低的层级前面，keyWindow显示在相同级别的最上层。

  



