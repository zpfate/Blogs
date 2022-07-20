## autorelease

* 自动释放池的主要底层数据结构是：__AutoReleasePool，AutoReleasePoolPage
* 调用了autorelease的对象最终都是通过AutoReleasePoolPage
* 每个AutoReleasePoolPage对象占用4096字节内存，除了用来存放他内部的成员变量，剩下的空间用来存放autorelease对象的地址
* 所有的AutoReleasePoolPage对象是通过双向链表的形式连接在一起

## AutoreleasePoolPage的结构

* `magic` 用来校验 AutoreleasePoolPage 的结构是否完整；
* `next` 指向最新添加的 autoreleased 对象的下一个位置，初始化时指向 begin() ；
* `thread` 指向当前线程；
* `parent` 指向父结点，第一个结点的 parent 值为 nil ；
* `child` 指向子结点，最后一个结点的 child 值为 nil ；
* `depth` 代表深度，从 0 开始，往后递增 1；
* hiwat` 代表 high water mark 。

![img](https://upload-images.jianshu.io/upload_images/5835116-170e421948bfb845.png?imageMogr2/auto-orient/strip|imageView2/2/w/836/format/webp)

* 调用push方法会将一个POOL_BOUNDARY入栈，并且会返回其存放的内存地址
* 调用pop放时传入一个POOL_BOUNDARY的内存地址，会从最后一个入栈的对象开始发送release消息，直到遇到这个POOL_BOUNDARY
* id *next指向了下一个能存放autorelease对象地址的区域

## RunLoop和AutoRelease

iOS在主线程的RunLoop注册了2个Observer

* 第一个Observer监听了KCFRunLoopEntry事件，会调用objc_autoreleasePoolPush()

* 第二个Observer监听了KCFRunLoopBeforeWaiting事件，会调用objc_autoreleasePoolPop()、objc_autoreleasePoolPush()，监听了KCFRunLoopBeforeExit事件，会调用objc_autoreleasePoolPop()

释放是由RunLoop来控制的

