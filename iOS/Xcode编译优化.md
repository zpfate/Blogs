## 代码层面优化

1. 减少头文件引入 代替使用@class
2. 减少使用xib，storyboard
3. 打包framework或者静态库（组件化）
4. 常用的头文件放入到pch文件中，开启预编译

## 编译器选项优化

1.  提高XCode编译时使用的线程数

   ```
   defaults write com.apple.Xcode PBXNumberOfParallelBuildSubtasks 8  
   ```

2. 加载RAM磁盘编译Xcode项目

   1. 将DeriveData下的文件删除：

      ```shell
      rm -rf ~/Library/Developer/Xcode/DerivedData/*
      ```

   2. 在~/Library/Developer/Xcode/DerivedData.上部署安装2 GB大小的RAM磁盘

      ```shell
      cd ~/Library/Developer/Xcode/DerivedData
      # 创建2 GB的RAM磁盘（size的计算公式 size = 需要分配的空间(M) * 1024 * 1024 / 512）
      hdid -nomount ram://4194304
      # 初始化磁盘 有以下输出：
      # Initialized /dev/rdisk3 as a 2 GB case-insensitive HFS Plus volume
      newfs_hfs -v DerivedData /dev/rdiskN
      # 安装磁盘
      diskutil mount -mountPoint ~/Library/Developer/Xcode/DerivedData /dev/diskN
      
      ```

      





## 使用Injection



![image-20220420102527639](../../../../../Users/twistedfate/Library/Application%20Support/typora-user-images/image-20220420102527639.png)

first time

![image-20220420094316493](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1650419027.png)



![image-20220420094514002](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1650419114.png)

![image-20220420094616359](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1650419177.png)

cocoapods-pre

![image-20220420100442221](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1650420283.png)

![image-20220420100612613](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1650420373.png)

![image-20220420100737322](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1650420458.png)

