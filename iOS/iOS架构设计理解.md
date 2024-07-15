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

