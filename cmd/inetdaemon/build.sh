# author: d1y<chenhonzhou@gmail.com>
# 该脚本只编译在小米路由器4c上
# 编写时间: 2021/03/13

GOOS=linux GOARCH=mipsle GOMIPS=softfloat CGO_ENABLED=0 go build -ldflags "-s -w" -o inetd .
# upx -9 inetd