# alloc和init

## alloc做了什么

![alloc](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1720764256400)

### callAlloc

```c
static ALWAYS_INLINE id
callAlloc(Class cls, bool checkNil, bool allocWithZone=false)// alloc 源码 第三步
{
#if __OBJC2__ //有可用的编译器优化
    
    // checkNil 为false，!cls 也为false ，所以slowpath 为 false，假值判断不会走到if里面，即不会返回nil
    if (slowpath(checkNil && !cls)) return nil;
    
    // 判断一个类是否有自定义的 +allocWithZone 实现，没有则走到if里面的实现
    if (fastpath(!cls->ISA()->hasCustomAWZ())) {
        return _objc_rootAllocWithZone(cls, nil);
    }
#endif

    // No shortcuts available. // 没有可用的编译器优化
    if (allocWithZone) {
        return ((id(*)(id, SEL, struct _NSZone *))objc_msgSend)(cls, @selector(allocWithZone:), nil);
    }
    return ((id(*)(id, SEL))objc_msgSend)(cls, @selector(alloc));
}
```



slowpath和fastpath理解

```c
//x很可能为真， fastpath 可以简称为 真值判断
#define fastpath(x) (__builtin_expect(bool(x), 1)) 
//x很可能为假，slowpath 可以简称为 假值判断
#define slowpath(x) (__builtin_expect(bool(x), 0)) 
```

> `__builtin_expect`指令是由`gcc`引入:
>
> * 目的：编译器可以对代码进行优化，以减少指令跳转带来的性能下降。即性能优化
>
> * 作用：允许程序员将最有可能执行的分支告诉编译器。
>
> * 指令的写法为：__builtin_expect(EXP, N)。表示 EXP==N的概率很大。
>
> * fastpath定义中__builtin_expect((x),1)表示 x 的值为真的可能性更大；即 执行if 里面语句的机会更大
>
> * slowpath定义中的__builtin_expect((x),0)表示 x 的值为假的可能性更大。即执行 else 里面语句的机会更大

### _objc_rootAllocWithZone

```objc
NEVER_INLINE
id _objc_rootAllocWithZone(Class cls, malloc_zone_t *zone __unused)
{
    // allocWithZone under __OBJC2__ ignores the zone parameter
    return _class_createInstanceFromZone(cls, 0, nil,
                                         OBJECT_CONSTRUCT_CALL_BADALLOC);
}
```

### _class_createInstanceFromZone

```c
_class_createInstanceFromZone(Class cls, size_t extraBytes, void *zone,
                              int construct_flags = OBJECT_CONSTRUCT_NONE,
                              bool cxxConstruct = true,
                              size_t *outAllocatedSize = nil)// alloc 源码 第五步
{
    ASSERT(cls->isRealized()); //检查是否已经实现

    // Read class's info bits all at once for performance
    //一次性读取类的位信息以提高性能
    bool hasCxxCtor = cxxConstruct && cls->hasCxxCtor();
    bool hasCxxDtor = cls->hasCxxDtor();
    bool fast = cls->canAllocNonpointer();
    size_t size;

    //计算需要开辟的内存大小，传入的extraBytes 为 0
    size = cls->instanceSize(extraBytes);
    if (outAllocatedSize) *outAllocatedSize = size;

    id obj;
    if (zone) {
        obj = (id)malloc_zone_calloc((malloc_zone_t *)zone, 1, size);
    } else {
        //申请内存
        obj = (id)calloc(1, size);
    }
    if (slowpath(!obj)) {
        if (construct_flags & OBJECT_CONSTRUCT_CALL_BADALLOC) {
            return _objc_callBadAllocHandler(cls);
        }
        return nil;
    }

    if (!zone && fast) {
        // 将 cls类 与 obj指针（即isa） 关联
        obj->initInstanceIsa(cls, hasCxxDtor);
    } else {
        // Use raw pointer isa on the assumption that they might be
        // doing something weird with the zone or RR.
        obj->initIsa(cls);
    }

    if (fastpath(!hasCxxCtor)) {
        return obj;
    }

    construct_flags |= OBJECT_CONSTRUCT_FREE_ONFAILURE;
    return object_cxxConstructFromClass(obj, cls, construct_flags);
}
```



这部分代码实现是alloc的核心操作，主要做了三件事:

>- `cls->instanceSize`：计算需要开辟的内存空间大小
>- `calloc`：申请内存，返回地址指针
>- `obj->initInstanceIsa`：将类与`isa`关联

![_class_createInstanceFromZone](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1720764722652)



### instanceSize

![instanceSize](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1720765326791)

```c
size_t instanceSize(size_t extraBytes) const {
    //编译器快速计算内存大小
    if (fastpath(cache.hasFastInstanceSize(extraBytes))) {
        return cache.fastInstanceSize(extraBytes);
    }
    
    // 计算类中所有属性的大小 + 额外的字节数0
    size_t size = alignedInstanceSize() + extraBytes;
    // CF requires all objects be at least 16 bytes.
    //如果size 小于 16，最小取16
    if (size < 16) size = 16;
    return size;
}
```

```c
size_t fastInstanceSize(size_t extra) const
{
    ASSERT(hasFastInstanceSize(extra));

    //Gcc的内建函数 __builtin_constant_p 用于判断一个值是否为编译时常数，如果参数EXP 的值是常数，函数返回 1，否则返回 0
    if (__builtin_constant_p(extra) && extra == 0) {
        return _flags & FAST_CACHE_ALLOC_MASK16;
    } else {
        size_t size = _flags & FAST_CACHE_ALLOC_MASK;
        // remove the FAST_CACHE_ALLOC_DELTA16 that was added
        // by setFastInstanceSize
        //删除由setFastInstanceSize添加的FAST_CACHE_ALLOC_DELTA16 8个字节
        return align16(size + extra - FAST_CACHE_ALLOC_DELTA16);
    }
}
```

```c
// 16字节对齐算法
static inline size_t align16(size_t x) {
    return (x + size_t(15)) & ~size_t(15);
}
```

## init

### 类方法init

```c
+ (id)init {
    return (id)self;
}
```

### 实例方法init

```c
- (id)init {
    return _objc_rootInit(self);
}
```

### _objc_rootInit

```c
id _objc_rootInit(id obj)
{
    // In practice, it will be hard to rely on this function.
    // Many classes do not properly chain -init calls.
    return obj;
}
```

可以看到，init方法直接返回了obj。

## new

```c

/// new方法和alloc不一样的在于callAlloc的allocWithZone参数为false
+ (id)new {
    return [callAlloc(self, false/*checkNil*/) init];
}
```

`new`函数中直接调用了`callAlloc`函数（即`alloc`中分析的函数），且调用了`init`函数，所以可以得出`new` 其实就等价于 `[alloc init]`的结论

## 总结

1. 通过分析可以得出，alloc主要做的事情就是计算申请的内存大小，开辟内存，关联isa指针。

   new方法直接调用了callAlloc和init方法。

2. init用于类实例的初始化



