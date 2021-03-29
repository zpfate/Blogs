### AVAudioSession

AVAudioSession是用来管理APP对音频硬件(扬声器 麦克风)的使用

#### AVAudioSession API见解以及使用

1.获取会话单例

```objective-c
AVAudioSession *session = [AVAudioSession sharedInstance];
```

2.设置category

```objective-c
NSError *sessionErr;
[session setCategory:AVAudioSessionCategoryPlayAndRecord error:&sessionErr];
```

![AudioSession的Category参数](https://i.loli.net/2020/08/07/TaCBfmrYbNOpUAK.png "AudioSession的Category参数")

3. 