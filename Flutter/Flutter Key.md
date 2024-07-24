# Flutter Key

## Key

![image-20240724163838499](https://raw.githubusercontent.com/zpfate/ImageService/master/uPic/1721810323781)

**Flutter key 子类包含 LocalKey 和 GlobalKey。**

1. 局部键（LocalKey）：ValueKey、ObjectKey、UniqueKey。
2. 全局键（GlobalKey）： GlobalKey、GlobalObjectKey。

* ValueKey （值key）：把一个值作为 key 。

* UniqueKey（唯一key）：程序生成唯一的 Key，当我们不知道如何指定 ValueKey 的时候就可以使用 UniqueKey。
* ObjectKey（对象 key）：把一个对象实例作为 key。
* GlobalKey：每个 Widget 都对应一个 Element ，我们可以直接对 Widget 进行操作，但是无法直接操作 Widget 对应的 Element 。而 GlobalKey 就是那把直接访问 Element 的钥匙。通过 GlobalKey可以获取到 Widget 对应的 Element 。

## Key的作用

### 唯一标识符
Key的主要作用是唯一标识Widget。当Widget结构发生变化时，Flutter会根据Key来识别旧Widget和新Widget的不同，以确定是否需要进行重建。如果没有Key，Flutter将无法区分不同的Widget，会将它们全部重建，导致性能下降。

### 优化性能
通过Key，Flutter可以高效地更新Widget树中的部分子树，而不是重新构建整个Widget树。这样可以避免不必要的重建，提高性能。

### 避免重建问题
在某些情况下，我们希望保持Widget的状态不变，即使Widget的位置发生了改变。Key可以帮助我们解决这个问题，通过给Widget添加唯一的Key，确保在更新Widget树时能够正确地保留原来的状态，而不是将其重置为默认值。
