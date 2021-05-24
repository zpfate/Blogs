##### 前情摘要

Flutter 2.0请使用GetX 4.0.0以及之后的版本，未升级Flutter 2.0的同学请使用GetX 4.0.0以前的版本。

[GetX GitHub地址](https://github.com/jonataslaw/getx)

## GetX使用

### GetX 计数器Demo

我们先通过一个计数器的Demo来学会GetX的简单使用。

#### 1.添加依赖

在`pubspec.yaml`文件中添加依赖

```yaml
  get: ^4.0.0
```

#### 2.导入头文件

在需要用到的文件中导入

```dart
import 'package:get/get.dart';
```

##### 3.使用GetMaterialApp

```dart
class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return GetMaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: Homepage(),
    );
  }
}
```

* GetMaterialApp的子组件是MaterialApp
* GetMaterialApp对于路由、snackbar、国际化、bottomSheet、对话框以及与路由相关的高级apis和没有上下文（context）的情况下是必要的。
* 如果你只用Get来进行状态管理或依赖管理，就没有必要使用GetMaterialApp

##### 3.创建业务逻辑类

使用简单的`.obs`可以使任何变量成为可观察的。

```dart
import 'package:get/get.dart';
class CountController extends GetxController {
  var count = 0.obs;
  increment() {
    count++;
  }
}
```

##### 4.创建页面

可以使用`StatelessWidget`来创建页面，节省内存，使用Get你可能不需要使用`StatefulWidget`。

HomePage页面代码：

```dart
import 'package:flutter/material.dart';
import 'package:getx_demo/count_controller.dart';
import 'package:get/get.dart';
import 'package:getx_demo/next_page.dart';
class Homepage extends StatelessWidget {
  const Homepage({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    /// 使用Get.put()实例化你的类，使其对当下的所有子路由可用。
    final CountController controller = Get.put(CountController());
    
    return Scaffold(
      appBar: AppBar(
        title: Text("GetX Counter Demo"),
      ),
      body: Center(
        child: Column(
          children: [
            ElevatedButton(onPressed: ()=> Get.to(NextPage()), child: Text("Get to NextPage")),
            Obx(()=> Text("${controller.count}")),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        child: Icon(Icons.add),
        onPressed: controller.increment,
      ),
    );
  }
}
```

NextPage页面代码：

```dart
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:getx_demo/count_controller.dart';

class NextPage extends StatelessWidget {
  /// Get找到一个正在被其他页面使用的Controller
  final CountController controller = Get.find();
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text("Next Page"),),
      body: Center(
        child: Text("${controller.count}"),
      )
    );
  }
}
```

![计数器Demo效果](/Users/twistedfate/Library/Application Support/typora-user-images/image-20210524214148004.png)

