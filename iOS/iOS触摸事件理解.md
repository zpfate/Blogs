## 如何确定由哪个视图响应触摸事件

解决这个问题之前，我们先看下几个与之相关的对象概念。

### UITouch

触摸，`UITouch`保存着跟手指相关的信息，比如触摸的位置、时间、阶段等。

1.  当手指移动时，系统会更新同一个`UITouch`对象，使之能够一直保存该手指在的触摸位置。
    当手指离开屏幕时，系统会销毁相应的UITouch对象。
2. 一个手指触摸一次屏幕就会生产一个`UITouch`对象，多个手指同时触摸屏幕，就会生成多个`UITouch`对象。
3. 单个手指多次触摸屏幕（类似于双击等操作），系统会根据触摸的位置判断是否更新同一个`UITouch`对象（`UITouch`的`tapCount`属性也会随之更新）

```objective-c
// 触摸产生时所处的窗口
@property(nonatomic,readonly,retain) UIWindow    *window;

// 触摸产生时所处的视图
@property(nonatomic,readonly,retain) UIView      *view;

// 短时间内点按屏幕的次数，可以根据tapCount判断单击、双击或更多的点击
@property(nonatomic,readonly) NSUInteger          tapCount;

// 记录了触摸事件产生或变化时的时间，单位是秒
@property(nonatomic,readonly) NSTimeInterval      timestamp;

// 当前触摸事件所处的状态
@property(nonatomic,readonly) UITouchPhase        phase;

// 返回值表示触摸在view上的位置
// 这里返回的位置是针对view的坐标系的（以view的左上角为原点(0, 0)）
// 调用时传入的view参数为nil的话，返回的是触摸点在UIWindow的位置
- (CGPoint)locationInView:(UIView *)view;

// 记录了前一个触摸点的位置
- (CGPoint)previousLocationInView:(UIView *)view;
```

### UIEvent

事件，记录事件产生的时刻和事件的类型。每产生一个事件，都会生成一个`UIEvent`对象。

```objective-c
// 事件类型
@property(nonatomic,readonly) UIEventType     type API_AVAILABLE(ios(3.0));
@property(nonatomic,readonly) UIEventSubtype  subtype API_AVAILABLE(ios(3.0));
// 事件产生的时间
@property(nonatomic,readonly) NSTimeInterval  timestamp;
```

### UIResponder（响应者对象）

着重介绍下`UIResponder`，在iOS中不是任何对象都能处理事件的，只有继承了`UIResponder`的对象才能接受并处理事件，称之为响应者对象。老生常谈，打开`UIResponder`类的官方文档：

```objective-c
// 响应链中下一个响应者
@property(nonatomic, readonly, nullable) UIResponder *nextResponder;
// 是否可以成为第一响应者
@property(nonatomic, readonly) BOOL canBecomeFirstResponder; 
- (BOOL)becomeFirstResponder;
// 是否可以放弃第一响应者
@property(nonatomic, readonly) BOOL canResignFirstResponder;    
// 放弃第一响应者
- (BOOL)resignFirstResponder;
// 是否是第一响应者
@property(nonatomic, readonly) BOOL isFirstResponder;

// 触摸事件
// 触摸开始
- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
// 触摸移动
- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
// 触摸结束
- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
// 触摸取消
- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
// 3D触摸事件
- (void)touchesEstimatedPropertiesUpdated:(NSSet<UITouch *> *)touches API_AVAILABLE(ios(9.1));

// 按压事件
- (void)pressesBegan:(NSSet<UIPress *> *)presses withEvent:(nullable UIPressesEvent *)event API_AVAILABLE(ios(9.0));
- (void)pressesChanged:(NSSet<UIPress *> *)presses withEvent:(nullable UIPressesEvent *)event API_AVAILABLE(ios(9.0));
- (void)pressesEnded:(NSSet<UIPress *> *)presses withEvent:(nullable UIPressesEvent *)event API_AVAILABLE(ios(9.0));
- (void)pressesCancelled:(NSSet<UIPress *> *)presses withEvent:(nullable UIPressesEvent *)event API_AVAILABLE(ios(9.0));

// 加速计事件
- (void)motionBegan:(UIEventSubtype)motion withEvent:(nullable UIEvent *)event API_AVAILABLE(ios(3.0));
- (void)motionEnded:(UIEventSubtype)motion withEvent:(nullable UIEvent *)event API_AVAILABLE(ios(3.0));
- (void)motionCancelled:(UIEventSubtype)motion withEvent:(nullable UIEvent *)event API_AVAILABLE(ios(3.0));

// 远程控制事件
- (void)remoteControlReceivedWithEvent:(nullable UIEvent *)event API_AVAILABLE(ios(4.0));

// 通过这个方法告诉UIMenuController它内部应该显示什么内容,”复制”、”粘贴”等
- (BOOL)canPerformAction:(SEL)action withSender:(nullable id)sender API_AVAILABLE(ios(3.0));

// 默认的实现是调用canPerformAction:withSender:方法来确定对象是否可以调用action操作。如果我们想要重写目标的选择方式，则应该重写这个方法。
- (nullable id)targetForAction:(SEL)action withSender:(nullable id)sender API_AVAILABLE(ios(7.0));

// 要求响应者从菜单中添加和删除项目。
- (void)buildMenuWithBuilder:(id<UIMenuBuilder>)builder API_AVAILABLE(ios(13.0));
// 要求响应者验证命令
- (void)validateCommand:(UICommand *)command API_AVAILABLE(ios(13.0));

// UIResponder提供了一个只读方法来获取响应链中共享的undo管理器，公共的事件撤销管理者
@property(nullable, nonatomic,readonly) NSUndoManager *undoManager API_AVAILABLE(ios(3.0));

// 编辑手势
@property (nonatomic, readonly) UIEditingInteractionConfiguration editingInteractionConfiguration API_AVAILABLE(ios(13.0));
```

