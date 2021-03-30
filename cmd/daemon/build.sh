# create by d1y<chenhonzhou@gmail.com>
# write date 2021/03/30
# 该脚本只编译在嵌入式设备里用来不断发送心跳包

GOOS=linux

# 如果编译在路由器上
GOARCH=mipsle
GOMIPS=softfloat
# ==============

CGO_ENABLED=0

go build -ldflags "-s -w" -o idaemon .