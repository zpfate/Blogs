## Info.plist文件

### 介绍

info.plist文件主要描述的是一些工程的配置。

### 位置以及打开

![info.plist文件位置](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_26_17_04_47.png "info.plist文件位置")

另一种打开方式Open As Source Code：

![info.plist文件open as source code](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_26_17_06_47.png "info.plist文件open as source code")

*git操作冲突的时候打不开plist文件的时候可以使用该方法解决冲突，或者在Finder中右击选择打开方式，选择其他诸如VSCode之类的编辑器打开。*

或者直接在Project配置中查看：

![在Project配置中查看info.plist文件](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_26_17_05_36.png "在Project配置中查看info.plist文件")

### 具体配置

#### 项目配置相关

`Bundle identifier` ： App的Bundle I在Apple Developer中申请的

`Bundle name`:  包名

`Bundle version string (short)` ： 版本号

`Bundle version `：build号

`Launch screen interface file base name` ： 启动页

`Main storyboard file base name `： 启动的根视图storyboard

`Supported interface orientations `： 设备支持的方向

上面这些我们一般都不会在info.plist文件中修改，如果需要修改的话，直接在Project设置的General选项中修改，具体可以看一下下图：

![Project General配置](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_27_10_05_43.png "Project General配置")

#### 权限配置相关

隐私是一个大家越来越关注的问题，所以我们在app开发中经常需要配置隐私权限，才能使用对应的API。

我们可以在info.plist文件中点击"+"来添加对应的权限字段

![在info.plist文件中添加字段](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_27_10_17_13.png "在info.plist文件中添加字段")

一些常用的权限字段说明：

`Privacy - Calendars Usage Description` ：日历使用权限

`Privacy - Bluetooth Peripheral Usage Description` : 蓝牙使用权限

`Privacy - Health Share Usage Description` ： 健康分享权限

`Privacy - Health Update Usage Description` : 健康数据更新权限

`Privacy - Siri Usage Description` : Siri使用权限

`Privacy - Face ID Usage Description` :  Face ID使用权限

`Privacy - Microphone Usage Description` ： 麦克风使用权限，录音，发送语音

`Privacy - Camera Usage Description`： 摄像头使用权限，拍照，录制等

`Privacy - Photo Library Usage Description` ：相册使用权限（iOS11之前读写权限，iOS11之后只有读的权限）

`Privacy - Photo Library Additions Usage Description` ： 保存图片到相册（iOS11新增写入相册的权限）

`Privacy - Contacts Usage Description`  ：通讯录使用权限

`Privacy - Location Always Usage Description` ： 永久使用地址位置信息，定位相关（iOS11之前前后台都能获取定位的权限）

`Privacy - Location When In Use Usage Description` ：仅在App使用期间访问地理位置信息

`Privacy - Location Always and When In Use Usage Description` ：永久使用App定位的权限（iOS11之后）

* *位置权限配置的不同调用API也会有所区别，同时在调用API时的弹窗也会有所区分*
* *在配置相关权限的字段时，value中一定要写清楚具体的用途，而不是单纯的使用xx权限，否则会面临审核被拒的问题*

##### iOS14新增的权限相关

PHPhotoLibraryPreventAutomaticLimitedAccessAlert：

iOS14 中当用户选择 `PHAuthorizationStatusLimited` 时，如果未进行适配，有可能会在每次触发相册功能时都进行弹窗询问用户是否需要修改照片权限。在info.plist文件中设置`PHPhotoLibraryPreventAutomaticLimitedAccessAlert`为YES可以阻止该弹窗反复弹出。

#### ATS

由于安全原因，系统会拦截http请求，解决方案便是在info.plist文件中添加`App Transport Security Settings`的`Allow Arbitrary Loads`属性为YES。

![ATS设置](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_27_11_02_50.png "ATS设置")

#### 跳转白名单

iOS之后新增加的App间跳转的白名单设置`LSApplicationQueriesSchemes`

![App跳转白名单](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_11_27_11_11_59.png "App跳转白名单")