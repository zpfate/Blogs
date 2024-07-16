## copy和mutableCopy

1. copy 不可变拷贝， 产生不可变副本
2. mutableCopy可变拷贝，产生可变副本

|             | NSString                                             | NSMutableString                                      | NSArray                                             | NSMutableArray                                      | NSDictionary                                             | NSMutableDictionary                                      |
| ----------- | ---------------------------------------------------- | ---------------------------------------------------- | --------------------------------------------------- | --------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| copy        | NSString<br /><font color = green>浅拷贝</font>      | NSString<br /><font color = red>深拷贝</font>        | NSArray<br /><font color = green>浅拷贝</font>      | NSArray<br /><font color = red>深拷贝</font>        | NSDictionary<br /><font color = green>浅拷贝</font>      | NSDictionary<br /><font color = red>深拷贝</font>        |
| mutableCopy | NSMutableString<br /><font color = red>深拷贝</font> | NSMutableString<br /><font color = red>深拷贝</font> | NSMutableArray<br /><font color = red>深拷贝</font> | NSMutableArray<br /><font color = red>深拷贝</font> | NSMutableDictionary<br /><font color = red>深拷贝</font> | NSMutableDictionary<br /><font color = red>深拷贝</font> |



## weak

weak是弱引用，用weak来修饰的引用对象的计数器不会增加，而且weak会在引用对象释放的时候，自动置为nil



苹果为了管理所有对象的计数器和weak指针，苹果创建了一个全局的哈希表，我们暂且叫它SideTables，里面装的是的名为SideTable的结构体。用对象的地址作为key，可以取出sideTable结构体，这个结构体用来管理引用计数和weak指针。

### objc_initWeak

1.初始化一个weak对象时，runtime会调用一个objc_initWeak函数，初始化一个新的weak指针指向该对象的地址

```c++
id objc_initWeak(id *location, id newObj) {
    // 查看对象实例是否有效
    // 无效对象直接导致指针释放
    if (!newObj) {
        *location = nil;
        return nil;
    }
    // 这里传递了三个 bool 数值
    // 使用 template 进行常量参数传递是为了优化性能
    return storeWeak<false/*old*/, true/*new*/, true/*crash*/>
        (location, (objc_object*)newObj);
}
```

### objc_storeWeak

2.在objc_initWeak函数中会继续调用objc_storeWeak函数，在这个过程是用来更新weak指针的指向，同时创建对应的弱引用表

```c++
template <bool HaveOld, bool HaveNew, bool CrashIfDeallocating>
static id 
storeWeak(id *location, objc_object *newObj)
{
    assert(HaveOld  ||  HaveNew);
    if (!HaveNew) assert(newObj == nil);

    Class previouslyInitializedClass = nil;
    id oldObj;
    SideTable *oldTable;
    SideTable *newTable;

    // Acquire locks for old and new values.
    // Order by lock address to prevent lock ordering problems. 
    // Retry if the old value changes underneath us.
 retry:
    if (HaveOld) {
        oldObj = *location;
        oldTable = &SideTables()[oldObj];
    } else {
        oldTable = nil;
    }
    if (HaveNew) {
        newTable = &SideTables()[newObj];
    } else {
        newTable = nil;
    }

    SideTable::lockTwo<HaveOld, HaveNew>(oldTable, newTable);

    if (HaveOld  &&  *location != oldObj) {
        SideTable::unlockTwo<HaveOld, HaveNew>(oldTable, newTable);
        goto retry;
    }

    // Prevent a deadlock between the weak reference machinery
    // and the +initialize machinery by ensuring that no 
    // weakly-referenced object has an un-+initialized isa.
    if (HaveNew  &&  newObj) {
        Class cls = newObj->getIsa();
        if (cls != previouslyInitializedClass  &&  
            !((objc_class *)cls)->isInitialized()) 
        {
            SideTable::unlockTwo<HaveOld, HaveNew>(oldTable, newTable);
            _class_initialize(_class_getNonMetaClass(cls, (id)newObj));

            // If this class is finished with +initialize then we're good.
            // If this class is still running +initialize on this thread 
            // (i.e. +initialize called storeWeak on an instance of itself)
            // then we may proceed but it will appear initializing and 
            // not yet initialized to the check above.
            // Instead set previouslyInitializedClass to recognize it on retry.
            previouslyInitializedClass = cls;

            goto retry;
        }
    }

    // Clean up old value, if any.
    if (HaveOld) {
        weak_unregister_no_lock(&oldTable->weak_table, oldObj, location);
    }

    // Assign new value, if any.
    if (HaveNew) {
        newObj = (objc_object *)weak_register_no_lock(&newTable->weak_table, 
                                                      (id)newObj, location, 
                                                      CrashIfDeallocating);
        // weak_register_no_lock returns nil if weak store should be rejected

        // Set is-weakly-referenced bit in refcount table.
        if (newObj  &&  !newObj->isTaggedPointer()) {
            newObj->setWeaklyReferenced_nolock();
        }

        // Do not set *location anywhere else. That would introduce a race.
        *location = (id)newObj;
    }
    else {
        // No new value. The storage is not changed.
    }
    
    SideTable::unlockTwo<HaveOld, HaveNew>(oldTable, newTable);

    return (id)newObj;
}
```



