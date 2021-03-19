## 谓词

首先，我们可以看看官方文档的解释

![NSPredicate官方描述](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_17_10_35_15.png "NSPredicate官方描述")

大概意思是： 一种逻辑条件的定义，用于约束查找或者内存中筛选的搜索条件。个人理解就是一个用来匹配查询结果的条件。

### NSPredicate的API

```objective-c
// 通过表达式创建谓词实例
+ (NSPredicate *)predicateWithFormat:(NSString *)predicateFormat, ...;

// 表达式中的%@从arguments中取值
+ (NSPredicate *)predicateWithFormat:(NSString *)predicateFormat argumentArray:(nullable NSArray *)arguments;

// 表达式中的%@占位符将从argList中取值
+ (NSPredicate *)predicateWithFormat:(NSString *)predicateFormat arguments:(va_list)argList;

// 通过一个元数据查询字符串来创建谓词实例，仅macos可以使用
+ (nullable NSPredicate *)predicateFromMetadataQueryString:(NSString *)queryString API_AVAILABLE(macos(10.9)) API_UNAVAILABLE(ios, watchos, tvos);

// 创建一个固定结果为BOOL值的谓词实例
+ (NSPredicate *)predicateWithValue:(BOOL)value;    

// 根据block的回调的返回值创建一个结果为YES/NO的谓词实例
+ (NSPredicate*)predicateWithBlock:(BOOL (^)(id _Nullable evaluatedObject, NSDictionary<NSString *, id> * _Nullable bindings))block API_AVAILABLE(macos(10.6), ios(4.0), watchos(2.0), tvos(9.0));

// 获取谓词的表达式，只读属性
@property (readonly, copy) NSString *predicateFormat;    

// 用常量值代替变量，用字典中的键值对替换用$声明的变量
- (instancetype)predicateWithSubstitutionVariables:(NSDictionary<NSString *, id> *)variables;    // substitute constant values for variables

// 比较评估对象是否符合谓词
- (BOOL)evaluateWithObject:(nullable id)object;    

// 比较评估对象是否符合谓词
- (BOOL)evaluateWithObject:(nullable id)object substitutionVariables:(nullable NSDictionary<NSString *, id> *)bindings API_AVAILABLE(macos(10.5), ios(3.0), watchos(2.0), tvos(9.0)); 

// 强制使用已安全解码的谓词以进行评估
- (void)allowEvaluation API_AVAILABLE(macos(10.9), ios(7.0), watchos(2.0), tvos(9.0)); 
```

### 常用方法

上面官方暴露的API我们可以看到创建NSPredicate基本都需要一个表达式

#### 创建谓词

```objective-c
 NSPredicate *predicate1 = [NSPredicate predicateWithFormat:@"SELF CONTAINS 'world'"];
 NSPredicate *predicate2 = [NSPredicate predicateWithFormat:@"SELF CONTAINS %@", @"world"];
 NSPredicate *predicate3 = [NSPredicate predicateWithFormat:@"SELF CONTAINS %@" argumentArray:@[@"world"]];
 NSLog(@"predicate1谓词表达式:%@", predicate1.predicateFormat);
 NSLog(@"predicate2谓词表达式:%@", predicate2.predicateFormat);
 NSLog(@"predicate3谓词表达式:%@", predicate3.predicateFormat);
```

打印结果为：![谓词创建实例比较](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_17_14_54_01.png "谓词创建实例比较")

#### 创建谓词并使用

```objective-c
NSString *str = @"Hello, world";
NSPredicate *predicate = [NSPredicate predicateWithFormat:@"SELF CONTAINS 'world'"];
NSLog(@"谓词评估结果:%hhd",[predicate evaluateWithObject:str]);
```

打印结果为：**谓词评估结果:1**

```objective-c
NSPredicate *predicate = [NSPredicate predicateWithFormat:@"SELF CONTAINS $variable"];
predicate = [predicate predicateWithSubstitutionVariables:@{@"variable":@"world"}]; NSLog(@"谓词表达式为%@", [predicate predicateFormat]);
```

```objective-c
NSPredicate *predicate = [NSPredicate predicateWithBlock:^BOOL(NSArray * _Nullable evaluatedObject, NSDictionary<NSString *,id> * _Nullable bindings) {
      // 通过evaluatedObject和bindings来判断是否过滤
      if ([evaluatedObject containsObject:[bindings objectForKey:@"key"]]) {
         	return YES;
      } else {
          return NO;
       }
}];
BOOL result = [predicate evaluateWithObject:@[@"1", @"2"] substitutionVariables:@{@"key":@"2"}];
NSLog(@"谓词评估结果%hhd", result);
```

