### 声明

```swift
let emptyString = ""
let str = "abc"
let str1: String = "abc"
// str2赋值后才可以使用 并且不能改变
let str2: String? 

```

#### 声明多行字符串

```swift
let numbers = let numbers = """
1
2
3
4
5
"""
print("numbers === \(numbers)")
```

输出结果如下所示：

![多行字符串输出结果](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_09_03_09_34_59.png "多行字符串输出结果")

### 字符串插值

使用`\()`来进行插值，与`Objective-C`的`stringWithFormat`方法类似，其实就是使用占位符格式化输出.

```
let str = "6 times 7 is \(6 * 7)."
print(str)
```

输出结果：![输出插值字符串](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_09_03_10_25_48.png "输出插值字符串")

插值使用起来确实比`OC`的更加方便灵活

```swift
let str: String = "abc"
print("str is \(str)")
```

#### Raw String

`swift`的转义写法，可以在字符串中输出特殊字符

```swift
print(#"6 times 7 is \(6 * 7)."#)
```

输出结果如下所示：

![Raw String输出](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_09_03_10_41_04.png "Raw String输出")

这个时候如果想要使插值生效的话，可以在反斜杠后面添加与首尾同等数量的#

```swift
print(#"6 times 7 is \#(6 * 7)."#)	
```

输出结果：![Raw String使插值生效](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_09_03_10_51_20.png "Raw String使插值生效")

### 字符串处理

#### 判空

```swift
let emptyString = ""
if emptyString.isEmpty {
    print("emptyString is empty")
}
```

#### 取出字符串中某个字符

1.取出字符串中某个字符

```
let welcome = "Hello, world!"
// 取出第一个字符
welcome[welcome.startIndex]
// 取出最后一个字符
welcome[welcome.index(before: greeting.endIndex)]

```

