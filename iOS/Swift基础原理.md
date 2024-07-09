# Swift基础整理

## 闭包和Objective-C中Block的差别

### Swift闭包使用

```swift
// 后置闭包
var names = ["AT", "AE", "D", "S", "BE"]
names.sort(){$0 > $1}
print("\(names)")	
```

#### 逃逸闭包和非逃逸闭包的区别：

逃逸闭包不需要在函数结束前被调用，可以等到特定时机时才被调用。



### 与OC Block差别

Swift闭包都是捕获的是“引用”，而不是他们引用的对象，闭包就像是oc给外部变量默认添加了`__block`或者`__weak`。



## Swift枚举原理

### 简单实用

```swift
enum Direction {
    case north
    case south
    case east
    case west
}
```

等价写法：

```swift
enum Direction {
    case north, south, east, west
}
```

### RawValue

```swift
enum Season: String {
    case spring = "spring"
    case summer = "summer"
    case autumn = "autumn"
    case winter = "winter"
}

print(Season.spring.rawValue) // spring
print(Season.summer.rawValue) // summer
print(Season.autumn.rawValue) // autumn
print(Season.winter.rawValue) // winter
```

### 关联值

```swift
enum Season {
    case spring(month: Int)
    case summer(startMonth: Int, endMonth: Int)
}

var sping: Season = Season.spring(month: 1)
switch sping {
case .spring(let month):
    print(month)
case .summer(let startMonth, let endMonth):
    print(startMonth, endMonth)
}
```

**关联值和枚举值的区别：**

<font color=red>原始值RawValue不占用枚举的内存，是一个计算属性，而关联值占据枚举的内存。</font>



## 反射Reflection

```swift

class Person {
    let name: String?
    let age: Int?
    
    init(name: String?, age: Int?) {
        self.name = name
        self.age = age
    }
}

// main
let p = Person(name: "Zed", age: 222)
let mirror = Mirror(reflecting: p)
print("mirror.children: \(mirror.children)")

for child in mirror.children {
    
    if let key = child.label {
    print("\(child.label!): \(child.value)")
}
```



## 高阶函数

### map

map函数: 遍历数组 并可以对其中的元素做一次处理

```swift
var values = [1, 2, 3, 4, 5]
// map: [11, 12, 13, 14, 15]
print("map: \(values.map({ e in return e + 10}))")
```





### compactMap

compactMap函数: 

1. 实现降维操作,二维数组变为一维数组 

2. 强制解包,过滤掉nil元素

```swift
var strs = ["1", "2", "second", "3", "five", "4"]
// compactMap: [1, 2, 3, 4]
print("compactMap: \(strs.compactMap{Int($0)})")
```



### filter

filter函数: 过滤操作

```swift
var numbers = [1, 2, 3, 4, 5]
// filter: [3, 4, 5]
print("filter: \(numbers.filter{$0 > 2})")
```



### reduce

reduce函数: 可以设置一个初始值, 返回两个结果变量(result, current)

```swift
var prices = [10, 20, 30]
// reduce: 60
print("reduce: \(prices.reduce(0, { partialResult, value in return partialResult + value}))")
// reduce: 70
print("reduce: \(prices.reduce(10){$0 + $1})")
```



