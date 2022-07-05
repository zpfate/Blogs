一个循环 do-while

 全局的Dictionary，key 是 pthread_t， value 是 CFRunLoopRef

 RunLoop 之间是一一对应的，其关系是保存在一个全局的 Dictionary 里。线程刚创建时并没有 RunLoop，如果你不主动获取，那它一直都不会有。RunLoop 的创建是发生在第一次获取时，RunLoop 的销毁是发生在线程结束时。你只能在一个线程的内部获取其 RunLoop（主线程除外）。

![RunLoop_1](https://blog.ibireme.com/wp-content/uploads/2015/05/RunLoop_1.png)