# iOS架构设计

## iOS的系统架构

**iOS的系统架构分为四层，由上到下一次为：可触摸层(Cocoa Touch layer)、媒体层(Media layer)、核心服务层(Core Services layer)、核心操作系统层(Core OS layer)。**1

- **可触摸层**：为应用程序开发提供了各种常用的框架并且大部分框架与界面有关，本质上来说它负责用户在iOS设备上的触摸交互操作。它包括以下组件：Multi-Touch Events、Core Motion、Camera、View Hierarchy、Localization、Alerts、Web Views、Image Picker、Multi-Touch Controls等。
- **媒体层**：通过它我们可以在应用程序中使用各种媒体文件，进行音频与视频的录制，图形的绘制，以及制作基础的动画效果。它包括以下组件：Core Audio、OpenGL、Audio Mixing、Audio Recording、Video Playback、JPG、PNG、TIFF、PDF、Quartz Core Animation、OpenGL ES等。
- **核心服务层**：我们可以通过它来访问iOS的一些服务。它包括以下组件：Collections、Address Book、Networking、File Access、SQLite、Core Location、Net Services、Threading、Preferences、URL Utilities等。
- **核心操作系统层**：包括内存管理、文件系统、电源管理以及一些其他的操作系统任务。它可以直接和硬件设备进行交互。核心操作系统层包括以下组件：OS X Kernel、Mach 3.0、BSD、Sockets、Power Mgmt、File System、Keychain、Certificates、Security、Bonjour等。

## 设计模式

### 工厂模式

定义了一个创建对象的抽象方法，由子类决定要实例化的类。工厂方法模式将对象的实例化推迟到子类。

比如说多种cell的工厂模式， 工厂类作为抽象类。

### 单例模式

确保一个类最多只有一个实例，并提供一个全局访问点。

可以使用懒加载或者预加载。

> iOS开发中常用的GCD创建单例的方式就是。

### 责任链模式

职责链模式的英文是Chain Of Responsibility Design Pattern。针对一个事务，让多个接收者都可以处理，将这些接收者对象串成一条链，并沿着这条链传递这个事务，直到链上的某个对象完成处理。

比如说Flutter Dio处理请求的拦截器。

### 代理模式

代理模式给某一个对象提供一个代理对象，并由代理对象控制对原对象的引用。通俗的来讲代理模式就是我们生活中常见的中介。



## 架构理解

### 好的架构

1. 代码整齐有规范，分类明确，没有`common`

   > `common`的问题是随着项目业务发展,架构的发展成长，会变得一团糟糕。

2. 思路和方法统一，不要多元

   > 解决问题的方案有很多，但是确认了一种方案，就不用在其他地方采用别的方案。

3. 没有横向依赖、减少跨层访问

   > 避免横向的依赖，有时候跨层访问是不可避免的，一旦出现类似的情况，尽快交给上层或者下层，尽量不要出现跨层的情况。
   >
   > 跨层访问会增加耦合度，否则需要整体替换的时候，牵涉面就会很大。

4. 便于测试以及拓展

   >  高度的模块化，尽可能减少依赖关系，便于mock

5. 保持一定的超前性

   > 关注行业动态，把握技术走向。保证适度的技术的超前性，能够使架构的更新变得轻松。

### 三层架构

* 数据管理者
* 数据加工者
* 数据展示者



## View层的架构设计

1. View代码结构的规定

   * 提高业务方View层的可读性和可维护性
   * 防止业务代码对架构产生腐蚀
   * 确保传承
   * 保持架构发展的方向不轻易被不合理的意见左右

   > #pragma mark - Life Cycle
   >
   > #pragma mark - Event Response
   >
   > #pragma mark - Delegate
   >
   > #pragma mark - Getter & Setter



## MVC

M应该做的事：

1. 给ViewController提供数据
2. 给ViewController存储数据提供接口
3. 提供经过抽象的业务基本组件，供Controller调度



C应该做的事：

1. 管理View Container的生命周期
2. 负责生成所有的View实例，并放入View Container
3. 监听来自View与业务有关的事件，通过与Model的合作，来完成对应事件的业务。



V应该做的事：

