## NSURLProtocol介绍

NSURLProtocol是一个抽象类，作为URL Loading System系统的一部分，能够帮助我们拦截所有的URL Loading System的请求，在此进行各种自定义的操作，是网络层实现AOP的利器。。

> URL Loading System可以理解为一系列的类和协议，主要用来访问通过URL定位获取的资源。详细看下图
>
> AOP：面向切面编程。

![URL Loading System](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_04_13_10_31_50.webp "URL Loading System")

## NSURLProtocol能拦截的网络请求

UIWebView，NSURLConnection， NSURLSession都可以被NSURLProtocol拦截，包含），除了WkWebView，因为它使用的是Webkit

* UIWebView
* NSURLConnection
* NSURLSession
* 一些第三方库库（AFNetworking等），本质上还是有NSURLSession实现的，拦截与其他稍有区别

拦截下网络请求，我们就可以做一些自定义的操作，类似于重定向网络请求，全局的网络请求设置等等，下面来看看如何使用NSProtocol

## 使用NSURLProtocol

1. 自定义子类继承于NSURLProtocol

2. 实现抽象类NSURLProtocol的方法

   ```objective-c
   + (BOOL)canInitWithRequest:(NSURLRequest *)request;
   + (NSURLRequest *)canonicalRequestForRequest:(NSURLRequest *)request;
   + (BOOL)requestIsCacheEquivalent:(NSURLRequest *)a toRequest:(NSURLRequest *)b;
   - (void)startLoading;
   - (void)stopLoading;
   + (nullable id)propertyForKey:(NSString *)key inRequest:(NSURLRequest *)request;
   + (void)setProperty:(id)value forKey:(NSString *)key inRequest:(NSMutableURLRequest *)request;
   + (void)removePropertyForKey:(NSString *)key inRequest:(NSMutableURLRequest *)request;
   ```

   

## 参考文档

[CustomHTTPProtocol官方Demo](https://developer.apple.com/library/archive/samplecode/CustomHTTPProtocol/Introduction/Intro.html)

[NSURLProtocol 全攻略](https://blog.csdn.net/intheair100/article/details/80888742)