日常接触到的继承自`UIResponder`的类有三个:

* AppDelegate
* UIApplication
* UIVIewController
* UIView

## 事件的产生和传递

前面介绍到这里，一次触摸，产生的事件，以及响应者都已经具备，唯一的问题就是如何确定哪个是我们想要的响应者，那么如何找到最合适的控件来处理事件。

应用接收到事件后，会先将其加入事件队列中以等待处理（队列的好处是先进先出）。

1. `UIApplication`会事件传递给当前显示的窗口（`UIWindow`），如果存在多个窗口，优先询问后显示的窗口。
2. 如果窗口可以响应事件，则传递给子视图；不能响应则将事件传递给其他窗口。
3. 若子视图能响应，则从后往前询问当前视图的子视图；否则将事件传递给上一个同级的子视图，重复该步骤。
4. 如果没有能响应的子视图，则自身就是最合适的响应者

```objectivec
UIApplication ——> UIWindow ——> 子视图 ——> ... ——> 子视图
```

日常开发中我们都知道想要响应事件必须满足下面的条件：

1. 触摸点在该控件范围内
2. 该控件能接收触摸事件

### Hit-Testing

以下几种状态的视图无法响应事件：

* **不允许交互**：userInteractionEnabled = NO

* **隐藏**：hidden = YES 如果父视图隐藏，那么子视图也会隐藏，隐藏的视图无法接收事件

* **透明度**：alpha < 0.01 如果设置一个视图的透明度<0.01，会直接影响子视图的透明度。alpha：0.0~0.01为透明。

#### hitTest:withEvent:

每个UIView对象都有一个 `hitTest:withEvent:` 方法，这个方法是`Hit-Testing`过程中最核心的存在，其作用是询问事件在当前视图中的响应者，同时又是作为事件传递的桥梁。大概实现的思路如下：

```objective-c
- (UIView *)hitTest:(CGPoint)point withEvent:(UIEvent *)event{
    // 无法响应事件
     if (self.userInteractionEnabled == NO || self.hidden == YES ||  self.alpha <= 0.01) return nil; 
    // 触摸点若不在当前视图上则无法响应事件
    if ([self pointInside:point withEvent:event] == NO) return nil; 
    // 倒序遍历 
    int count = (int)self.subviews.count; 
    for (int i = count - 1; i >= 0; i--) { 
        // 获取子视图
        UIView *childView = self.subviews[i]; 
        // 坐标系的转换， 把触摸点在当前视图上坐标转换为在子视图上的坐标
        CGPoint childPoint = [self convertPoint:point toView:childView]; 
        // 询问子视图层级中的最佳响应视图
        UIView *responder = [childView hitTest:childPoint withEvent:event]; 
        if (responder) {
            // 如果子视图中有更合适的就返回
            return responder; 
        }
    } 
    // 子视图中没有找到合适的，返回自身
    return self;
}
```

## 参考链接

[Using Responders and the Responder Chain to Handle Events](https://developer.apple.com/documentation/uikit/touches_presses_and_gestures/handling_touches_in_your_view)

[iOS触摸事件全家桶](https://www.jianshu.com/p/c294d1bd963d)