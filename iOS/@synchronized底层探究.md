

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

### @synchronized分析成C++

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

```c++
static SyncData* id2data(id object, enum usage why)
    {
        spinlock_t *lockp = &LOCK_FOR_OBJ(object);
        SyncData **listp = &LIST_FOR_OBJ(object);
        SyncData* result = NULL;
        
#if SUPPORT_DIRECT_THREAD_KEYS //环境定义--针对单个线程
        // Check per-thread single-entry fast cache for matching object
        
        bool fastCacheOccupied = NO;
        // 检查每线程单项快速缓存中是否有匹配的对象--
        SyncData *data = (SyncData *)tls_get_direct(SYNC_DATA_DIRECT_KEY);
        if (data) {//
            fastCacheOccupied = YES;
            
            if (data->object == object) {
                // Found a match in fast cache.
                uintptr_t lockCount;
                
                result = data;
                //取出线程单项快速缓存中加锁的次数
                lockCount = (uintptr_t)tls_get_direct(SYNC_COUNT_DIRECT_KEY);
                if (result->threadCount <= 0  ||  lockCount <= 0) {
                    _objc_fatal("id2data fastcache is buggy");
                }
                
                switch(why) {
                    case ACQUIRE: {//加锁，则对加锁次数+1，然后重新缓存到线程的单项快速缓存中
                        lockCount++;
                        tls_set_direct(SYNC_COUNT_DIRECT_KEY, (void*)lockCount);
                        break;
                    }
                    case RELEASE://解锁，如果解锁，则对加锁次数-1，然后重新缓存到线程的单项快速缓存中
                        lockCount--;
                        tls_set_direct(SYNC_COUNT_DIRECT_KEY, (void*)lockCount);
                        if (lockCount == 0) {//lockCount为0时候就从tls中删除解锁
                            // remove from fast cache
                            tls_set_direct(SYNC_DATA_DIRECT_KEY, NULL);
                            // atomic because may collide with concurrent ACQUIRE--记录的线程数去掉，threadCount是记录的线程的情况
                            OSAtomicDecrement32Barrier(&result->threadCount);
                        }
                        break;
                    case CHECK:
                        // do nothing
                        break;
                }
                return result;
            }
        }
#endif//针对所有线程
        // Check per-thread cache of already-owned locks for matching object
        // 检查已拥有锁的每个线程高速缓存中是否有匹配的对象---跟上面操作类似
        SyncCache *cache = fetch_cache(NO);
        if (cache) {
            unsigned int i;
            for (i = 0; i < cache->used; i++) {
                SyncCacheItem *item = &cache->list[i];
                if (item->data->object != object) continue;
                // Found a match.
                result = item->data;
                if (result->threadCount <= 0  ||  item->lockCount <= 0) {
                    _objc_fatal("id2data cache is buggy");
                }
                switch(why) {
                    case ACQUIRE:
                        item->lockCount++;
                        break;
                    case RELEASE:
                        item->lockCount--;
                        if (item->lockCount == 0) {
                            // remove from per-thread cache
                            cache->list[i] = cache->list[--cache->used];
                            // atomic because may collide with concurrent ACQUIRE
                            OSAtomicDecrement32Barrier(&result->threadCount);
                        }
                        break;
                    case CHECK:
                        // do nothing
                        break;
                }
                return result;
            }
        }
        
        // Thread cache didn't find anything.
        // Walk in-use list looking for matching object
        // Spinlock prevents multiple threads from creating multiple
        // locks for the same new object.
        // We could keep the nodes in some hash table if we find that there are
        // more than 20 or so distinct locks active, but we don't do that now.
        //加锁
        lockp->lock();
        {
            /*
             再次从缓存中取，看看有没有，避免有其他线程已经对obj加锁
             如果找到了就对threadCount+1
             如果找不到就在SyncData接一个data
             
             如果不是加锁咋不管
             如果是加锁，则将新的data的objc设置为object，并且将threadCount加1
             
             如果一个SyncData都没有那么就创建一个，设置好数据后保存在map中
             最后如果是新创建的，那么就会在tls中保存数据
             */
            SyncData* p;
            SyncData* firstUnused = NULL;
            for (p = *listp; p != NULL; p = p->nextData) {
                if ( p->object == object ) {//再次从缓存中取一下
                    result = p;
                    // atomic because may collide with concurrent RELEASE
                    //线程+1
                    OSAtomicIncrement32Barrier(&result->threadCount);
                    goto done;
                }
                //如果已经有了SyncData但是没有对应的object
                if ( (firstUnused == NULL) && (p->threadCount == 0) )
                    firstUnused = p;
            }
            
            // no SyncData currently associated with object
            //如果不是加锁则不管
            if ( (why == RELEASE) || (why == CHECK) )
                goto done;
            
            // an unused one was found, use it
            if ( firstUnused != NULL ) {//则将object配置在后面的SyncData中
                result = firstUnused;
                result->object = (objc_object *)object;
                result->threadCount = 1;
                goto done;
            }
        }
        
        // Allocate a new SyncData and add to list.
        // XXX allocating memory with a global lock held is bad practice,
        // might be worth releasing the lock, allocating, and searching again.
        // But since we never free these guys we won't be stuck in allocation very often.
        //如果一个SyncData都没有，那就创建一个
        posix_memalign((void **)&result, alignof(SyncData), sizeof(SyncData));
        result->object = (objc_object *)object;
        result->threadCount = 1;
        new (&result->mutex) recursive_mutex_t(fork_unsafe_lock);
        result->nextData = *listp;
        *listp = result;
        
    done:
        lockp->unlock();
        if (result) {
            // Only new ACQUIRE should get here.
            // All RELEASE and CHECK and recursive ACQUIRE are
            // handled by the per-thread caches above.
            if (why == RELEASE) {
                // Probably some thread is incorrectly exiting
                // while the object is held by another thread.
                return nil;
            }
            if (why != ACQUIRE) _objc_fatal("id2data is buggy");
            if (result->object != object) _objc_fatal("id2data is buggy");
            
#if SUPPORT_DIRECT_THREAD_KEYS
            if (!fastCacheOccupied) {//如果不是从缓存中读取的，那么就将SyncData缓存到tls中
                // Save in fast thread cache
                tls_set_direct(SYNC_DATA_DIRECT_KEY, result);
                tls_set_direct(SYNC_COUNT_DIRECT_KEY, (void*)1);
            } else
#endif
            {
                // Save in thread cache
                if (!cache) cache = fetch_cache(YES);
                cache->list[cache->used].data = result;
                cache->list[cache->used].lockCount = 1;
                cache->used++;
            }
        }
        return result;
    }
```



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

