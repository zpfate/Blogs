## 为什么有这个指南

当然是因为Flutter项目运行起来太难了~ 尤其针对只有单平台原生能力甚至无移动端方向能力的开发者

## 前置要求

电脑上Flutter环境已经配置好，简单点说就是执行`flutter doctor`命令已经成功了，如果尚未配置好Flutter环境请左转 [Flutter中文网](https://flutterchina.club/get-started/install/)，或者本人博客

下面直接步入正题：

## 常见错误

### /Users/{用户名}/flutter/packages/flutter_tools/bin/xcode_backend.sh: No such file or directory

这是一个最常见的，基本切换分支或者拉了别人的代码经常遇到的一个问题，原因是每个人的Flutter安装环境不同

解决方案：

        在Android Studio中全局搜索`FLUTTER_ROOT`


​       ![全局搜索flutter_root](https://user-gold-cdn.xitu.io/2020/3/23/1710522e22e63643?w=605&h=269&f=png&s=48770 "全局搜索flutter_root")

找出非本地环境的那一行改之即可。

### 









