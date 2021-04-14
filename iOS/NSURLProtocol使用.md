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

首先看一下NSURLProtocol类常用的方法

```objective-c
/// 注册该类，使之对URL加载系统可见
+ (BOOL)registerClass:(Class)protocolClass;

/// 取消该类
+ (void)unregisterClass:(Class)protocolClass;

/// 过滤方法 返回YES由该类处理请求，否则URL Loading System使用系统默认的行为处理
+ (BOOL)canInitWithRequest:(NSURLRequest *)request;

/// 在该方法中自定义网络请求, 对请求进行修改，如URL重定向、添加Header
/// 无需额外处理可直接返回request
+ (NSURLRequest *)canonicalRequestForRequest:(NSURLRequest *)request;

/// 判断两个请求是否相同，相同的话可以使用缓存数据，一般直接返回父类实现
+ (BOOL)requestIsCacheEquivalent:(NSURLRequest *)a toRequest:(NSURLRequest *)b;

/// 开始加载请求
- (void)startLoading;

/// 取消加载请求
- (void)stopLoading;

/// 给指定的请求设置与指定键相关联的属性
+ (void)setProperty:(id)value forKey:(NSString *)key inRequest:(NSMutableURLRequest *)request;

/// 返回与指定的请求中指定的关键字关联的属性。如果没有该key，返回nil
+ (nullable id)propertyForKey:(NSString *)key inRequest:(NSURLRequest *)request;

/// 移除给指定的请求的指定key相关联的属性
+ (void)removePropertyForKey:(NSString *)key inRequest:(NSMutableURLRequest *)request;


```

##### 使用步骤

1. 自定义子类继承于NSURLProtocol

   ![自定义子类继承于NSURLProtocol](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_04_14_15_13_46.png "自定义子类继承于NSURLProtocol")

2. 实现抽象类NSURLProtocol的方法

   ```objective-c
   
   ```

   

3. 在网络请求调用之前注册该类

## 参考文档

[CustomHTTPProtocol官方Demo](https://developer.apple.com/library/archive/samplecode/CustomHTTPProtocol/Introduction/Intro.html)

[NSURLProtocol 全攻略](https://blog.csdn.net/intheair100/article/details/80888742)