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

category的底层结构是_category_t的结构体，所含的方法在运行中通过runtime动态的将方法合并到类对象和元类对象中

![image-20220223111528215](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_23_11_15_29.png)

### 区别

1. extension扩展就是类的一部分，在编译期和头文件声明以及类实现一起形成一个完整的类，一般用来隐藏类的私有信息，无法为系统类添加extension

2. category类目则不一样，它是在运行期决定的，无法添加实例变量。

   

## load和initialize

### 区别

#### 调用方式

1. load方式是根据函数地址直接调用
2. initialize是通过objc_msgSend调用

#### 调用时刻

1. load是runtime加载类、分类的时候调用（只会调用一次）
2. initialize是类第一次收到消息的时候调用，每一个类只会初始化一次（父类的initialize方法可能被调用多次）

#### 调用顺序

##### load

1. 先调用类的load，先编译的类先调用
2. 调用子类的load之前，会先调用父类的load
3. 再调用分类的load，先编译的分类，先调用

##### initialize

1. 先初始化父类
2. 再初始化子类（子类没有则调用父类的initialize方法）

## 启动优化

### 启动过程

#### 1. pre-main阶段

* 加载应用的可执行文件
* 加载动态连接器dyld