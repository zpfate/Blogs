## 启动流程

![三包项目启动流程&&业务模块](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_10_29_10_33_24.png "三包项目启动流程&&业务模块")

## 项目结构

1. 图片资源置于项目根目录的`assets`文件夹中

   图片命名规范，模块名\_页面名\_标识，比如：`account_edit_save`

2. 源代码置于项目根目录下的`lib`文件夹中

`lib`文件夹细分如下：

* `base`：基础通用组件

* `common`：常量全局配置

* `pages`：页面模块代码，文件夹内按业务模块分文件夹，model类置于model文件夹中

  *大致如上面流程图所示，下面八大模块外加登录 侧边栏*

* `utils`：工具类文件夹（注入插件封装等）

* `views`：轮子

## 插件引用

使用插件：

* 网络请求： dio
* 项目信息配置：package_info
* 隐私权限类:  permission_handler
* 浏览webview: webview_flutter
* 图片选择器：image_picker
* 本地持久化：shared_preferences
* url跳转插件：url_launcher
* 业务逻辑分离组件：flutter_bloc/
* ......







