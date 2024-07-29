# Flutter setState流程

## setState

1. 先执行闭包函数`fn`，我们都知道`setState`包裹的会先执行再刷新
2. `StatefulElement`执行`markNeedsBuild`方法

```dart

  @protected
  void setState(VoidCallback fn) {
		....
    final Object? result = fn() as dynamic;
		.... 
    _element!.markNeedsBuild();
  }
```

## markNeedsBuild

1. 对当前树进行标脏，如果已经是标脏过则直接返回
2. 执行`BuildOwner`的`scheduleBuildFor`方法

```dart
void markNeedsBuild() {
...
if (!_active)

    return;

    if (dirty)

    return;
    ...
    _dirty = true;
    owner.scheduleBuildFor(this);
}
```

*查看*BuildOwner的源码可以发现这个是一个Element的管理类

```dart
class BuildOwner {
 	....
  VoidCallback? onBuildScheduled;

  final _InactiveElements _inactiveElements = _InactiveElements();

  final List<Element> _dirtyElements = <Element>[];
  bool _scheduledFlushDirtyElements = false;

  bool? _dirtyElementsNeedsResorting;
  
  void scheduleBuildFor(Element element) {....}
  
  void buildScope(Element context, [ VoidCallback? callback ]) {....}
  
  // 还有一些GlobalKey相关的
  void _registerGlobalKey(GlobalKey key, Element element) {....}
  void _unregisterGlobalKey(GlobalKey key, Element element) {...}
	....
}
```

## scheduleBuildFor

```dart
final List _dirtyElements = [];
......

void scheduleBuildFor(Element element) {
    ......
    if (element._inDirtyList) {
        _dirtyElementsNeedsResorting = true;
        return;
    }
    .....
    if (!_scheduledFlushDirtyElements && onBuildScheduled != null) {
            _scheduledFlushDirtyElements = true;
        // 是一个回调方法
        onBuildScheduled();
    }
    _dirtyElements.add(element);
    element._inDirtyList = true;
}
```

1. 将`element`添加到脏列表_dirtyElements中,并触发`WidgetsBinding`的回调。

   > `onBuildScheduled()`回调在`WidgetsBinding`类中执行

2. `onBuildScheduled()`回调在`WidgetsBinding`类中执行。

```dart
mixin WidgetsBinding on BindingBase, ServicesBinding, SchedulerBinding, GestureBinding, RendererBinding, SemanticsBinding {
    @override
    void initInstances() {
    		// 此方法会添加一个持久的帧回调。    
        super.initInstances();
        _instance = this;
        // 初始化_buildOwner对象
        _buildOwner = BuildOwner();
        // 配置onBuildScheduled回调
        buildOwner!.onBuildScheduled = _handleBuildScheduled;
        window.onLocaleChanged = handleLocaleChanged;
        window.onAccessibilityFeaturesChanged = handleAccessibilityFeaturesChanged;
        .....
    }

    void _handleBuildScheduled() {
        // 执行ensureVisualUpdate
        ensureVisualUpdate();
    }
}
```

## ensureVisualUpdate

