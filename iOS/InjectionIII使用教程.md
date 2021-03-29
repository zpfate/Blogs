

## 什么是InjectionIII
InjectionIII是John Holdsworth开发的一款可以动态的将`Swift`或者`Objective-C`代码在已运行的程序中执行的工具，可以不用重新运行程序，加快调试速度。这款工具可以在`Mac AppStore`免费下载，如下所示：

![image-20200525192215540](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_05_26_15_08_36.png)

同时作者也将该软件开源了，[源代码地址](https://github.com/johnno1962/InjectionIII)

[^备注]: Xcode 10.2 and later (Swift 5+)



## 使用方法

1.从`Mac AppStore`下载完成后，打开软件，在电脑右上方任务栏会出现针孔的图标，点击`Open Project`，或者`Open Recent`，选择工程文件的根文件夹，再点击`File Watcher`至选中状态。

![image-20200525192810711](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_05_26_15_09_01.png)

2.打开Xcode工程，在`Appdelegate.m`文件`didFinishLaunchingWithOptions`方法里添加如下代码

```objc
#if DEBUG
    // iOS
    [[NSBundle bundleWithPath:@"/Applications/InjectionIII.app/Contents/Resources/iOSInjection.bundle"] load];
    // tvOS
    //[[NSBundle bundleWithPath:@"/Applications/InjectionIII.app/Contents/Resources/tvOSInjection.bundle"] load];
    // macOS
    //[[NSBundle bundleWithPath:@"/Applications/InjectionIII.app/Contents/Resources/macOSInjection.bundle"] load];
#endif
```

```swift
#if DEBUG
Bundle(path: "/Applications/InjectionIII.app/Contents/Resources/iOSInjection.bundle")?.load()
//for tvOS:
Bundle(path: "/Applications/InjectionIII.app/Contents/Resources/tvOSInjection.bundle")?.load()
//Or for macOS:
Bundle(path: "/Applications/InjectionIII.app/Contents/Resources/macOSInjection.bundle")?.load()
#endif
```

3.在你编辑的`ViewController`中添加如下方法

```objc
- (void)injected {

}
```

在该方法中添加你想要更改的代码，按`Command + S`运行。

**注意**

> 1.重新运行不会调用`injected`方法，因为该方法并未被实际调用
>
> 2.只能在模拟器中使用