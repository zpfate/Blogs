本文主要介绍iOS App加载常用格式文件的问题，比如说TXT，PDF文件等。

## 使用webView加载

使用webView加载是最简单的方式，这是最大的优点，使用方便，缺点就是效果体验不怎么好。

***注意：***

UIWebview组件已经在今年被apple禁用，使用webView无法提交App Store，本文下面使用的都是基于`WKWebView`。

首先看一下 webView加载的方法，

### \- (nullable WKNavigation *)loadFileURL:(NSURL *)URL allowingReadAccessToURL:(NSURL *)readAccessURL

![921EC7F6-6D72-4D7E-88CA-120F18C50BBF](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_09_27_17_20_59.png)

该方法第一个参数URL自不用说，是我们想要加载的文件的路径，注意看apple的解释，这是一个`file URL`,我们应该使用下面这个方法的结果作为参数。

```objective-c
    NSURL *fileUrl = [NSURL fileURLWithPath:@"path"];
```

第二个参数readAccessURL，是用来传入页面引用资源的路径，正常是前面`File URL`的上级目录，备注也写的很详细，如果该参数传入的是单一的文件，那么只有这个文件能被WebKit加载。





