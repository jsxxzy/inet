# 校园网登录客户端

先连接学校的网络, `wifi`和`有线`都可以(当然咯!你也可以使用宽带拨号上网)


```shell

# cli
go get github.com/jsxxzy/inet/cmd

# use package
go get github.com/jsxxzy/inet

```

使用例子

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

