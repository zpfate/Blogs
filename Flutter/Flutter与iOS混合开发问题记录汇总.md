## Q1：Flutter多页面开发问题

目前两个方案：flutter_boost和多引擎方案

`flutter_boost`的问题是已经不更新了，该库可以在 [Dart Packages](https://pub.dev/packages/flutter_boost)查看

![image-20230717101443786](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1689650371266)

多引擎方案自己管理，有诸多坑要填，先附上[Flutter多引擎方案官方链接](https://flutter.cn/docs/development/add-to-app/multiple-flutters)

## Q2：原生侧滑手势与Flutter页面侧滑冲突

该问题可能出现的情况是，打开`FlutterViewController`之后，侧滑直接pop掉整个`FlutterViewController`，安卓好像是物理按键也是一样的效果，物理返回键会导致直接退出App。

## Q3：Flutter端数据无法公用

这个产生的原因应该是各个`engine`之间的内存是隔离的，所以各个Flutter页面之间的代码是无法公用的。解决方案有两种：

1. 通过原生和Flutter交互的`channel`传参，可以在`ViewController`的viewWillAppear方法中调用，可以确保每次刷新

2. 本地缓存

   比如说`NSUserDefaults`和`shared_prefrences`联动

   **需要注意的是key值需要添加flutter.**

## Q4：Flutter Attach总是失败

这个问题很难妥善解决，感觉应该是多引擎方案的自带bug，attach的时候会出现同一个bundle id的多个程序，可能后面多个后缀（2）（3）



## Q5：弹框比如说Get.bottomSheet上有输入框，弹出键盘后不再移动

不管是GetX的bottomSheet或者flutter原生的showModalBottomSheet，在bottomSheet上有输入框且sheet比较高的时候，键盘弹出widget整体上移，会和状态栏重合在iOS上无法点击

解决方案：使用margin，把剩下的空间利用margin顶住，这样键盘弹起的时候`widget`不会随之上移。



