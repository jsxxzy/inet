# ==>> 校园网登录客户端

先连接学校的网络, `wifi`和`有线`都可以(当然咯!你也可以使用宽带拨号上网)


```shell

# cli
go get github.com/jsxxzy/inet/cmd/inet

# use package
go get github.com/jsxxzy/inet

```

`cli` 使用

```shell
~/c/o/j/inet ❯❯❯ inet

=============>
 login: 登录
   get: 获取保存的账号密码
  save: 保存账号密码
   fix: 初始化账号密码
 check: 查询是否登录
  info: 查询信息
logout: 注销
=============>

Usage of inet:
  -password string
        密码
  -username string
        用户名


# 使用账号密码登录
inet -username "666@jszy" -password "6666" login

# 存储到本地, 下次就可以直接 inet login
inet -username "666@jszy" -password "6666" save

# 查询本地的账号密码
inet get

# 初始化, 既重置本地的账号密码
inet fix

# 检测是否登录
inet check

# 查询使用信息(登录之后)
inet info

# 注销
inet logout
```

包使用例子

```go

var username, password = "用户名", "密码"
info, _ := inet.Login(username, password)

fmt.Println(info.GetMsg())
fmt.Println(info.Error())

data, _ := inet.QueryInfo()

fmt.Println("data", data)

if inet.HasLogin() {
  inet.Logout()
}


  ```

![upload](https://i.loli.net/2020/10/28/wDsoYXupFA7f5aG.png)

