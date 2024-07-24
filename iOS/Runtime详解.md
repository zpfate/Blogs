## Runtime

`Objective-C`是一门动态性编程语言，动态性由`Runtime API`支撑

`Runtime`是一套c语言的`api`，封装了很多动态性相关的函数

平时编写的`Objc`代码，底层都是转成了`Runtime API`进行调用

### isa

在arm64位架构之前，isa就是一个普通的指针，存储着`Class`、`Meta-Class`对象的内存地址

从arm64架构开始，对isa进行了优化，变成了一个共用体（union）结构，还使用位域存储更多信息

![image-20220321154521223](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1647848721.png)

### Class的结构

![image-20220324093942864](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648085983.png)

### objc_msgSend

* OC中的方法调用，其实都是转换为objc_msgSend函数的调用
* objc_msgSend的执行流程可分为三大阶段
  * 消息发送
  * 动态方法解析
  * 消息转发

![image-20220324110152676](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648090913.png)

#### 动态解析

![image-20220324140711722](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648102031.png "动态解析流程")

```objc
// 动态方法解析
+ (BOOL)resolveInstanceMethod:(SEL)sel {
    if (sel == @selector(test)) {
        Method method = class_getInstanceMethod(self, @selector(instanceTest));
        class_addMethod(
                        self,
                        sel,
                        method_getImplementation(method),
                        method_getTypeEncoding(method)
                        );
        return YES;
       
    }
    return [super resolveInstanceMethod:sel];
}

+ (BOOL)resolveClassMethod:(SEL)sel {
    
    if (sel == @selector(test)) {
        Method method = class_getClassMethod(self, @selector(classTest));
        class_addMethod(
                        object_getClass(self),
                        sel,
                        method_getImplementation(method),
                        method_getTypeEncoding(method)
                        );
        return YES;
    }
    return [super resolveClassMethod:sel];
}
```

![image-20220325101031211](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648174231.png "Objective-C type encodings")



#### 消息转发

NSInvocation封装了一个方法调用，包括了*方法调用者、方法、方法参数*。

![image-20220324172139395](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648113699.png)



```objective-c
// 消息转发
- (id)forwardingTargetForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
//        return [[Cat alloc] init];
        return nil;
    }
    return [super forwardingTargetForSelector: aSelector];
}

// 类方法的消息转发
+ (id)forwardingTargetForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
//        return [Cat class];
    }
    return [super forwardingTargetForSelector:aSelector];
}
// 方法签名: 返回值类型 参数类型
- (NSMethodSignature *)methodSignatureForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
        Method method = class_getInstanceMethod(object_getClass(self), @selector(instanceTest));
        return [NSMethodSignature signatureWithObjCTypes:method_getTypeEncoding(method)];
        // 也可以这么生成方法签名
        return [[[Cat alloc] init] methodSignatureForSelector:aSelector];
    }
    return [super methodSignatureForSelector:aSelector];
}

+ (NSMethodSignature *)methodSignatureForSelector:(SEL)aSelector {
    if (aSelector == @selector(test)) {
        Method method = class_getClassMethod([Cat class], @selector(test));
        return [NSMethodSignature signatureWithObjCTypes:method_getTypeEncoding(method)];
        // 也可以这么生成方法签名
//        return [[[Cat alloc] init] methodSignatureForSelector:aSelector];
    }
    return [super methodSignatureForSelector:aSelector];
}
// NSInvocation封装了一个方法调用
// anInvocation.target 方法调用者
// anInvocation.selector 方法名
// [anInvocation getArgument:NULL atIndex:0] 方法参数 参数顺序receiver,selector,other
// [anInvocation getReturnValue:&value]; 获取返回值
- (void)forwardInvocation:(NSInvocation *)anInvocation {
//    anInvocation.target = [[Cat alloc] init];
//    [anInvocation invoke];
    [anInvocation invokeWithTarget:[[Cat alloc] init]];
}

+ (void)forwardInvocation:(NSInvocation *)anInvocation {
    [anInvocation invokeWithTarget:[Cat class]];
}
```

![img](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1721721172142)



#### super

![image-20220325151039767](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648192240.png)

[super message]底层实现是消息发送的时候，从父类开始寻找方法实现，消息接收者仍然是子类对象。

super底层调用的是objc_msgSendSuper函数

![image-20220326213958881](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648302138.png "super底层调用的是objc_msgSendSuper")

该函数需要传入一个objc_super的结构体

![image-20220326214837125](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648302519.png "objc_super")

最后还是需要看class与superClass方法的内部实现:

![image-20220326220001204](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648303204.png "class和superClass的实现")

#### isKindOfClass、isMemberOfClass

两个方法的源码实现如下：

![image-20220326220222191](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648303343.png)

可以明显的看出，类方法（+方法）会获取元类对象，所以后面比较的也应该是元类对象。

![image-20220328100112973](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648432873.png)

### runtime应用

* 通过runtime找出控件私有成员变量 kvc修改属性
* 字典转模型
* 替换方法， method_swizzling方法交换（hook系统原生实现）
* 自动归档解档
* 关联对象添加属性
* 利用消息转发机制解决方法找不到异常管理

<font color = red> hook类簇:</font>

![image-20220329093504754](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1648517705.png)

