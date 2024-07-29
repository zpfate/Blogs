## 什么是block

Block是将函数极其执行上下文封装起来的对象。

![image-20210315163225613](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_15_16_32_26.png)

通过`clang -rewrite-objc`命令编译该.m文件

![clang -rewrite-objc命令](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_15_16_36_24.png "clang -rewrite-objc命令")

发现`block`被编译成如下所示：

![block通过clang-rewrite-objc编译](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2021_03_15_16_40_31.png "block通过clang-rewrite-objc编译")

`block`本质是一个`OC`对象，内部也有`isa`指针，封装了函数调用以及函数调用环境

![img](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_02_16_09_53_59.jpg)

1. `isa`指针，所有对象都有该指针，用于实现对象相关的功能。
2. `flags`，用于按`bit`位表示一些`block`的附加信息，本文后面介绍`block copy`的实现代码可以看到对该变量的使用。
3. `reserved`，保留变量。
4. `invoke`，函数指针，指向具体的`block`实现的函数调用地址。
5. `descriptor`，表示该`block`的附加描述信息，主要是`size`大小，以及`copy`和`dispose`函数的指针。
6. `variables`，`capture`捕获过来的变量，`block`能够访问它外部的局部变量，就是因为将这些变量（或变量的地址）复制到了结构体中。

### 变量捕获

![image-20220302170846862](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_02_17_08_47.png)

<font color=red>*auto默认的变量申明关键字，与static相对*</font>

![image-20220303092726423](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_03_09_27_26.png

![image-20220307135319690](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_07_13_53_20.png)

### 三种block

1. `NSGlobalBlock` 全局的静态`block`，不会访问任何`auto`变量。
2. `NSStackBlock` 保存在栈中的`block`，访问了`auto`变量，当函数返回时会被销毁。
3. `NSMallocBlock` 保存在堆中的`block`，`NSStackBlock`调用了`copy`变成`NSMallocBlock`， 当引用计数为0时会被销毁。

![image-20220303101903412](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2022_03_03_10_19_03.png)

### __block

* `__block`可以用于解决`block`内部无法修改`auto`变量值问题
* `__block`不能修饰全局变量、静态变量（`static`）
* 编译器会将`__block`变量包装成一个对象

### 修饰词用copy

没有`copy`操作，就不会在堆上，无法控制生命周期。

### block内存管理

* 当`block`在栈上时，不会对`__block`变、对象类型的`auto`变量产生强引用
* 当`block`被`copy`到堆时 

  * 会调用到`block`内部的`copy`函数
  * `copy`函数会调用`_Block_object_assign`函数
  * `_Block_object_assign`函数会对`__block`变量形成强引用(`retain`）
* 当`block`从堆中移除时

  * 会调用`block`内部的`dispose`函数
  * `dispose`函数内部会调用`_Block_object_dispose`函数
  * `_Block_object_dispose`函数会自动释放引用的`__block`变量(`release`)


##### `__block`修饰对象

* 当`__block`变量在栈上时，不会对指向的对象产生强引用

* 当`__block`变量被`copy`到堆时

  * 会调用`__block`内部的`copy`函数

   * `copy`函数内部会调用`_Block_object_assign`函数

   * `_Block_object_assign`函数会根据所指向对象的修饰符（`__strong`、`__weak`、`__unsafe_unretained`）做出相应操作，形成强引用（`retain`)或者弱引用 

     **注意 ** 这里仅限于`ARC`时会`retain`，`MRC`时不会`retain`

* 当`__block`变量从堆上移除时

  * 会调用`__block`内部的`dispose`函数
  * `dispose`函数内部会调用`_Block_object_dipose`函数
  * `_Block_object_dipose`函数会自动释放指向的对象（`release`）

#### 解决循环引用

##### MRC

* 用`__unsafe_unretained`
* 用`__block`解决

##### ARC

* `__unsafe_unretained` (不会产生强引用，不安全，指向的对象销毁时，指针存储的地址值不变)

  ```objective-c
  __unsafe_unretained typeof(self) weakSelf = self;
  ```

* `__weak`(不会产生强引用，指向的对象销毁时会自动将指针置为`nil`)

  ```objective-c
  __weak typeof(self) weakSelf = self;
  ```

* `__block`（必须调用`block`，并将`__block`修饰的变量置为`nil`）

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

