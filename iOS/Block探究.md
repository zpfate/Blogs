## 什么是block

Block是将函数极其执行上下文封装起来的对象。

![image-20210315163225613](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_15_16_32_26.png)

通过`clang -rewrite-objc`命令编译该.m文件

![clang -rewrite-objc命令](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_15_16_36_24.png "clang -rewrite-objc命令")

发现block被编译成如下所示：

![block通过clang-rewrite-objc编译](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_15_16_40_31.png "block通过clang-rewrite-objc编译")

