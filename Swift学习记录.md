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