### SideTable

```c++
struct SideTable {
	  // 因为操作对象的引用计数频率很快，因此系统在这里设置了一把自旋锁，保证是原子操作
    spinlock_t slock; 
	  // 引用计数器哈希表,根据对象地址查找对象的引用计数
    RefcountMap refcnts; 
	  // 维护weak指针的结构体
    weak_table_t weak_table; 
}
```

### weak_table_t

```c++
struct weak_table_t {
  	// 保存所有指向某一个对象的weak指针的一个数组，循环遍历此数组可找到指定对象的weak_entry_t结构体实例
    weak_entry_t *weak_entries; 
	  // 用来维护数组的size始终在一个合理的大小
    size_t    num_entries; 
    uintptr_t        mask;
    uintptr_t        max_hash_displacement;
};
```



### weak_entry_t

```c++
struct weak_entry_t {
   // 被指向的对象的地址，前面循环遍历查找的时候就是判断目标地址是否和它相等。
    union {
    DisguisedPtr<objc_object> referent;
        struct {
          // 可变数组,里面保存着所有指向这个对象的弱引用的地址。当这个对象被释放的时候，referrers里的所有指针都会被设置成nil。
            weak_referrer_t *referrers; 
            uintptr_t        out_of_line : 1;
            uintptr_t        num_refs : PTR_MINUS_1;
            uintptr_t        mask;
            uintptr_t        max_hash_displacement;
        };
        struct {
            // 只有4个元素的数组，默认情况下用它来存储弱引用的指针。当大于4个的时候使用referrers来存储指针。
            weak_referrer_t  inline_referrers[WEAK_INLINE_COUNT]; 
        };
    };
};
```



#### 当 weak 指向的对象被释放时，如何让 weak 指针置为 nil 的呢

* 调用 objc_release
* 因为对象的引用计数为0，所以执行 dealloc
* 在 dealloc 中，调用了 _objc_rootDealloc 函数
* 在 _objc_rootDealloc 中，调用了 object_dispose 函数
* 调用 objc_destructInstance
* 最后调用 objc_clear_deallocating,详细过程如下：
  a. 从 weak 表中获取废弃对象的地址为键值的记录
  b. 将包含在记录中的所有附有 weak 修饰符变量的地址，赋值为 nil
  c. 将 weak 表中该记录删除
  d. 从引用计数表中删除废弃对象的地址为键值的记录

```c++
{
    SideTable& table = SideTables()[this];

    // clear any weak table items
    // clear extra retain count and deallocating bit
    // (fixme warn or abort if extra retain count == 0 ?)
    table.lock();
    RefcountMap::iterator it = table.refcnts.find(this);
    if (it != table.refcnts.end()) {
        if (it->second & SIDE_TABLE_WEAKLY_REFERENCED) {
            weak_clear_no_lock(&table.weak_table, (id)this);
        }
        table.refcnts.erase(it);
 }
```





```c++
void 
weak_clear_no_lock(weak_table_t *weak_table, id referent_id) 
{
    //1、拿到被销毁对象的指针
    objc_object *referent = (objc_object *)referent_id;
 
    //2、通过 指针 在weak_table中查找出对应的entry
    weak_entry_t *entry = weak_entry_for_referent(weak_table, referent);
    if (entry == nil) {
        /// XXX shouldn't happen, but does with mismatched CF/objc
        //printf("XXX no entry for clear deallocating %p\n", referent);
        return;
    }
 
    //3、将所有的引用设置成nil
    weak_referrer_t *referrers;
    size_t count;
 
    if (entry->out_of_line()) {
        //3.1、如果弱引用超过4个则将referrers数组内的弱引用都置成nil。
        referrers = entry->referrers;
        count = TABLE_SIZE(entry);
    } 
    else {
        //3.2、不超过4个则将inline_referrers数组内的弱引用都置成nil
        referrers = entry->inline_referrers;
        count = WEAK_INLINE_COUNT;
    }
 
    //循环设置所有的引用为nil
    for (size_t i = 0; i < count; ++i) {
        objc_object **referrer = referrers[i];
        if (referrer) {
            if (*referrer == referent) {
                *referrer = nil;
            }
            else if (*referrer) {
                _objc_inform("__weak variable at %p holds %p instead of %p. "
                             "This is probably incorrect use of "
                             "objc_storeWeak() and objc_loadWeak(). "
                             "Break on objc_weak_error to debug.\n", 
                             referrer, (void*)*referrer, (void*)referent);
                objc_weak_error();
            }
        }
    }
 
    //4、从weak_table中移除entry
    weak_entry_remove(weak_table, entry);
}
```

