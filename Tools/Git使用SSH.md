

先放一下github的添加ssh key教程链接[添加ssh key教程链接](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)

### 查看有无ssh key

```shell
open ./.ssh
```



### 创建ssh key

```shell
ssh-keygen -t rsa -b 2048 -C "email@example.com"
```

运行后显示(这是选择存储key的位置，可以直接`Enter`)

```shell
Generating public/private rsa key pair.
Enter file in which to save the key (/home/user/.ssh/id_rsa):
```

下一步（是否输入密码，建议直接`Enter`跳过，不然使用git的时候一直让输入密码）：

```shell
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
```

### 复制ssh key

```shell
pbcopy < ~/.ssh/id_ed25519.pub
```

### 测试

```shell
ssh -T expample@github.com
```

