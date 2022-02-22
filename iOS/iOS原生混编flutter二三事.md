### 背景

一直以来，公司项目都是原生和原生和flutter并行发展，这一次技术团队想整点不一样的东西（需求紧张，解决人生不足），于是决定采用在原生种混编flutter的方案。

### 遇到的问题



1. 原生如何跳转到flutter端不同页面

   首先这边因为iOS只能初始化一个flutterEngine，一个engine只有一个flutterViewController，想要打开不同的页面，我们这边采用的是flutter向原生通信解决的。

2. 第一次跳转flutter页面，有点慢设置有点卡顿

   初始化一个FlutterViewController的子类，设置好背景色

   

   

   

   

   

   