1. 在`SchedulerBinding`中触发`ensureVisualUpdate`的方法，根据[不同的调度任务执行不同的操作。

```dart
void ensureVisualUpdate() {
    switch (schedulerPhase) {
    // 没有正在处理的帧，可能正在执行的是 WidgetsBinding.scheduleTask，
    // scheduleMicrotask，Timer，事件 handlers，或者其他回调等
    case SchedulerPhase.idle:
    // 主要是清理和计划执行下一帧的工作
    case SchedulerPhase.postFrameCallbacks:
        scheduleFrame();//请求新的帧渲染。
        return;
    // SchedulerBinding.handleBeginFrame 过程， 处理动画状态更新
    case SchedulerPhase.transientCallbacks:
    // 处理 transientCallbacks 阶段触发的微任务（Microtasks)
    case SchedulerPhase.midFrameMicrotasks:
    // WidgetsBinding.drawFrame 和 SchedulerBinding.handleDrawFrame 过程，
    // build/layout/paint 流水线工作
    case SchedulerPhase.persistentCallbacks:
        return;
    }
}
```

## scheduleFrame

1. 在`SchedulerBinding`请求新的frame，注册 Vsync 信号。

```dart
void scheduleFrame() {
    if (_hasScheduledFrame || !framesEnabled)
   	 	return;
    //确保帧渲染的回调已经被注册
    ensureFrameCallbacksRegistered();
    // window调度帧
    window.scheduleFrame();
    _hasScheduledFrame = true;
}

@protected
void ensureFrameCallbacksRegistered() {
    //处理渲染前的任务
  window.onBeginFrame ??= _handleBeginFrame;
    //核心渲染流程
  window.onDrawFrame ??= _handleDrawFrame;
}
```

* Flutter在window上注册一个`onBeginFrame`和一个`onDrawFrame`回调，在`onDrawFrame`回调中最终会调用`drawFrame`。
* 当我们调用`window.scheduleFrame()`方法之后，`Flutter`引擎会在合适的时机（可以认为是在屏幕下一次刷新之前，具体取决于`Flutter`引擎的实现）来调用`onBeginFrame`和`onDrawFrame`。

## onBeginFrame、onDrawFrame

1. 当下一次刷新之前，会调用`onBeginFrame`和`onDrawFrame`方法。

```dart
void _handleBeginFrame(Duration rawTimeStamp) {
    if (_warmUpFrame) {
        assert(!_rescheduleAfterWarmUpFrame);
        _rescheduleAfterWarmUpFrame = true;
        return;
    }
    handleBeginFrame(rawTimeStamp);
}


void handleBeginFrame(Duration? rawTimeStamp) {
    .....
    _hasScheduledFrame = false;

    try { _frameTimelineTask?.start('Animate', arguments: timelineArgumentsIndicatingLandmarkEvent);
    		/// 执行动画的回调方法 schedulerPhase.transientCallbacks;
        _schedulerPhase = SchedulerPhase.transientCallbacks;
        final Map callbacks = _transientCallbacks;
        _transientCallbacks = {};
        callbacks.forEach((int id, _FrameCallbackEntry callbackEntry) {
        if (!_removedIds.contains(id))
            _invokeFrameCallback(callbackEntry.callback, _currentFrameTimeStamp!, callbackEntry.debugStack);
        });
        _removedIds.clear();
    } finally {
        // 最后,  SchedulerPhase.midFrameMicrotasks;
        _schedulerPhase = SchedulerPhase.midFrameMicrotasks;
    }
}
```

* 主要就是遍历_transientCallbacks，执行相应的Animate操作， 可通过`scheduleFrameCallback()`/`cancelFrameCallbackWithId()`来完成添加和删除成员 最后将调度状态更新到`SchedulerPhase`.`midFrameMicrotasks`.

* handleDrawFrame():

```dart
void handleDrawFrame() {

    Timeline.finishSync(); // end the "Animate" phase
    try {
        // 遍历_persistentCallbacks，执行相应的回调方法
        // _persistentCallbacks 主要包括build、layout、draw流程
        _schedulerPhase = SchedulerPhase.persistentCallbacks;

        for (final FrameCallback callback in _persistentCallbacks)
            _invokeFrameCallback(callback, _currentFrameTimeStamp);
            _schedulerPhase = SchedulerPhase.postFrameCallbacks;
            final List localPostFrameCallbacks =
            List.from(_postFrameCallbacks);
          _postFrameCallbacks.clear();
        for (final FrameCallback callback in localPostFrameCallbacks)
        _invokeFrameCallback(callback, _currentFrameTimeStamp);

    } finally {
        _schedulerPhase = SchedulerPhase.idle;
        Timeline.finishSync(); // end the Frame
        _currentFrameTimeStamp = null;
    }
}
```

* `_postFrameCallbacks`可通过`addPersistentFrameCallback()`注册，一旦注册后不可移除，后续每一次`frame`回调都会执行,`handleDrawFrame()`执行完成后会清空`_postFrameCallbacks`内容`_postFrameCallbacks`主要是状态清理，准备调度下一帧`frame`绘制请求。_
* 当遍历`_persistentCallbacks`时会执行相应的回调方法也就是通过`addPersistentFrameCallback`添加注册的`handlePersistentFrameCallback()`方法。最后会执行`RendererBinding`的`drawFrame`方法

## 调用RendererBinding和drawFrame方法

```dart
@override

void drawFrame() {

    TimingsCallback firstFrameCallback;
  	
    if (_needToReportFirstFrame) {
        firstFrameCallback = (List timings) {
            if (!kReleaseMode) {
                developer.Timeline.instantSync('Rasterized first useful frame');
                developer.postEvent('Flutter.FirstFrame', {});
            }
            SchedulerBinding.instance.removeTimingsCallback(firstFrameCallback);
              firstFrameCallback = null;
               _firstFrameCompleter.complete();
        };
        SchedulerBinding.instance.addTimingsCallback(firstFrameCallback);
    }

    try {
        if (renderViewElement != null)
            buildOwner.buildScope(renderViewElement);
        // 调用RendererBinding的drawFrame方法
        super.drawFrame();
        buildOwner.finalizeTree();
    } finally {

    }

    if (!kReleaseMode) {
        if (_needToReportFirstFrame && sendFramesToEngine) {
        	 developer.Timeline.instantSync('Widgets built first useful frame');
        }
    }
  
    _needToReportFirstFrame = false;
    if (firstFrameCallback != null && !sendFramesToEngine) {
        SchedulerBinding.instance.removeTimingsCallback(firstFrameCallback);
    }
}
```