1. 响应与业务无关的事件，并因此引发动画效果，点击反馈（如果合适的话，尽量还是放在View去做）等。
2. 界面元素表达



## MVCS

基于MVC衍生出来的一套架构。从概念上来说，它拆分的部分是Model部分，拆出来一个Store。这个Store专门负责数据存取。但从实际操作的角度上讲，它拆开的是Controller。



### 胖model

`胖Model包含了部分弱业务逻辑`。胖Model要达到的目的是，`Controller从胖Model这里拿到数据之后，不用额外做操作或者只要做非常少的操作，就能够将数据直接应用在View上`。

胖Model相对比较难移植，虽然只是包含弱业务，但好歹也是业务，迁移的时候很容易拔出萝卜带出泥。另外一点，MVC的架构思想更加倾向于Model是一个Layer，而不是一个Object，不应该把一个Layer应该做的事情交给一个Object去做。最后一点，软件是会成长的，FatModel很有可能随着软件的成长越来越Fat，最终难以维护。

### 瘦model

`瘦Model只负责业务数据的表达，所有业务无论强弱一律扔到Controller`。瘦Model要达到的目的是，`尽一切可能去编写细粒度Model，然后配套各种helper类或方法来对弱业务做抽象，强业务依旧交给Controller`。

slimModel在一定程度上违背了DRY（Don't Repeat Yourself）的思路，Controller仍然不可避免在一定程度上出现代码膨胀。



## MVVM

MVVM本质上也是从MVC中派生出来的思想，MVVM着重想要解决的问题是尽可能地减少Controller的任务。不管MVVM也好，MVCS也好，他们的共识都是`Controller会随着软件的成长，变很大很难维护很难测试`。只不过两种架构思路的前提不同，MVCS是认为Controller做了一部分Model的事情，要把它拆出来变成Store，MVVM是认为Controller做了太多数据加工的事情，所以MVVM把`数据加工`的任务从`Controller`中解放了出来，使得`Controller`只需要专注于数据调配的工作，`ViewModel`则去负责数据加工并通过通知机制让View响应ViewModel的改变。

MVVM是`基于胖Model的架构思路建立的，然后在胖Model中拆出两部分：Model和ViewModel`。关于这个观点我要做一个额外解释：胖Model做的事情是先为Controller减负，然后由于Model变胖，再在此基础上拆出ViewModel，跟业界普遍认知的`MVVM本质上是为Controller减负`这个说法并不矛盾，因为胖Model做的事情也是为Controller减负。

另外，我前面说`MVVM把数据加工的任务从Controller中解放出来`，跟`MVVM拆分的是胖Model`也不矛盾。要做到解放Controller，首先你得有个胖Model，然后再把这个胖Model拆成Model和ViewModel。

不用ReactiveCocoa也能MVVM，用ReactiveCocoa能更好地体现MVVM的精髓。

ViewModel本质上算是Model层（因为是胖Model里面分出来的一部分），所以View并不适合直接持有ViewModel，那么View一旦产生数据了怎么办？扔信号扔给ViewModel，用谁扔？ReactiveCocoa。

在MVVM中使用ReactiveCocoa的第一个目的就是如上所说，View并不适合直接持有ViewModel。第二个目的就在于，ViewModel有可能并不是只服务于特定的一个View，使用更加松散的绑定关系能够降低ViewModel和View之间的耦合度。



`在MVC的基础上，把C拆出一个ViewModel专门负责数据处理的事情，就是MVVM。`然后，为了让View和ViewModel之间能够有比较松散的绑定关系，于是我们使用ReactiveCocoa，因为苹果本身并没有提供一个比较适合这种情况的绑定方法。iOS领域里KVO，Notification，block，delegate和target-action都可以用来做数据通信，从而来实现绑定，但都不如ReactiveCocoa提供的RACSignal来的优雅，如果不用ReactiveCocoa，绑定关系可能就做不到那么松散那么好，但并不影响它还是MVVM。





要做一个View层架构，主要就是从以下三方面入手：

1. 制定良好的规范
2. 选择好合适的模式（MVC、MVCS、MVVM、VIPER）
3. 根据业务情况针对ViewController做好拆分，提供一些小工具方便开发
