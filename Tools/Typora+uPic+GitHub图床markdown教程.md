  今天给大家推荐一款写博客神器Typora，真的是非常好使的一款markdown编辑器。

### 优势：

1.及时预览，虽然现在很多网站都可以分屏显示，但是体验还是有点缺陷（偏移量不一致）

2.本地可查看

3.当然是完全免费的啦

下载地址：[官网下载](https://www.typora.io/)

下载完成之后，我们直接打开Typoro便可以进行编写，markdown的语法还是一样的。下面我们进行编写时候插入图片的设置，打开 偏好设置->图片，如下图：

![偏好设置 图像设置](https://user-gold-cdn.xitu.io/2020/4/7/17153bbdfd6ea979?w=700&h=480&f=png&s=51090)

上传服务的设定有好几种方案，

![上传图片服务设置](https://user-gold-cdn.xitu.io/2020/4/7/17153bbe1107726e?w=470&h=113&f=png&s=53970)

这边我因为是mac ox系统，选择的是uPic，这边介绍一下uPic的使用方式，windows用户的同学可以使用PicGo,无论是command line还是app方式都可以，方式应该都差不多。

在终端中输入：

```shell

```
brew cask install uPic
```

```

安装完成后，打开uPic程序，在系统偏好设置 扩展中进行设置：

![系统偏好设置 扩展](https://user-gold-cdn.xitu.io/2020/4/7/17153bbd9bff9e19?w=668&h=508&f=png&s=69729)

勾选访达扩展中的 uPic访达扩展，到了这里我们就可以使用图片上传了。

在右上角中选择uPic图标

![uPic使用](https://user-gold-cdn.xitu.io/2020/4/7/17153bbdabb34940?w=231&h=252&f=png&s=105781)

* 注意：截图上传的话需要给录制权限。

  因为此时上传使用的是默认的SMMS的图床，我们可以使用GitHub作为自己的图床，打开uPic的偏好设置

  ![uPic设置GitHub作为图床](https://user-gold-cdn.xitu.io/2020/4/7/17153bbd99830627?w=599&h=744&f=png&s=37894)

用户名就是你的GitHub用户名，再填上你的仓库名称，分支的话一定要使用已经创建的分支才可以上传成功，Token就是GitHub的Personal Access Token，[去这里设置](https://github.com/settings/tokens/new)，新建成功后复制粘贴到这里即可。

* 注意：Personal Access Token只会在创建的时候显示一次

* 保存路径自行添加，注意别出现重复名称

  到了这里，我们的自己的图床配置就已经完成了，我们可以点击下面的验证按钮，上传成功后可以在GitHub仓库中发现多了一张uPic logo的图片。

  ![uPic logo](https://user-gold-cdn.xitu.io/2020/4/7/17153bbd9b1c1973?w=980&h=172&f=png&s=2652)

到这里，大功告成~
如果打开uPic显示 `xxx.app可能会损坏您的电脑，您应该将他移入废纸篓`云云，解决方法如下：

> ### 1. macOS Mojave 10.14及以下系统：
>
> 打开「终端.app」，输入以下命令并回车，输入开机密码回车
>
> ```shell
> sudo xattr -rd com.apple.quarantine 空格 软件的路径
> ```
>
> ### 2. macOS Catalina 10.15系统：
>
> 打开「终端.app」，输入以下命令并回车，输入开机密码回车
>
> ```shell
> sudo xattr -rd com.apple.quarantine 空格 软件的路径
> ```
>
> ### 3. macOS Catalina 10.15.4 系统：
>
> ```
> sudo codesign --force --deep --sign - (应用路径)
> ```
>
> ***错误解决：***
>
> 如出现以下错误提示：
>
> ![错误提示](https://user-gold-cdn.xitu.io/2020/5/26/1724fce7e7d1247a?w=520&h=35&f=png&s=16993)

> 那么，先在终端执行：
>
> ```shell
> xattr -cr /文件位置（直接将应用拖进去即可）
> ```
>
> 然后再次执行如下指令即可：
>
> ```shell
> codesign --force --deep --sign - /文件位置（直接将应用拖进去即可）
> ```