*步骤1*: 当加锁的时候，系统会先找到当前线程的tls（线程局部存储） 中查找obj 对应的`SyncData`。

>1. 如果找到了执行 *步骤2*
>2. 如果没有找到就从全局的`SyncCache`中遍历寻找`obj`对应的`SyncCacheItem`。
>   ①如果取到了就到 *步骤2*。
>   ②如果没有取到就到 *步骤3*。

*步骤2*: 找到并取出`SyncData`或者`SyncCachetItem`，先取出对象的`lockCount`，代表的是加锁的次数。

> 1. 如果是加锁，则把`lockCount` +1 ，保存到线程的`tls` 
> 2. 如果是解锁，则把`lockCount` -1 ，保存到线程的`tls`中。如果减完之后的`lockCount`为0则将`tls`中`lockCount`置空（`nil`），并且通过 `OSAtomicDecrement32Barrier`将`SyncData`的线程数置空，然后将`SyncData`或者`SyncCacheItem`返回

 *步骤3*: 假如没找到，说明该对象以前没加锁过 或者是 在其他的`SyncData`中，为了保证不重复创建`SyncData`，进行如下操作。
    1. 先从`listp` 中遍历`SyncData`（源码中有注释，以前的做法是在多个`data`里对同一个对象加锁，所以找到`data`发下没有该对象的话继续去遍历因为有可能在其他的data中。所以不应该单纯的只查一个`data`而是去查整个链表。这个遍历过程加了锁，更安全），如果找到 obj 对应的`SyncData`（可能其他线程创建的），执行*步骤4*
    1. 如果没有找到`obj`对应的`SyncData`，但是有空白的`SyncData`，则将空白的`SyncData` 与obj 关联起来，并且将`SyncData`的`threadCount` (线程数)置为1，再执行 *步骤4*
    1. 假如没有空白的`SyncData`，则创建一个`SyncData`，然后与obj 关联起来将``SyncData`的`threadCount` 置为1，并且将创建的`SyncData`变为`listp`的第一个元素（`listp`是通过`obj` 哈希算法得到`index` 存储在`map`表中，因为不同`obj`可能得到相同的index，所以此时`listp`是已经有数据了，为了不让数据丢失，会吧数据赋值在取到数据的`nextData`中，类似栈特征，先进的在后面）


*步骤4*:判断是否是`release` 也就是是解锁那么直接返回nil（因为在这个当前线程没有找到，在其他线程中找到的。例如线程1的tls 中没有找到，在线程2中找到的）  准备好了`SyncData`，则解锁上面代码，如果是解锁，则不管直接`return`,如果是加锁，则将准备好的`SyncData`保存在当前线程tls的`SYNC_DATA_DIRECT_KEY` 中，并将tls中 `SYNC_DATA_DIRECT_KEY` 设置为1，然后将`SyncData`返回。

 记录`lockCount`进行递归调用，`id2data` （`why`)是根据加锁还是解锁存入不同给的`lockCount`可以重入，可以被锁多次，所以`@synchronized`是一个递归锁（可以被锁多次）。
 因为`@synchronized`记录了我们的线程和锁的情况，就能知道哪些线程造成了相互等待，这个时候只要通过处理这个相互等待就不会死锁（例如发送一个强大的命令什么的）

**由于@synchronized需要不断对map 表已经缓存进行读写操作，所以性能比较低**



1. synchronized 的 obj 为 nil 怎么办？加锁操作无效。
2. synchronized 会对 obj 做什么操作吗？会为obj生成递归自旋锁，并建立关联，生成 SyncData，存储在当前线程的缓存里或者全局哈希表里。
3. synchronized 和 pthread_mutex 有什么关系？SyncData里的递归互斥锁，使用 pthread_mutex 实现的。
4. synchronized 和 objc_sync 有什么关系？synchronized 底层调用了 objc_sync_enter() 和 objc_sync_exit()

