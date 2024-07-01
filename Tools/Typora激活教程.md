# Typora激活教程

[Typora下载链接](https://typoraio.cn/)

1. 打开LicenseIndex.xxxxxxxxxxx.chunk.js

![image-20240701171703871](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1719825426820)

接下来的路径看下面的shell命令行`Contents`开始一次点击

或者使用命令行

```shell
open /Applications/Typora.app/Contents/Resources/TypeMark/page-dist/static/js
```

打开该目录下的LicenseIndex.xxxxxxxxxxx.chunk.js，推荐使用`vscode`

2.  cmd+f 搜索 e.hasActivated

```javascript
//将
e.hasActivated=“true”=e.hasActivated
 
//改为
e.hasActivated = "true" == "true"	
```

到这里可以直接cmd+s保存退出使用了，但是会有激活弹框提示。

3. 关闭自动欢迎弹框

   还是上一步的那个文件, 使用**cmd+f 搜索 e.exports=n.p+"static/media/**

```js
//e.exports=n.p+"static/media/完整语句
e.exports=n.p+"static/media/icon.06a6aa23.png"}
 
//在png" 和 }间添加下面内容
window.onload = function () {setTimeout(() => {window.Setting.invoke("close");}, 1);};	
```

这里直接保存退出，可以正常使用了。

4. 关闭底部未激活弹框

   打开 **/Applications/Typora.app/Contents/Resources/zh-Hans.lproj**下的`Panel.strings`

   ```js
   //将 
   "UNREGISTERED" = "未激活";
    
   //改为
   "UNREGISTERED" = " ";
   ```

   