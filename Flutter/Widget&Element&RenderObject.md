# Widget & Element & RenderObject

## Widget

`Widget`是`Flutter`开发框架中最基本的概念，是视图的配置信息。前端框架中的常见名词，比如视图`View`、控制器`View Controller`、活动`Activity`、应用`Application`、布局`Layout`等，在`Flutter`中都是`Widget`。

<font color=redAccent>**Flutter的核心设计思想便是”一切皆Widget“**。</font>

**`Widget`是不可变的。**当视图渲染的配置信息发生变化时，`Flutter`会重建`Widget`树，`Widget`只是一份轻量级的数据结构，重建的成本很低。而且因为`Widget`的不可变性，可以以较低的成本进行复用。所以在一个真实的渲染树中，可能存在不同的`Widget`对应一个渲染节点的情况。



**`Flutter`的视图开发是声明式的，其核心设计思想就是将视图和数据分离，这与`React`的设计思路完全一致。**

### StatelessWidget

**无状态组件**

当你所要构建的用户界面不随任何状态信息的变化而变化时，选择使用`StatelessWidget`。简而言之，就是当你的初始化参数完全满足效果展示的时候使用`StatelessWidget`。

![StatelessWidget](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1720079154750)

### StatefulWidget

**有状态组件**

![StatefulWidget](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1720072959834)

> `Widget`是不可变的，更新则意味着销毁 + 重建（build）。`StatelessWidget`是静态的，一旦创建则无需更新；而对于`StatefulWidget`来说，在`State`类中调用 `setState`方法更新数据，会触发视图的销毁和重建，也将间接地触发其每个子`Widget`的销毁和重建。

### State

`StatefulWidget`是以`State`类代理`Widget`构建的设计方式实现的。

`StatefulWidget`没有使用`build`方法来创建视图，而是通过`createState()`的方法创建一个`State`对象，由`State`对象负责 视图的构建，`State`持有并处理`StatefulWidget`中的变化，一旦发生变化，就可以调用`setState`方法通知`Flutter`框架进行重构。

## Element

`Element`是`Widget`的一个实例化对象，承载了视图构建的上下文数据，是连接结构化的配置信息到完成最终渲染的桥梁。

### Flutter渲染过程

1. 通过`Widget`树生成对应的`Element`树
2. 创建相应的`RenderObject`并关联到`Element.renderObject`属性上
3. 构建成`RenderObject`树，完成最终的渲染。

```dart
abstract class Element extends DiagnosticableTree implements BuildContext {
	// ······
	@override
  Widget get widget => _widget!;
  // ······
  RenderObject? get renderObject {
  Element? current = this;
    while (current != null) {
      if (current._lifecycleState == _ElementLifecycle.defunct) {
        break;
      } else if (current is RenderObjectElement) {
        return current.renderObject;
      } else {
        current = current.renderObjectAttachingChild;
      }
    }
    return null;
	}
  // ·······
}
```

1. **可以看到`Element`同时持有`Widget`和`RenderObject`**

2. **`BuildContext`就是`Element`**

### Element树的意义

为什么不直接用`Widget`命令`RenderObject`去渲染，因为会极大的增加渲染带来的性能消耗。

因为 `Widget` 具有不可变性，但`Element` 却是可变的。实际上，`Element`树这一层将`Widget`树的变化（类似 React 虚拟 DOM diff）做了抽象，可以只将真正需要修改的部分同步到真实的`RenderObject`树中，最大程度降低对真实渲染视图的修改，提高渲染效率，而不是销毁整个渲染视图树重建。

## RenderObject

`RenderObject`是主要负责实现视图渲染的对象。

`Flutter` 通过控件树（`Widget`树）中的每个控件（`Widget`）创建不同类型的渲染对象，组成渲染对象树。而渲染对象树在`Flutter`的展示过程分为四个阶段，即布局、绘制、合成和渲染。 其中，布局和绘制在`RenderObject`中完成，`Flutter`采用深度优先机制遍历渲染对象树，确定树中各个对象的位置和尺寸，并把它们绘制到不同的图层上。绘制完毕后，合成和渲染的工作则交给`Skia`搞定。

### RenderObjectWidget

```dart
abstract class RenderObjectWidget extends Widget {
  @override
  RenderObjectElement createElement();
  @protected
  RenderObject createRenderObject(BuildContext context);
  @protected
  void updateRenderObject(BuildContext context, covariant RenderObject renderObject) { }
  ...
}
```

`RenderObjectWidget`是一个抽象类。我们通过源码可以看到，这个类中同时拥有创建`Element`、`RenderObject`，以及更新`RenderObject`的方法。

对于`Element`的创建，`Flutter`会在遍历`Widget`树时，调用`createElement`去同步`Widget`自身配置，从而生成对应节点的`Element`对象。而对于`RenderObject`的创建与更新，其实是在`RenderObjectElement`类中完成的。

```dart
abstract class RenderObjectElement extends Element {
  RenderObject _renderObject;

  @override
  void mount(Element parent, dynamic newSlot) {
    super.mount(parent, newSlot);
    _renderObject = widget.createRenderObject(this);
    attachRenderObject(newSlot);
    _dirty = false;
  }
   
  @override
  void update(covariant RenderObjectWidget newWidget) {
    super.update(newWidget);
    widget.updateRenderObject(this, renderObject);
    _dirty = false;
  }
  ...
}
```

在`Element`创建完毕后，`Flutter`会调用`Element`的`mount`方法。在这个方法里，会完成与之关联的`RenderObject`对象的创建，以及与渲染树的插入工作，插入到渲染树后的`Element`就可以显示到屏幕中了。



如果`Widget`的配置数据发生了改变，那么持有该`Widget`的`Element`节点也会被标记为`dirty`。在下一个周期的绘制时，`Flutter`就会触发`Element`树的更新，并使用最新的`Widget`数据更新自身以及关联的`RenderObject`对象，接下来便会进入`Layout`和`Paint`的流程。而真正的绘制和布局过程，则完全交由`RenderObject`完成：

```dart
abstract class RenderObject extends AbstractNode with DiagnosticableTreeMixin implements HitTestTarget {
  ...
  void layout(Constraints constraints, { bool parentUsesSize = false }) {...}
  
  void paint(PaintingContext context, Offset offset) { }
}
```

## 总结

`Widget`是`Flutter`世界里对视图的一种结构化描述，里面存储的是有关视图渲染的配置信息，`Element`则是`Widget`的一个实例化对象，将`Widget`树的变化做了抽象，能够做到只将真正需要修改的部分同步到真实的`RenderObject`树中，最大程度地优化了从结构化的配置信息到完成最终渲染的过程；而`RenderObject`，则负责实现视图的最终呈现，通过布局、绘制完成界面的展示。
