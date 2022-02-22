## Block相关

![img](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_16_09_53_59.jpg)

1. isa 指针，所有对象都有该指针，用于实现对象相关的功能。
2. flags，用于按 bit 位表示一些 block 的附加信息，本文后面介绍 block copy 的实现代码可以看到对该变量的使用。
3. reserved，保留变量。
4. invoke，函数指针，指向具体的 block 实现的函数调用地址。
5. descriptor， 表示该 block 的附加描述信息，主要是 size 大小，以及 copy 和 dispose 函数的指针。
6. variables，capture 过来的变量，block 能够访问它外部的局部变量，就是因为将这些变量（或变量的地址）复制到了结构体中。

### 三种block

1. _NSConcreteGlobalBlock 全局的静态 block，不会访问任何外部变量。
2. _NSConcreteStackBlock 保存在栈中的 block，当函数返回时会被销毁。
3. _NSConcreteMallocBlock 保存在堆中的 block，当引用计数为 0 时会被销毁。

**LLVM在ARC下NSConcreteStackBlock 的 block 会被 NSConcreteMallocBlock 类型的 block 替代，编译器已经默认将block copy到堆区**

### __block

在block中无法修改外部变量是因为变量是在申明block时，被复制到block的结构体中，而在变量前加上了__block的关键字，会增加一个结构体保存我们要修改的变量，在使用时引用的是该结构的指针，这样就达到了可以修改外部变量的作用。

## KVO相关

### 实现原理



## 内存管理

内存管理技术：GC垃圾回收， 引用计数法

引用计数法可以及时的回收引用计数为0的对象，减少查找次数。但是引用计数会带来循环引用的问题。



## OC对象

一个NSObject对象实力占16个字节



![image-20220222100823130](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_22_10_08_24.png)

## Category和Extension

### 区别

1. extension扩展就是类的一部分，在编译期和头文件声明以及类实现一起形成一个完整的类，一般用来隐藏类的私有信息，无法为系统类添加extension
2. category类目则不一样，它是在运行期决定的，无法添加实例变量。



## 启动优化

### 启动过程

#### 1. pre-main阶段

* 加载应用的可执行文件
* 加载动态连接器dyld