今天，给大家带来一篇关于Provider的底层分析。

我们都知道Provider是一个常用的状态管理的框架，那么什么是状态管理？

## 状态管理

简单明了的说，状态管理就是在某这个状态发生改变时，通知使用该状态的监听者发生改变，有一个联动的效果。

然后我们在讲解Provider之前，先了解一下想要实现这样的效果，需要了解的东西（InheritedWidget和ChangeNotifier）。

对这两已经明了的请直戳

### InheritedWidget

InheritedWidget是Flutter提供的非常重要的组件，提供了一种数据在widget树中从上到下传递、共享的方式，比如我们在应用的根widget中通过InheritedWidget共享了一个数据，那么我们便可以在任意子widget中来获取该共享的数据。

我们来通过一个计数器的Demo理解InheritedWidget的使用：

1. 先声明一个继承自InheritedWidget的widget

   ```dart
   class CounterInheritedWidget extends InheritedWidget {
     /// 需要共享的数据
     final int count;
   
     const CounterInheritedWidget(
         {Key? key, required this.count, required Widget child})
         : super(key: key, child: child);
   
     /// 默认的约定:如果状态是希望暴露的 应该提供一个静态of方法来获取对象
     /// 返回实例对象 方便子树中的widget获取共享数据
     static CounterInheritedWidget? of(BuildContext context) {
       return context.dependOnInheritedWidgetOfExactType<CounterInheritedWidget>();
     }
   
     /// 是否通知widget树中依赖该共享数据的子widget
     /// 这里当count发生变化时, 是否通知子树中所有依赖count的Widget重新build
     @override
     bool updateShouldNotify(covariant CounterInheritedWidget oldWidget) {
       return count != oldWidget.count;
     }
   }
   ```

2. 来两个监听父widget状态的子组件，分别继承自StatelessWidget和StatefulWidget

   ![继承自StatelessWidget的子Widget](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666690661.png "继承自StatelessWidget的子Widget")![继承自StatefulWidget的子Widget](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666690689.png "继承自StatefulWidget的子Widget")

3. 主页

   ![image-20221025173635786](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666690614.png)

最后点击几下按钮，效果如下

![InheritedWidget Demo效果](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666690890.png "InheritedWidget Demo效果")

## ChangeNotifier

ChangeNotifier 大概是 Flutter 中实现观察者模式最典型的例子了，它实现自 Listenable，内部维护一个 _listeners 列表用来存放观察者，并实现了 addListener、removeListener 等方法来完成其内部的订阅机制。

大概使用如下:

![继承自ChangeNotifier](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666747974.png)

![image-20221026093332062](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666748012.png)

## Provider

理解Provider底层实现的原理，我们从ChangeNotifierProvider这个来切入。

![image-20221026095214492](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1666749134.png)



