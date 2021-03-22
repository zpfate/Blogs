简单的设想一下，一次触摸事件，一根手指，一个响应手指触摸的对象。在iOS中，手指触摸屏幕是一个UITouch的对象，而响应手指触摸的对象则是UIResponder。

### UIResponder（响应者对象）

首先介绍下UIResponder，在iOS中不是任何对象都能处理事件的，只有继承了UIResponder的对象才能接受并处理事件，称之为响应者对象。老生常谈，打开UIResponder类的官方文档：

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
- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(nullable UIEvent *)event;
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

日常接触到的继承自UIResponder的类有三个:

* UIApplication
* UIViewController
* UIView

