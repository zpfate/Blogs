

# KVC&KVO底层探究

## KVC

`KVC(Key-Value Coding)`是利用`NSKeyValueCoding`非正式协议实现的一种机制，对象采用这种机制来提供对其属性的间接访问。

写下KVC代码并点击跟进setValue会发现NSKeyValueCoding是在Foundation框架下

* KVC通过对`NSObject`的扩展来实现的——所有集成了`NSObject`的类可以使用KVC
* `NSArray、NSDictionary、NSMutableDictionary、NSOrderedSet、NSSet`等也遵守KVC协议
* 少数类型(结构体)以外都可以使用KVC

### KVC常用方法

```objc
// 通过 key 设值
- (void)setValue:(nullable id)value forKey:(NSString *)key;
// 通过 key 取值
- (nullable id)valueForKey:(NSString *)key;
// 通过 keyPath 设值
- (void)setValue:(nullable id)value forKeyPath:(NSString *)keyPath;
// 通过 keyPath 取值
- (nullable id)valueForKeyPath:(NSString *)keyPath;
```

`NSKeyValueCoding`类别的其它方法

```objc
// 默认为YES。 如果返回为YES,如果没有找到 set<Key> 方法的话, 会按照_key, _isKey, key, isKey的顺序搜索成员变量, 返回NO则不会搜索
+ (BOOL)accessInstanceVariablesDirectly;
// 键值验证, 可以通过该方法检验键值的正确性, 然后做出相应的处理
- (BOOL)validateValue:(inout id _Nullable * _Nonnull)ioValue forKey:(NSString *)inKey error:(out NSError **)outError;
// 如果key不存在, 并且没有搜索到和key有关的字段, 会调用此方法, 默认抛出异常。两个方法分别对应 get 和 set 的情况
- (nullable id)valueForUndefinedKey:(NSString *)key;
- (void)setValue:(nullable id)value forUndefinedKey:(NSString *)key;
// setValue方法传 nil 时调用的方法
// 注意文档说明: 当且仅当 NSNumber 和 NSValue 类型时才会调用此方法 
- (void)setNilValueForKey:(NSString *)key;
// 一组 key对应的value, 将其转成字典返回, 可用于将 Model 转成字典
- (NSDictionary<NSString *, id> *)dictionaryWithValuesForKeys:(NSArray<NSString *> *)keys;
```



## KVO相关

Key-Value Observing 键值监听

### 实现原理

![image-20220314222151383](https://github.com/zpfate/uPic/2022%2003%20141647269033.png)

在程序运行中，动态创建一个NSKVONotifiying_XXX类