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

/// 开始加载请求，在该方法中发起一个新的请求
- (void)startLoading;

/// 取消加载请求
- (void)stopLoading;

/// 给指定的请求设置与指定键相关联的属性
+ (void)setProperty:(id)value forKey:(NSString *)key inRequest:(NSMutableURLRequest *)request;

/// 返回与指定的请求中指定的关键字关联的属性。如果没有该key，返回nil
+ (nullable id)propertyForKey:(NSString *)key inRequest:(NSURLRequest *)request;

/// 移除给指定的请求的指定key相关联的属性
+ (void)removePropertyForKey:(NSString *)key inRequest:(NSMutableURLRequest *)request;

/// 注册该类
+ (BOOL)registerClass:(Class)protocolClass;

/// 取消注册
+ (void)unregisterClass:(Class)protocolClass;
```

##### 使用步骤

1. 自定义子类继承于NSURLProtocol

   ![自定义子类继承于NSURLProtocol](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_04_14_15_13_46.png "自定义子类继承于NSURLProtocol")

2. 实现抽象类NSURLProtocol的方法

   ```objective-c
   /// 是否接受处理这个请求 返回NO则由URL Loading System使用系统默认的行为处理
   /// @param request 请求
   + (BOOL)canInitWithRequest:(NSURLRequest *)request {
       
       BOOL shouldAccept;
       NSURL *url = request.URL;
       shouldAccept = (request != nil) && (url != nil);
       
       if (shouldAccept) {
           /// 防止递归调用
           shouldAccept = ![self propertyForKey:TFCusomterProtocolKey inRequest:request];
       }
       
       if (shouldAccept) {
           NSString *scheme = url.scheme;
           shouldAccept = [scheme isEqualToString:@"http"];
       }
       return shouldAccept;
   }
   
   /// 自定义网络请求, 对请求进行修改，如URL重定向、添加Header
   /// 无需额外处理可直接返回request
   /// @param request 请求
   + (NSURLRequest *)canonicalRequestForRequest:(NSURLRequest *)request {
       /// 如果是http请求改成https
       if ([request.URL.scheme isEqualToString:@"http"]) {
           NSMutableURLRequest *mutableRequest = [request mutableCopy];
           NSString *urlString = mutableRequest.URL.absoluteString;
           urlString = [urlString stringByReplacingOccurrencesOfString:@"http" withString:@"https"];
           mutableRequest.URL = [NSURL URLWithString:urlString];
           return mutableRequest;
       }
       return request;
   }
   
   /// 此方法抽象类提供了默认实现,重写可以直接调用父类
   /// 一般不做特殊处理, 直接返回父类实现
   + (BOOL)requestIsCacheEquivalent:(NSURLRequest *)a toRequest:(NSURLRequest *)b {
       
       BOOL result = [super requestIsCacheEquivalent:a toRequest:b];
       return result;
   }
   
   /// 开始网络请求
   - (void)startLoading {
       
       NSMutableURLRequest *recursiveRequest = [[self request] mutableCopy];
       [[self class] setProperty:@YES forKey:TFCusomterProtocolKey inRequest:recursiveRequest];
       
       NSURLSessionConfiguration *configuration = [NSURLSessionConfiguration defaultSessionConfiguration];
       
       NSMutableArray *protocolClasses = [configuration.protocolClasses mutableCopy];
       [protocolClasses addObject:self];
       configuration.protocolClasses = @[self.class];
       
       self.session = [NSURLSession sessionWithConfiguration:configuration delegate:self delegateQueue:nil];
       NSURLSessionDataTask *task = [self.session dataTaskWithRequest:recursiveRequest];
       [task resume];
   }
   
   /// 停止相应请求
   - (void)stopLoading {
       
       [self.session invalidateAndCancel];
       self.session = nil;
   }
   ```

   

3. 实现NSURLSessionDataDelegate，当网络请求接收到服务端的响应时，将其通过"NSURLProtocolClient"协议，转发给URL Loading System

   ```objective-c
   #pragma mark -- NSURLSessionDataDelegate
   
   /// 接收到服务响应时调用的方法
   - (void)URLSession:(NSURLSession *)session dataTask:(NSURLSessionDataTask *)dataTask didReceiveResponse:(NSURLResponse *)response completionHandler:(void (^)(NSURLSessionResponseDisposition))completionHandler {
       [[self client] URLProtocol:self didReceiveResponse:response cacheStoragePolicy:NSURLCacheStorageAllowed];
       
       completionHandler(NSURLSessionResponseAllow);
   }
   
   ///接收到服务器返回数据的时候会调用该方法，如果数据较大那么该方法可能会调用多次
   - (void)URLSession:(NSURLSession *)session dataTask:(NSURLSessionDataTask *)dataTask didReceiveData:(NSData *)data {
       [[self client] URLProtocol:self didLoadData:data];
   }
   
   /// 当请求完成(成功|失败)的时候会调用该方法，如果请求失败，则error有值
   - (void)URLSession:(NSURLSession *)session task:(NSURLSessionTask *)task didCompleteWithError:(NSError *)error {
       
       if (error) {
           [[self client] URLProtocol:self didFailWithError:error];
       } else {
           [[self client] URLProtocolDidFinishLoading:self];
       }
   }
   ```

4. 在网络请求调用之前注册该类

   ```objective-c
   [NSURLProtocol registerClass:self];
   ```

   

#### 注意事项

对于AFNetworking这种网络库，可以采用hook掉NSURLSessionConfiguration的protocolClasses方法

#### 测试Demo

[NSURLProtocolDemo](https://github.com/zpfate/BlogDemo-iOS/tree/protocol)

## 参考文档

[CustomHTTPProtocol官方Demo](https://developer.apple.com/library/archive/samplecode/CustomHTTPProtocol/Introduction/Intro.html)

[NSURLProtocol 全攻略](https://blog.csdn.net/intheair100/article/details/80888742)