#### NSPredicate的扩展

```
@interface NSArray<ObjectType> (NSPredicateSupport)

// 对数组中所有元素进行谓词评估，返回数组
- (NSArray<ObjectType> *)filteredArrayUsingPredicate:(NSPredicate *)predicate; 

@end
@interface NSMutableArray<ObjectType> (NSPredicateSupport)
// 移除所有不符合谓词评估的数组元素
- (void)filterUsingPredicate:(NSPredicate *)predicate; 
@end

```

NSSet相关类中也有类似的方法，使用与数组这边雷同。

### 谓词表达式语法

**可以在关键字后面加在[cd]，c表示忽略大小写，即CAFE与cafe为相同字符串；d表示忽略读音，即café与cafe为相同字符串。**

#### 常用字面语义

| 表达式    | 用途                              |
| --------- | --------------------------------- |
| SELF      | 谓词评估的对象                    |
| 'text'    | 字符串text，使用双引号需要转义符\ |
| TRUE、YES | 逻辑真                            |
| FALSE、NO | 逻辑假                            |

 ##### 占位符

| 表达式 | 用途   |
| ------ | ------ |
| %K     | 属性名 |
| %@     | 属性值 |

```objective-c
NSString *name = @"name";
NSString *value = @"value";
NSPredicate *predicate = [NSPredicate predicateWithFormat:@"%K CONTAINS %@", name, value];
NSLog(@"谓词表达式%@", [predicate predicateFormat]);
BOOL result = [predicate evaluateWithObject:@{@"name" : @"value111"}];
NSLog(@"谓词结果为%hhd", result);
```

输出结果为：**谓词表达式name CONTAINS "value"**

​					  **谓词结果为1**

#### 基本比较运算

| 表达式   | 用途                     |
| -------- | ------------------------ |
| =、==    | 等于                     |
| \>=、 => | 大于或等于               |
| <=、=<   | 小于或等于               |
| >        | 大于                     |
| <        | 小于                     |
| !=、<>   | 不等于                   |
| BETWEEN  | 值在区间内(包含边界数值) |

这边附上一个`BETWEEN`的使用例子

```objective-c
NSPredicate *predicate = [NSPredicate predicateWithFormat:@"SELF BETWEEN {200, 250}"];
NSLog(@"谓词评估结果:%hhd",  [predicate evaluateWithObject:@(212)]); // 谓词评估结果:1
```

#### 布尔值

| 表达式         | 用途 |
| -------------- | ---- |
| TRUEPREDICATE  | 真   |
| FALSEPREDICATE | 假   |

[NSPredicate predicateWithValue:YES]`的表达式实际上就是`TRUEPREDICATE

##### 逻辑运算

| 表达式   | 用途 |
| -------- | ---- |
| AND、&&  | 与   |
| OR、\|\| | 或   |
| NOT、!   | 非   |

##### 字符串比较

| 表达式     | 用途                       |
| ---------- | -------------------------- |
| BEGINSWITH | 字符串以...开始            |
| ENDSWITH   | 字符串以...结束            |
| CONTAINS   | 字符串包含...              |
| LIKE       | 字符串相等                 |
| MATCHES    | 字符串匹配右边的正则表达式 |

##### 集合操作

| 表达式                     | 用途                             |
| -------------------------- | -------------------------------- |
| ANY、SOME                  | 集合中是否至少有一个元素满足条件 |
| ALL                        | 集合中所有元素满足条件           |
| NONE                       | 集合中所有元素均不满足条件       |
| IN                         | 被集合包含                       |
| SELF[index]                | 取对应索引的元素                 |
| SELF[FIRST]                | 第一个元素                       |
| SELF[LAST]                 | 最后一个元素                     |
| SELF[SIZE]                 | 元素的个数                       |
| SELF['key']或SELF.key或key | 字典对应key的元素                |

[参考链接1](https://www.jianshu.com/p/9f19a0842f52)

[参考链接2](https://www.jianshu.com/p/26da5d4e86b8)

***如有错误，欢迎指正***

iOS开发QQ交流群：`568571839 `







