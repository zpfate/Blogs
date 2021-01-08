### 源码下载

在[Apple Open Source](https://opensource.apple.com/)上选择对应自己的macOS版本下，下载objc4即可。

笔者是macOS Catalina 10.15.6，所以选择的是[objc4-787.1](https://opensource.apple.com/tarballs/objc4/objc4-787.1.tar.gz)版本，顺便在该页面中下载Libc，dyld，libauto，libclosure，libdispatch，libpthread，xnu，libplatform这几个依赖库，编译objc4的时候需要用到。

### 解决编译错误

下载完成，解压之后我们打开工程文件，运行报错`unable to find sdk 'macosx.internal'`

![unable to find sdk 'macosx.internal'](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_09_49_49.png "unable to find sdk 'macosx.internal'")

解决方法：Project->Target->Build Settings->Base SDK选择macOS

![解决unable to find sdk 'macosx.internal'](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_10_04_56.png "解决unable to find sdk 'macosx.internal'")

再次运行报错`'sys/reason.h' file not found`

!['sys/reason.h' file not found](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_10_15_55.png "'sys/reason.h' file not found")

在项目目录下新建include/sys文件夹，将`xnu-6153.141.1/bsd/sys/reason.h`添加到该目录下再次编译。

这次报错是还是在`objc-os.h`文件中，提示`'mach-o/dyld_priv.h' file not found`!['mach-o/dyld_priv.h' file not found](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_10_45_35.png "'mach-o/dyld_priv.h' file not found")

解决方法和上面一样，该文件在`dyld-750.6/include/mach-o/dyld_priv.h`路径。

运行后提示`'os/lock_private.h' file not found`，该文件在`libplatform-220.100.1/private/os/lock_private.h`路径，之后报错`'os/base_private.h' file not found`,也在该目录下。

再次运行后报错`'pthread/tsd_private.h' file not found`，该文件在`libpthread-416.100.3/private/tsd_private.h`

再次运行后报错`'System/machine/cpu_capabilities.h' file not found`， 该文件在 `xnu-6153.141.1/osfmk/machine/cpu_capabilities.h`路径。

之后是报错`'os/tsd.h' file not found`，该文件在`xnu-6153.141.1/libsyscall/os/tsd.h`路径。

报错`'pthread/spinlock_private.h' file not found`,该文件在 `libpthread-416.100.3/private/spinlock_private.h`路径。

报错`'System/pthread_machdep.h' file not found`,该文件在只能在8.x级以前的版本中找到，[LibC下载](https://opensource.apple.com/source/Libc/)。

路径为 `Libc-825.26/pthreads/pthread_machdep.h) `。该文件导入后会发现重复声明的报错以及宏定义重复的警告，报错为:

`Typedef redefinition with different types ('int' vs 'volatile OSSpinLock' (aka 'volatile int'))`

`Static declaration of '_pthread_has_direct_tsd' follows non-static declaration`

`Static declaration of '_pthread_getspecific_direct' follows non-static declaration`

`Static declaration of '_pthread_setspecific_direct' follows non-static declaration`

![注释掉重复声明](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_15_16_52.png)

![注释重复声明2](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_15_18_44.png)

![注释重复声明3](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_15_19_31.png)

![注释重复宏1](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_15_19_58.png)

![注释重复宏3](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_07_15_20_24.png)

再次运行报错`'CrashReporterClient.h' file not found`，该文件在`Libc-825.26/include/CrashReporterClient.h`

报错`'_simple.h' file not found`，该文件在`libplatform-220.100.1/private/_simple.h`路径。

报错`'objc-shared-cache.h' file not found`,该文件在`dyld-750.6/include/objc-shared-cache.h`路径。

报错`Use of undeclared identifier 'DYLD_MACOSX_VERSION_10_11'`,在`dyld_priv.h`文件中添加下面的宏定义

```c
#define DYLD_MACOSX_VERSION_10_11 0x000A0B00
#define DYLD_MACOSX_VERSION_10_12 0x000A0C00
#define DYLD_MACOSX_VERSION_10_13 0x000A0D00
#define DYLD_MACOSX_VERSION_10_14 0x000A0E00
```



报错`Static_assert failed due to requirement 'bucketsMask >= ((unsigned long long)140737454800896ULL)' "Bucket field doesn't have enough bits for arbitrary pointers."`， 直接注释掉该断言。

报错`Use of undeclared identifier 'CRGetCrashLogMessage'`,解决方法在Project->Target->Build Settings->Preprocessor Macros加入LIBC_NO_LIBCRASHREPORTERCLIENT

![设置Preprocessor Macros](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_08_10_10_06.png "设置Preprocessor Macros")

报错`Mismatch in debug-ness macros`,解决方法是将改行代码注释掉

报错`'kern/restartable.h' file not found`, 该文件在` xnu-6153.141.1/osfmk/kern/restartable.h `路径。

报错`'Block_private.h' file not found`,该文件在`libclosure-74/Block_private.h`路径。

报错`can't open order file: /Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX11.1.sdk/AppleInternal/OrderFiles/libobjc.order`，解决方法是Project->Target->Build Settings->Order File添加路径`$(SRCROOT)/libobjc.order`。

报错`library not found for -lCrashReporterClient`,解决方法是在Project->Target->Build Settings->Other Linker Flags中删除`lCrashReporterClient`。

报错`'_static_assert' declared as an array with a negative size`，解决方法是直接注释掉该断言。

报错xcodebuild:1:1: `SDK "macosx.internal" cannot be located.`，解决方法
修改`Project->Target->Build Settings->Run Script(markgc)` 中的`macosx.internal`,改为`macosx`。

到了这里我们发现已经可以BUILD SUCCESSED了。

### 新建调试

我们新建一个target：

![新建target](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_08_13_39_48.png "新建target")

再为新建的Target添加好依赖

![为新建的target添加依赖](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_08_13_52_04.png "为新建的target添加依赖")

这时候我们创建一个Person类，在main.m中导入头文件

![导入Person类](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_08_14_00_55.png "导入Person类")

这时候我们发现断点无法使用，解决方案：

![移动main.m文件](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_08_14_03_37.png "移动main.m文件")

这时候再次运行会发现虽然断点生效了,但是无法进入alloc方法的实现，解决方案：

![设置Enable Hardened Runtime](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_01_08_14_07_04.png "设置Enable Hardened Runtime")

到这里就大功告成了。[objc4-787.1调试demo](https://github.com/ZpFate/LessonRuntime)