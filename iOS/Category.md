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