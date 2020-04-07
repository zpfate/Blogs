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

