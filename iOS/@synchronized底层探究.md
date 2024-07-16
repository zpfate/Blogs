

# @synchronized底层探究

## @synchronized分析

先将下面的代码使用`clang`命令分析一下：

```objc
// main.m
#import <UIKit/UIKit.h>
#import "AppDelegate.h"
    
int main(int argc, char * argv[]) {
    NSString * appDelegateClassName;
    @autoreleasepool {
        
        appDelegateClassName = NSStringFromClass([AppDelegate class]);
        @synchronized (appDelegateClassName) {
        }
    }
    return UIApplicationMain(argc, argv, nil, appDelegateClassName);
}
```

### @Synchronized分析成c++

```objc
int main(int argc, char * argv[]) {
    NSString * appDelegateClassName;
    / * @autoreleasepool * / { __AtAutoreleasePool __autoreleasepool;
        
        appDelegateClassName = NSStringFromClass(((Class (*)(id, SEL))(void *)objc_msgSend)((id)objc_getClass("AppDelegate"), sel_registerName("class")));
        {
            id _rethrow = 0;
            id _sync_obj = (id)appDelegateClassName;
            objc_sync_enter(_sync_obj);
            try {
                struct _SYNC_EXIT { _SYNC_EXIT(id arg) : sync_exit(arg) {}
                    ~_SYNC_EXIT() {
                        objc_sync_exit(sync_exit);
                    }
                    id sync_exit;
                } _sync_exit(_sync_obj);
                ......
```

对照`clang`命令转换成的c++代码可以看出:

**@synchronized 是如下两个方法包裹实现的,objc_sync_enter(); objc_sync_exit();**

### objc_sync_enter、objc_sync_exit源码

这两个方法可以去objc源码的objc-sync.mm文件中可以查找到实现:

```c++
// Begin synchronizing on 'obj'. 
// Allocates recursive mutex associated with 'obj' if needed.
// Returns OBJC_SYNC_SUCCESS once lock is acquired.  
int objc_sync_enter(id obj)
{
    int result = OBJC_SYNC_SUCCESS;

    if (obj) {
        SyncData* data = id2data(obj, ACQUIRE);
        assert(data);
        data->mutex.lock();
    } else {
        // @synchronized(nil) does nothing
        if (DebugNilSync) {
            _objc_inform("NIL SYNC DEBUG: @synchronized(nil); set a breakpoint on objc_sync_nil to debug");
        }
        objc_sync_nil();
    }

    return result;
}
// End synchronizing on 'obj'. 
// Returns OBJC_SYNC_SUCCESS or OBJC_SYNC_NOT_OWNING_THREAD_ERROR
int objc_sync_exit(id obj)
{
    int result = OBJC_SYNC_SUCCESS;
    
    if (obj) {
        SyncData* data = id2data(obj, RELEASE); 
        if (!data) {
            result = OBJC_SYNC_NOT_OWNING_THREAD_ERROR;
        } else {
            bool okay = data->mutex.tryUnlock();
            if (!okay) {
                result = OBJC_SYNC_NOT_OWNING_THREAD_ERROR;
            }
        }
    } else {
        // @synchronized(nil) does nothing
    }
 

    return result;
}
```

### objc_sync_enter

1. 1.在obj上开始同步锁
2. obj为nil，加锁不会成功 
3. obj不是nil，初始化递归互斥锁（recursive mutex），并关联obj

### SyncData

```c++
//objc-sync.mm
typedef struct SyncData {
     //下一条同步数据
    struct SyncData* nextData;
    //锁的对象
    DisguisedPtr<objc_object> object;
    //等待的线程数量
    int32_t threadCount;  // number of THREADS using this block
    //互斥递归锁
    recursive_mutex_t mutex;
} SyncData;
```

* SyncData是一个结构体，类似链表

* nextData：SyncData的指针节点，指向下一条数据 

* object：锁住的对象 

* threadCount：等待的线程数量

* mutex：使用的递归互斥锁

### recursive_mutex_t

```c++
//objc-os.h
using recursive_mutex_t = recursive_mutex_tt<DEBUG>;
class recursive_mutex_tt : nocopy_t {
    pthread_mutex_t mLock;

  public:
    recursive_mutex_tt() : mLock(PTHREAD_RECURSIVE_MUTEX_INITIALIZER) { }
    void lock() {
        lockdebug_recursive_mutex_lock(this);
        int err = pthread_mutex_lock(&mLock);
        if (err) _objc_fatal("pthread_mutex_lock failed (%d)", err);
    }
    //这里省略......
}
```

## @synchronized的原理

内部为每一个obj分配一把recursive_mutex递归互斥锁。针对每个obj，通过这个recursive_mutex递归互斥锁进行加锁、解锁

## @synchronized总结

1. synchronized 的 obj 为 nil 怎么办？加锁操作无效。
2. synchronized 会对 obj 做什么操作吗？会为obj生成递归自旋锁，并建立关联，生成 SyncData，存储在当前线程的缓存里或者全局哈希表里。
3. synchronized 和 pthread_mutex 有什么关系？SyncData里的递归互斥锁，使用 pthread_mutex 实现的。
4. synchronized 和 objc_sync 有什么关系？synchronized 底层调用了 objc_sync_enter() 和 objc_sync_exit()

