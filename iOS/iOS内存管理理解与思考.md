#### alloc/retain/release/dealloc

| 对象操作       | Objective-C Method           |
| -------------- | ---------------------------- |
| 生成并持有对象 | alloc/new/copy/mutableCopy等 |
| 持有对象       | retain                       |
| 释放对象       | release                      |
| 废弃对象       | dealloc                      |

使用类方法可以取得谁都不持有的对象，通过`autorelease`实现，比如NSMutableArray的array方法等，大致实现如下：

```objective-c
+ (id)object {
	id obj = [[NSObject alloc] init];
	[obj autorelease];
	return obj;
}
```



* Autorelease：

  1、生成并持有NSAutoreleasePool对象

  2、调用已分配的autorelease实例方法

  3、废弃NSAutoreleasePool对象

