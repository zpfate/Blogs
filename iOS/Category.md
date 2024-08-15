

## Category

category的底层结构是_category_t的结构体，所含的方法在运行中通过runtime动态的将方法合并到类对象和元类对象中
![image-20220223111528215](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_23_11_15_29.png)

### 和Extension区别
1.  extension（扩展）就是类的一部分，在编译期和头文件声明以及类实现一起形成一个完整的类，一般用来隐藏类的私有信息，无法为系统类添加extension
2.  category（类目）则不一样，它是在运行期决定的，无法添加实例变量。
3.  category不能直接添加成员变量，可以通过runtime关联对象间接添加

### 关联对象原理

实现关联对象技术的核心对象：

-   AssociationsManager
-   AssociationHashMap
-   ObjectAssociationMap
-   ObjcAssociation

![image-20220224154445087](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_24_15_44_46.png)

![image-20220224154855636](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_24_15_48_55.png)



![img](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1722416316433)

## Category不能添加属性的原因

1. **编译时特性**：Objective-C的属性（@property）是在编译时由编译器处理的。当你在类中声明一个属性时，编译器会自动为你生成相应的实例变量（在.m文件中以_前缀表示）以及getter和setter方法。这个过程是在编译期间完成的，而category是在运行时加载的，因此编译器无法为category中的属性生成实例变量和访问方法。

2. **内存布局**：类的内存布局在编译时就已经确定，包括所有实例变量的位置和大小。如果在运行时通过category添加属性，将会破坏这个内存布局，可能导致内存访问错误或其他不可预测的行为。

3. **动态特性**：Objective-C的运行时（Runtime）允许动态地添加方法，因为方法的添加不会影响到对象的内存布局。但是，属性的添加涉及到实例变量的创建，这是在运行时无法动态改变的。

