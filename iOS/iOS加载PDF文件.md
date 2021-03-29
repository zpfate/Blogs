本文主要介绍一下iOS加载PDF的方法。

## 使用webView加载

使用webView加载是最简单的方式，这是最大的优点，使用方便，并且可以加在网络文件，缺点就是效果体验不怎么好，单一的上下滑动效果。

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

如果出现了PDF乱码的情况，使用下面这个方法进行加载：

```objective-c
   [self.wkWebView loadData:data MIMEType:@"application/pdf" characterEncodingName:@"GBK" baseURL:fileUrl];	
```

## 使用QLPreviewController

使用QLPreviewController需要导入QuickLook框架

```objective-c
#import <QuickLook/QuickLook.h>		
```

遵循QLPreviewControllerDelegate和QLPreviewControllerDataSource，

使用方式如下：

```objective-c
@property (nonatomic, strong) QLPreviewController *previewController;
```



```objc
- (QLPreviewController *)previewController {
    if (!_previewController) {
        _previewController = [[QLPreviewController alloc] init];
        _previewController.delegate = self;
        _previewController.dataSource = self;
    }
    return _previewController;
}	
```

实现协议：

```objective-c
/// 返回文件数量
- (NSInteger)numberOfPreviewItemsInPreviewController:(QLPreviewController *)controller {
    return 1;
}

- (id<QLPreviewItem>)previewController:(QLPreviewController *)controller previewItemAtIndex:(NSInteger)index {
    NSString *pdfPath = [[NSBundle mainBundle] pathForResource:@"test" ofType:@"pdf"];
    return [NSURL fileURLWithPath:pdfPath];
}	
```

还有一种方式是使用UIDocumentInteractionController进行预览，效果与QLPreviewController一样，但是不支持多文件预览。

使用时也需要遵循UIDocumentInteractionControllerDelegate。

```objective-c
- (UIViewController *)documentInteractionControllerViewControllerForPreview:(UIDocumentInteractionController *)controller {
    return self;
}
- (UIView*)documentInteractionControllerViewForPreview:(UIDocumentInteractionController*)controller {
    return self.view;
}
- (CGRect)documentInteractionControllerRectForPreview:(UIDocumentInteractionController*)controller {
    return self.view.frame;
}
```



UIDocumentInteractionController继承于NSObject，而非UIViewController，使用该类自带的presentPreviewAnimated方法展示。

这两种的加载方法优势是，

1. 集成方便，并且自带系统分享功能

2. QLPreviewController支持多文件预览

劣势的话就是体验比较差，与webView类似的上下滑动的预览方式

## CGContexDrawPDFPage

如果我们需要一个比较好的体验效果，类似于翻书一样的动画效果，这个时候我们就需要使用CGContexDrawPDFPage来进行自己绘制。

1. 获取PDF文件

   ```objective-c
     /// 获取路径
       CFStringRef path = CFStringCreateWithCString(NULL, filePath.UTF8String, kCFStringEncodingUTF8);
       CFURLRef url = CFURLCreateWithFileSystemPath(NULL, path, kCFURLPOSIXPathStyle, NO);
       /// 获取PDF
       CGPDFDocumentRef document = CGPDFDocumentCreateWithURL(url);
       /// 主动释放
       CFRelease(path);
       CFRelease(url);	
   ```

   

2. 获取PDF页数

   ```objc
   NSInteger totoalPages = CGPDFDocumentGetNumberOfPages(document);
   ```

3. 获取页数内容

   ```objective-c
   CGPDFPageRef page = CGPDFDocumentGetPage(document, _currentPage);
   ```

4. 获取当前页的大小

   ```objective-c
   CGRect rect = CGPDFPageGetBoxRect(page, kCGPDFMediaBox);
   ```

5. 绘制当前页内容

   ```objective-c
   - (void)drawRect:(CGRect)rect {
       
       CGContextRef context = UIGraphicsGetCurrentContext();
       [[UIColor whiteColor] set];
       CGContextFillRect(context, rect);
       
       /// 转换坐标系
       CGContextTranslateCTM(context, 0.0, rect.size.height);
       CGContextScaleCTM(context, 1.0, -1.0);
       
       CGPDFPageRef page = CGPDFDocumentGetPage(document, _currentPage);
       CGAffineTransform pdfTransform = CGPDFPageGetDrawingTransform(page, kCGPDFCropBox, rect, 0, true);
       CGContextConcatCTM(context, pdfTransform);
       CGContextDrawPDFPage(context, page);
   }
   ```

6. 最后一步加上翻页动画

   ```objective-c
   /// 返回上一页
       [self transitionWithType:@"pageUnCurl" subtype:kCATransitionFromRight];
   /// 下一页
       [self transitionWithType:@"pageCurl" subtype:kCATransitionFromRight];
   
   ```

   ```objective-c
   - (void)transitionWithType:(NSString *)type subtype:(CATransitionSubtype)subType {
   
       CATransition *animation = [CATransition animation];
       /// 设置动画时长
       animation.duration = 0.8;
       /// 设置动画样式
       /***
        使用CATransitionType
        kCATransitionPush 推入效果
        kCATransitionMoveIn 移入效果
        kCATransitionReveal 截开效果
        kCATransitionFade 渐入渐出效果
        或者直接使用以下字符串:
        cube 方块
        suckEffect 三角
        rippleEffect 水波抖动
        pageCurl 上翻页
        pageUnCurl 下翻页
        oglFlip 上下翻转
        cameraIrisHollowOpen 镜头快门开
        cameraIrisHollowClose 镜头快门开
       */
       animation.type = type;
       
       /// 设置动画方向
       /**
        kCATransitionFromRight     从右边
        kCATransitionFromLeft      从左边
        kCATransitionFromTop       从上面
        kCATransitionFromBottom    从下面
        */
       if (subType) {
           animation.subtype = subType;
       }
       /// 动画的速度
       /** CAMediaTimingFunction:
           kCAMediaTimingFunctionLinear 匀速
           kCAMediaTimingFunctionEaseIn 慢进快出
           kCAMediaTimingFunctionEaseOut 快进慢出
           kCAMediaTimingFunctionEaseInEaseOut 慢进慢出 中间加速
           kCAMediaTimingFunctionDefault 默认
        */
       animation.timingFunction = [CAMediaTimingFunction functionWithName: kCAMediaTimingFunctionEaseIn];
       [self.layer addAnimation:animation forKey:@"animation"];
   }
   ```

   [附上简单的Demo](https://github.com/ZpFate/FileBrowser)

   





