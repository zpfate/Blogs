## 编写Swift

### 1.REPL

在终端中直接输入`swift`敲回车就可以进入REPL，可以在里面编写`Swift`代码

![REPL](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_04_15_17_02.png "REPL")

### 2.Xcode Playground

打开`Xcode`选择第一个`playground`就可以进行编写`Swift`代码，（当然也可以新建一个工程编写Swift）

在`playground`中想要使用`liveView`的功能需要

```
import PlaygroundSupport
```

## Optional

### 为什么`Swift`引入`Optional`

* Objective-C中nil是无类型空指针, 无法知道具体是哪个类型
* 数组字典集合不允许放nil
* Objective-C所有对象都可以为nil
* 运行的时候不知道是不是nil,需要先进行空值判断
* Objective-C中nil只能对象使用, 不能用在其他地方, 比如数组在查找元素没有找到 使用NSNotFound来标识没有找到

### Optional使用

#### 可选链

```swift
let str: String? = "abc"
let count = str?.count
```

需要注意的是这里获得的`count`也是可选类型

#### Optional绑定

```swift
let str: String? = "abc"
if let actualStr = str {
  let count = actualStr.count
  print(count)
}
```

#### 强制解包

```swift
let count = str!.count
```

使用感叹号代表默认`str`是有值的，如果为*nil*，程序会挂掉

#### 隐式解包

在声明的时候使用`!`代表默认始终有值

```swift
let str: String! = "abc"
let count = str.count
```



### Optional实现原理

`Optional`是`Swift`标准库里的一个枚举类型

![image-20200529092932060](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_05_29_09_29_32.png)

`Optioanl.none`就是为nil

`Optional.some`就是包装了实际的值

```swift
var optionalString: Optional<String> = "abc"
if let actualOptionalString = optionalString {
    let count = actualOptionalString.count
    print(count)
}
```

理论上可以直接通过`unsafelyUnwrapped`获取可选项的值

```swift
let optionalStringCount = optionalString.unsafelyUnwrapped.count
```



## 字符串处理

### 声明字符串

```swift
let str: String = "abc"
```

#### 声明多行字符串

```swift
let numbers = """
1
2
3
4
5
"""
print("numbers === \(numbers)")
```

![多行字符串声明定义](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_04_15_20_52.png)

#### 字符串插值

类似`NSString StringWithFormat`的占位符用法

```swift
print("numbers === \(numbers)")
```

#### Raw String

```swift
print(#"6 times 7 is \(6 * 7)."#)
```

运行结果：

![image-20200902201943162](https://i.loli.net/2020/09/02/nSbEWlY1t6PVQTc.png)

在反斜杠后面添加与首尾同等数量的#号， 使插值生效

```
print(#"6 times 7 is \#(6 * 7)."#)
```

运行结果：

![image-20200902202244031](https://i.loli.net/2020/09/02/Ake61FoLYXzcEJ5.png)

## 注释

swift支持注释嵌套

swift Playground注释支持markup语法 与markdown相似

```
单行注释 
//:开启markup
多行注释
/*:
开启markup
*/
```

点击上方菜单Editor ->show Rendered

