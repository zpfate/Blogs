## Block相关

block本质是一个OC对象，内部也有isa指针，封装了函数调用以及函数调用环境

![img](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_16_09_53_59.jpg)

1. isa 指针，所有对象都有该指针，用于实现对象相关的功能。
2. flags，用于按 bit 位表示一些 block 的附加信息，本文后面介绍 block copy 的实现代码可以看到对该变量的使用。
3. reserved，保留变量。
4. invoke，函数指针，指向具体的 block 实现的函数调用地址。
5. descriptor， 表示该 block 的附加描述信息，主要是 size 大小，以及 copy 和 dispose 函数的指针。
6. variables，capture捕获 过来的变量，block 能够访问它外部的局部变量，就是因为将这些变量（或变量的地址）复制到了结构体中。

### 变量捕获

![image-20220302170846862](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_02_17_08_47.png)

<font color=red>*auto 默认的变量申明关键字，与static相对*</font>



![image-20220303092726423](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_03_09_27_26.png

![image-20220307135319690](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_07_13_53_20.png)

### 三种block

1. NSGlobalBlock 全局的静态 block，不会访问任何auto变量。
2. NSStackBlock 保存在栈中的 block，访问了auto变量，当函数返回时会被销毁。
3. NSMallocBlock 保存在堆中的 block，NSStackBlock调用了copy变成NSMallocBlock， 当引用计数为 0 时会被销毁。

![image-20220303101903412](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_03_10_19_03.png)

### __block

* \_\_block可以用于解决block内部无法修改auto变量值问题__
* block不能修饰全局变量、静态变量（static）
* 编译器会将__block变量包装成一个对象

### 修饰词用copy

没有copy操作，就不会在堆上，无法控制生命周期。

### block内存管理

* 当block在栈上时，不会对__block变量， 对象类型的auto变量产生强引用
* 当block被copy到堆时 

   * 会调用到block内部的copy函数
   * copy函数会调用_Block_object_assign函数
   * _Block_object_assign函数会对__block变量形成强引用(retain）
* 当block从堆中移除时

  * 会调用block内部的dispose函数
  * dispose函数内部会调用_Block_object_dispose函数
  * _Block_object_dispose函数会自动释放引用的__block变量(release)


##### __block修饰对象

* 当__block变量在栈上时，不会对指向的对象产生强引用

* 当__block变量被copy到堆时

  * 会调用__block内部的copy函数

   * copy函数内部会调用_Block_object_assign函数

   * _Block_object_assign函数会根据所指向对象的修饰符（\_\_Strong、\_\_weak、__unsafe_unretained）做出相应操作，形成强引用（retain)或者弱引用 

     **注意 ** 这里仅限于ARC时会retain，MRC时不会retain

* 当__block变量从堆上移除时

   * 会调用__block内部的dispose函数
   * dispose函数内部会调用_Block_object_dipose函数
   * _Block_object_dipose函数会自动释放指向的对象（release）

#### 解决循环引用

##### MRC

* 用__unsafe_unretained
* 用__block解决

##### ARC

* __unsafe_unretained (不会产生强引用，不安全，指向的对象销毁时，指针存储的地址值不变)

  ```objective-c
  __unsafe_unretained typeof(self) weakSelf = self;
  ```

* __weak (不会产生强引用，指向的对象销毁时会自动将指针置为nil)

  ```objective-c
  __weak typeof(self) weakSelf = self;
  ```

* __block （必须调用block，并将\_\_block修饰的变量置为nil）

  ```objective-c
    __block Person *p = [[Person alloc] init];
    p.age = 20;
    p.block = ^{
       NSLog(@"age is %zd", p.age);
       p = nil;
    };
    p.block();
  ```

  

  ![image-20220316165156112](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1647850786.png)



## KVO相关

Key-Value Observing 键值监听

### 实现原理

![image-20220314222151383](https://github.com/zpfate/uPic/2022%2003%20141647269033.png)

在程序运行中，动态创建一个NSKVONotifiying_XXX类



## 内存管理

内存管理技术：GC垃圾回收， 引用计数法

引用计数法可以及时的回收引用计数为0的对象，减少查找次数。但是引用计数会带来循环引用的问题。



## OC对象

### NSObject



### class、metaclass



### isa、superclass



![image-20220222100823130](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_22_10_08_24.png)

## Category和Extension

category的底层结构是_category_t的结构体，所含的方法在运行中通过runtime动态的将方法合并到类对象和元类对象中

![image-20220223111528215](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_23_11_15_29.png)

### 区别

1. extension扩展就是类的一部分，在编译期和头文件声明以及类实现一起形成一个完整的类，一般用来隐藏类的私有信息，无法为系统类添加extension
2. category类目则不一样，它是在运行期决定的，无法添加实例变量。
3. category不能直接添加成员变量，可以通过runtime关联对象间接添加

### 关联对象原理

实现关联对象技术的核心对象：

* AssociationsManager

* AssociationHashMap

* ObjectAssociationMap

* ObjcAssociation

  ![image-20220224154445087](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_24_15_44_46.png)

  ![image-20220224154855636](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_24_15_48_55.png)

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



## Runtime

Objective-C是一门动态性编程语言，动态性由Runtime API支撑

Runtime主要由c、c++、汇编来编写

### isa

在arm64位架构之前，isa就是一个普通的指针，存储着Class、Meta-Class对象的内存地址

从arm64架构开始，对isa进行了优化，变成了一个共用体（union）结构，还使用位域存储更多信息

![image-20220321154521223](https://cdn.jsdelivr.net/gh/zpfate/ImageService@master/uPic/1647848721.png)



## 启动优化

### 启动过程

#### 1. pre-main阶段

* 加载应用的可执行文件
* 加载动态连接器dyld