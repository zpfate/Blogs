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