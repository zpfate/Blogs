## copy和mutableCopy

1. copy 不可变拷贝， 产生不可变副本
2. mutableCopy可变拷贝，产生可变副本

|             | NSString                                             | NSMutableString                                      | NSArray                                             | NSMutableArray                                      | NSDictionary                                             | NSMutableDictionary                                      |
| ----------- | ---------------------------------------------------- | ---------------------------------------------------- | --------------------------------------------------- | --------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| copy        | NSString<br /><font color = green>浅拷贝</font>      | NSString<br /><font color = red>深拷贝</font>        | NSArray<br /><font color = green>浅拷贝</font>      | NSArray<br /><font color = red>深拷贝</font>        | NSDictionary<br /><font color = green>浅拷贝</font>      | NSDictionary<br /><font color = red>深拷贝</font>        |
| mutableCopy | NSMutableString<br /><font color = red>深拷贝</font> | NSMutableString<br /><font color = red>深拷贝</font> | NSMutableArray<br /><font color = red>深拷贝</font> | NSMutableArray<br /><font color = red>深拷贝</font> | NSMutableDictionary<br /><font color = red>深拷贝</font> | NSMutableDictionary<br /><font color = red>深拷贝</font> |



## weak

weak是弱引用，用weak来修饰的引用对象的计数器不会增加，而且weak会在引用对象释放的时候，自动置为nil

#### 当 weak 指向的对象被释放时，如何让 weak 指针置为 nil 的呢

* 调用 objc_release
* 因为对象的引用计数为0，所以执行 dealloc
* 在 dealloc 中，调用了 _objc_rootDealloc 函数
* 在 _objc_rootDealloc 中，调用了 object_dispose 函数
* 调用 objc_destructInstance
* 最后调用 objc_clear_deallocating,详细过程如下：
  a. 从 weak 表中获取废弃对象的地址为键值的记录
  b. 将包含在记录中的所有附有 weak 修饰符变量的地址，赋值为 nil
  c. 将 weak 表中该记录删除
  d. 从引用计数表中删除废弃对象的地址为键值的记录

