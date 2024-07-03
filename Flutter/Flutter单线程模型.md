# 单线程模型

## Event Loop机制

Dart是单线程的，**单线程和异步不冲突**。

App应用程序绝大多数时间都在等待，等用户点击、等网络请求返回，文件IO结果等。基于这些特点，单线程模型可以在等待中做别的事情。等待的过程并不是阻塞的，这个行为就是`Event Loop`驱动的。事件队列`Event Queue`会把其他异步操作的完成、需要主线程响应的事件加入其中。

有点像iOS的`RunLoop`，不断轮询事件队列，取出事件（比如，键盘事件、I\O 事件、网络事件等），在主线程同步执行其回调函数。



![简化版 Event Loop](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1719902213899 "简化版 Event Loop")

## 异步任务

![Microtask Queue 与 Event Queue](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1719903548321 "Microtask Queue 与 Event Queue")



### microtask

首先，我们看看微任务队列。微任务顾名思义，表示一个短时间内就会完成的异步任务。从上面的流程图可以看到，微任务队列在事件循环中的优先级是最高的，只要队列中还有任务，就可以一直霸占着事件循环。微任务是由 `scheduleMicroTask()`或者`Future.microtask()` 建立的。

### Future

**<font color=redaccent>Dart 为 Event Queue 的任务建立提供了一层封装，叫作 Future。</font>**

Future 还提供了链式调用的能力，可以在异步任务执行完毕后依次执行链路上的其他函数体。

正常情况下，一个 Future 异步任务的执行是相对简单的：**在我们声明一个 Future 时，Dart 会将异步任务的函数执行体放入事件队列，然后立即返回，后续的代码继续同步执行**。而当同步执行的代码执行完毕后，事件队列会按照加入事件队列的顺序（即声明顺序），依次取出事件，最后同步执行 Future 的函数体及后续的 then。这意味着，then 与 Future 函数体共用一个事件循环。而如果 Future 有多个 then，它们也会按照链式调用的先后顺序同步执行，同样也会共用一个事件循环。



如果 Future 执行体已经执行完毕了，但你又拿着这个 Future 的引用，往里面加了一个 then 方法体，这时 Dart 会如何处理呢？

面对这种情况，**Dart 会将后续加入的 then 方法体放入微任务队列，尽快执行**。

### Stream

`Stream`中的执行异步的模式就是`scheduleMicrotask`。因为microtask的优先级又高于event。所以，如果 `microtask 太多`就可能会对触摸、绘制等外部事件造成`阻塞卡顿`。

`Stream` 有**同步流和异步流**之分。它们的区别在于**同步流**会在执行``add`，`addError` 或 `close` 方法时立即向流的监听器 `StreamSubscription` 发送事件，而**异步流**总是在事件队列中的代码执行完成后在发送事件。`

>**Future和Stream区别**
>
>1. `Future`中的任务会加入下一轮事件循环，`Stream`中的任务加入的是微任务队列
>2. `Future` 用于表示单个运算的结果，而 `Stream` 则表示多个结果的序列。

#### Stream订阅模式

`Stream`的订阅模式分为`Single`和`Broadcast`两种类型， 前者只允许订阅（`listen`）一次，后者允许多次订阅。

#### Stream单订阅，多次订阅会出现什么结果？

会报错，单订阅只能有一次订阅.
即使取消了第一个监听器，也不允许在单订阅流上设置其他的监听器。

`Stream` 可以通过 `transform()` 方法（返回另一个 `Stream`）进行连续调用。
通过 `Stream.asBroadcastStream()` 可以将一个单订阅模式的 `Stream` 转换成一个多订阅模式的 `Stream.isBroadcast` 属性可以判断当前 `Stream` 所处的模式。

## 异步函数

### await

对于一个异步函数来说，其返回时内部执行动作并未结束，因此需要返回一个 `Future` 对象，供调用者使用。调用者根据`Future`对象，来决定：是在这个`Future`对象上注册一个 `then`，等`Future`的执行体结束了以后再进行异步处理；还是一直同步等待`Future`执行体结束。对于异步函数返回的`Future`对象，如果调用者决定同步等待，则需要在调用处使用`await`关键字，并且在调用处的函数体使用`async`关键字。

因为 **Dart 中的`await`并不是阻塞等待，而是异步等待**。

```dart

// 声明了一个延迟2秒返回Hello的Future，并注册了一个then返回拼接后的Hello 2019
Future<String> fetchContent() => 
  Future<String>.delayed(Duration(seconds:2), () => "Hello")
    .then((x) => "$x 2019");
// 异步函数会同步等待Hello 2019的返回，并打印
func() async => print(await fetchContent());

main() {
  print("func before");
  func();
  print("func after");
}

```

### await for

`await for`是用来不断获取stream流中的数据，然后执行循环体中的操作。它一般用在直到`Stream`什么时候完成，并且必须等待传递完成后才能使用，不然会阻塞。

```dart

Stream<String> stream = new Stream<String>.fromIterable(['1', '2', '3','4']);

main() async {
    await for(String s in stream){
    print(s);
  }
}
```



## Ioslate

尽管 Dart 是基于单线程模型的，但**为了进一步利用多核 CPU，将 CPU 密集型运算进行隔离，Dart 也提供了多线程机制，即 Isolate**。在 Isolate 中，资源隔离做得非常好，**每个 Isolate 都有自己的 Event Loop 与 Queue，Isolate 之间不共享任何资源，只能依靠消息机制通信，因此也就没有资源抢占问题**。

```dart
Isolate isolate;

start() async {
  ReceivePort receivePort= ReceivePort();// 创建管道
  // 创建并发Isolate，并传入发送管道
  isolate = await Isolate.spawn(getMsg, receivePort.sendPort);
  // 监听管道消息
  receivePort.listen((data) {
    print('Data：$data');
    receivePort.close();// 关闭管道
    isolate?.kill(priority: Isolate.immediate);// 杀死并发Isolate
    isolate = null;
  });
}
// 并发Isolate往管道发送一个字符串
getMsg(sendPort) => sendPort.send("Hello");
```

```dart

// 并发计算阶乘
Future<dynamic> asyncFactoriali(n) async{
  final response = ReceivePort();//创建管道
  // 创建并发Isolate，并传入管道
  await Isolate.spawn(_isolate,response.sendPort);
  // 等待Isolate回传管道
  final sendPort = await response.first as SendPort;
  // 创建了另一个管道answer
  final answer = ReceivePort();
  // 往Isolate回传的管道中发送参数，同时传入answer管道
  sendPort.send([n,answer.sendPort]);
  return answer.first;//等待Isolate通过answer管道回传执行结果
}

// Isolate函数体，参数是主Isolate传入的管道
_isolate(initialReplyTo) async {
  final port = ReceivePort();// 创建管道
  initialReplyTo.send(port.sendPort);// 往主Isolate回传管道
  final message = await port.first as List;// 等待主Isolate发送消息(参数和回传结果的管道)
  final data = message[0] as int;// 参数
  final send = message[1] as SendPort;// 回传结果的管道 
  send.send(syncFactorial(data));// 调用同步计算阶乘的函数回传结果
}

// 同步计算阶乘
int syncFactorial(n) => n < 2 ? n : n * syncFactorial(n-1);
main() async => print(await asyncFactoriali(4));// 等待并发计算阶乘结果

```



## 

