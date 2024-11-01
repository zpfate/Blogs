# SwiftUI属性包装器



![img](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1723531971808)



## @State

`@State`用于在当前视图中声明可变的状态属性，并且在属性值改变时自动更新视图。

```swift
struct StateView: View {

    @State private var count = 0
    
    var body: some View {
        
        HStack (alignment: .center, spacing: 40) {
            Button(action: {
                count += 1
            }, label: {
                Image(systemName: "plus.square.fill").resizable().frame(width: 40, height: 40)
            })
            
            Button(action: {
                count -= 1
            }, label: {
                Image(systemName: "minus.square.fill").resizable().frame(width: 40, height: 40)
            })
        }
        
        Text("Current count: ").padding().font(.largeTitle)
        
        Text("\(count)").font(Font.system(size: 60)).foregroundStyle(Color.blue)
    }
}
```

## @Binding

`@Binding`用于在视图之间传递和共享可读写的值。对使用`@Binding`的属性进行赋值，改变的不是属性，而是它的引用，可以向外传递这个改变。

使用`@State`可以实现对当前`View`视图的状态管理，但是如果需要将状态传递到子视图，可以使用`@Binding`进行实现双向绑定。

```swift
struct BindingView: View {
    
    @State private var isVisible = false
    
    var body: some View {
        
        VStack (alignment: .center) {
            Toggle(isOn: $isVisible) {
                Text("Show child view")
            }.padding()
                        
            if isVisible {
                BindingChildView(isVisible: $isVisible)
            } else {
                Rectangle().frame(height: 60).foregroundColor(Color.clear)
            }
        }
    }
}

struct BindingChildView: View {
    
    @Binding var isVisible: Bool
    
    var body: some View {
        
        Text("Child View")
        
        Button {
            isVisible = false
        } label: {
            Image(systemName: "eye.slash").resizable().frame(width: 40, height: 30)
        }
    }
}
```

## @Environment

`@Environment`是声明一个属性来获取特定环境变量的值，而不需手动传递。当环境变量的值发生变化时，相关的视图就会自动更新。

```swift
struct EnvironmentView: View {
      
    @Environment(\.colorScheme) private var colorScheme
    
    var body: some View {
        
        if colorScheme == .dark {
            Text("Dark Mode")
                .foregroundColor(.black)
                .background(Color.white)
                .padding()
        } else {
            Text("Light Mode")
                .foregroundColor(.white)
                .background(Color.black)
                .padding()
        }
    }
}
```

## @ObservableObject

`@Published`用于在`ObservableObject`类中声明属性的一个属性包装器，当这个属性的值发生改变时，会自动触发视图的更新。

当你需要一个数据模型或对象需要在不同视图之间共享，并且需要在数据改变时更新用户界面事，`@ObservableObject`是非常有用的。

```swift
class UserData: ObservableObject {
    @Published var name: String = ""
}

struct ObservedObjectView: View {
    
    @ObservedObject var user = UserData()
    
    var body: some View {
        
        TextField("Enter your name", text: $user.name).frame(height: 50).border(Color.gray, width: 1).padding()
        Text("Current Text is:").padding()
        Text("\(user.name)").font(.largeTitle).foregroundStyle(Color.red).padding()
    }
}
```

## @StateObject

`@StateObject`用于在视图中创建一个持久化的可观察对象，并在视图的生命周期内保持持久性。它类似于`@ObservedObject`，但在对象生命周期中只会创建一次。

```swift
class DataModel: ObservableObject {
  
    @Published var data: [String] = []
    var count: Int = 0
}

struct StateObjectView: View {
    
    @StateObject var dataModel = DataModel()
    
    private func add() {
        dataModel.count += 1
    }
    
    var body: some View {
        
        VStack {
            Button(action: {
                add()
                dataModel.data.append("New Item \(dataModel.count)")
            }, label: {
                Text("Add Item")
            })
        }
        
        ForEach(dataModel.data, id: \.self) { item in
            Text(item)
        }
    }
}
```

## @EnvironmentObject

`@EnvironmentObject`用于在整个应用程序中共享全局环境对象。

```swift
final class Person: ObservableObject {
    @Published var name = "Alex Taylor"
}

struct EnvironmentObjectView: View {
    
    var body: some View {
        
        VStack {
            let p = Person()
            MapView().environmentObject(p)
        }
    }
}

struct MapView: View {
    @EnvironmentObject var p: Person
    
    var body: some View {
        VStack {
            Text(p.name)
            Button("点我") {
                p.name = "12345"
            }
        }
    }
}
```

## @Environment与@EnvironmentObject区别：
> `@Environment`属性包装器只能获取已经存在于环境中的环境变量的值
>
> `@EnvironmentObject`属性包装器结合`ObservableObject`协议来创建和传递自定义的环境对象

## @StateObject和@ObservedObject区别:
> `@SateObject`针对引用类型设计，被View持有，当View更新时，实例不会被销毁，与State类似，使得View本身拥有数据，这使得`@StateObject`适用于那些需要持久状态且与视图密切相关的数据对象，比如页面导航、用户输入等。`@StateObject`告诉SwiftUI，当这个视图更新时，你希望保留这个对象的一个示例
>
> `@ObservedObject`只是作为View的数据依赖，不被View持有，View更新时`ObservedObject`对象可能会被销毁,适合数据在SwiftUI外部存储，把`@ObservedObject`包裹的数据作为视图的依赖，比如数据库中存储的数据，当SwiftUI视图“更新”时，实际发生的是创建并显示视图的新示例。
> 这意味着当您通过`@ObservableObject`声明视图模型时，您将获得数据对象Model的一个新示例
