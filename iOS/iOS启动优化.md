### 启动流程

main函数之前：

* 加载dyld
* 创建启动闭包（更新App/重启手机需要）
* 加载动态库
* Bind & Rebase & Runtime初始化
* +load和静态初始化

### 优化思路

1. 删除启动项
2. 延迟初始化
3. 利用多线程
4. 更快的执行启动任务

### 优化方案

#### main之前

1. 减少动态库

   少动态库数量可以加减少启动闭包创建和加载动态库阶段的耗时，官方建议动态库数量小于 6 个。

   推荐的是动态库转静态库，还可以额外减少包大小。**不要链接那些用不到的库**（包括系统）

2. 下线代码

   下线代码可以减少 Rebase & Bind & Runtime 初始化的耗时。使用AppCode检索

   另外的一种静态扫描是基于 Mach-O 的：

   - `_objc_selrefs` 和`_objc_classrefs` 存储了引用到的 sel 和 class
   - `__objc_classlist` 存储了所有的 sel 和 class

   二者做个差集就知道那些类/sel 用不到，但**objc 支持运行时调用，删除之前还要在二次确认**。

3. +load迁移，静态初始化迁移

   **+load 除了方法本身的耗时，还会引起大量 Page In**，另外 +load 的存在对 App 稳定性也是冲击，因为 Crash 了捕获不到。

   静态初始化和 +load 方法一样也会引起大量 Page In，一般来自 C++代码，比如网络或者特效的库。另外**有些静态初始化是通过头文件引入进来的**，可以通过预处理来确认


#### main之后

延迟三方SDK初始化

高频次方法，虽然单个方法不耗时，可以做缓存

锁

线程数量： 线程的数量和优先级都会影响启动时间。可以通过设置 QoS 来配置优先级，两个高优的 QoS 是 User Interactive/Initiated，启动的时候，**需要主线程等待的子线程任务都应该设置成高优的**。

**高优的线程数量不应该多于 CPU 核心数量**，可以通过 System Trace 的 System Load 来分析这种情况。


图片

启动难免会用到很多图，有没有办法优化图片加载的耗时呢？

**用 Asset 管理图片而不是直接放在 bundle 里**。Asset 会在编译期做优化，让加载的时候更快，此外在 Asset 中加载图片是要比 Bundle 快的，因为 UIImage imageNamed 要遍历 Bundle 才能找到图。**加载 Asset 中图的耗时主要在在第一次张图，因为要建立索引**，可以通过把启动的图放到一个小的 Asset 里来减少这部分耗时。

每次创建 UIImage 都需要 IO，在首帧渲染的时候会解码。所以可以通过提前子线程预加载（创建 UIImage）来优化这部分耗时。

如下图，启动只有到了比较晚的阶段“RootWindow 创建”和“首帧渲染”才会用到图片，**所以可以在启动的早期开预加载的子线程启动任务**。

### 首帧渲染

不同 App 的业务形态不同，首帧渲染优化方式也相差的比较多，几个常见的优化点：

- LottieView：lottie 是 airbnb 用来做 AE 动画的库，但是加载动画的 json 和读图是比较慢的，可以**先显示一帧静态图，启动结束后再开始动画，或者子线程预先把图和 json 设置到 lottie cache 里**
- Lazy 初始化 View：**不要先创建设置成 hidden，这是很不好的习惯**
- AutoLayout：AutoLayout 的耗时也是比较高的，但这块往往历史包袱比较重，可以**评估 ROI 看看要不要改成 frame**
- Loading 动画：App 一般都会有个 loading 动画表示加载中，这个**动画最好不要用 gif**，线下测量一个 60 帧的 gif 加载耗时接近 70ms





