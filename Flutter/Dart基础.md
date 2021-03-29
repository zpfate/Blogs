

## final，const和static理解以及区别

### final

赋值后不能改变，不要求编译期是常数。

1. 在文件中声明的时候必须赋值
2. 如果是类的成员变量，可以声明时就赋值，也可以通过构造函数赋值

![final使用](https://cdn.jsdelivr.net/gh/ZpFate/ImageService@master/uPic/img_2020_10_28_13_35_57.png "final使用")

### const

const与final不一样的是，要求被修饰的对象在编译时就确定，并且对象完全不可变。

### static

static修饰的对象表示是类变量。static和上面两个关键字的差别还是很大的。

### final与const更具体的区别

举个例子会更直观一点：



## Future

Dart中Future是一个抽象类

Future提供了async和await关键字来支持异步编程。