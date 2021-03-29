## 使用

1. 打开下面的目录, 将工程内的Templates目录直接粘贴在该目录下

\```

~/Library/Developer/Xcode

```

2. 创建工程时选择下图选项

   ![image-20200629175409949](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_29_17_59_31.png)

3. 选择创建类文件时选择

   ![创建类](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_29_18_00_05.png "创建类文件")

## 效果

1. 选择`Single View App No Scene`创建的就是没有`SceneDelegate`等文件的空白工程

![image-20200630152200173](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_22_00.png "没有SceneDelegate相关文件的空白工程")

2. 选择`Custom Cocoa Touch Class`创建的类文件则自带了一些初始化方法以及`pragma mark`标签备注

   ![image-20200630152408712](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_24_08.png "创建类文件")

   

   ## 制作

   大概说一下粗略的方法，没太仔细研究，主要是行个方便，制作出统一的模板，代码也会更规范点。

   `Xcode`本身自带的模板在`/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/Library/Xcode/Templates/`

   目录中，分为`Project Templates`和`File Projects`分别对应创建工程和新建文件。

   可以直接将该目录复制到`~/Library/Developer/Xcode`下，然后在里面改就可以了，`Xcode`加载的时候自带去重了，比如创建工程经常选择的`Single View App`对应如下:

   *** No Scene是我加上去的，作为新的名字  ***

   ![image-20200630153807553](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_38_07.png "Single View App No Scene")

   

   修改创建文件时里面显示的东西则是在该文件中配置

   ![image-20200630154009170](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_40_09.png "TemplateInfo.plist文件")

   

![image-20200630154158766](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_41_59.png "TemplateInfo.plist具体介绍")

![image-20200630154437203](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_44_37.png "info.plist文件配置")

同样创建类文件的修改在`File Templates`中：

![Custom Cocoa Touch Class](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_46_06.png "Custom Cocoa Touch Class")

在类文件模板中添加代码，比如说我想要创建继承自`NSObjec`的类，在下图中的文件中添加你想要创建时自带的代码即可。

![在类文件模板添加代码](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_47_10.png "在类文件模板添加代码")

我这里只给他添加了延展

![NSObject添加延展](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_06_30_15_49_06.png "NSObject添加延展")

*** ___FILEBASENAMEASIDENTIFIER___ 会自动替换为你输入的类名***

