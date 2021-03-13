// Author: d1y<chenhonzhou@gmail.com>
// 守护进程, 自动登录, 发送心跳包

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jsxxzy/inet"
	"github.com/jsxxzy/inet/cmd/inetdaemon/logging"
	"github.com/robfig/cron"
)

//go:generate go-bindata template/

// 默认端口
var defaultPort = 2333

// 目标端口
// var targetPort int

// 默认的日志目录
var defaultLogDir = "runlog"

var logMiddleware *logging.Logging

type User struct {
	Username string
	Password string
}

func getUserConfig() User {
	var user = os.Getenv("dr_user")
	var psd = os.Getenv("dr_password")
	if len(user) <= 0 || len(psd) <= 0 {
		panic("额, 先配置好环境变量的用户名和密码\n用户名: dr_user\n密码: dr_password")
	}
	return User{Username: user, Password: psd}
}

var App = gin.Default()

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	var indexData, _ = AssetString("template/index.tmpl")
	_, err := t.New("index").Parse(indexData)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func main() {

	var user = getUserConfig()

	t, err := loadTemplate()
	if err != nil {
		log.Println("加载web模板失败")
		panic(err)
	}

	App.SetHTMLTemplate(t)

	App.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	App.GET("/", func(ctx *gin.Context) {
		var logString = logMiddleware.GetLog()
		ctx.HTML(http.StatusOK, "index", gin.H{
			"title": "大专人雄起",
			"logs":  logString,
		})
	})
	var port = fmt.Sprintf(":%d", defaultPort)

	c := cron.New()

	// 间歇检测机制
	var timeKeepAlive = 20
	var arg = os.Args
	if len(arg) >= 2 {
		var t = arg[1]
		var i, err = strconv.Atoi(t)
		if err == nil {
			timeKeepAlive = i
		}
	}

	spec1 := fmt.Sprintf("*/%d * * * * ?", timeKeepAlive)

	c.AddFunc(spec1, func() {
		if !inet.HasLogin() {
			logMiddleware.Info("当前未登录, 尝试登录中")
			var loginInfo, errWare = inet.Login(user.Username, user.Password)
			if errWare != nil {
				logMiddleware.Error("登录错误, " + errWare.Error())
			} else {
				logMiddleware.Info("登录成功, " + loginInfo.GetMsg())
			}
		} else {
			logMiddleware.Info("当前已登录")
		}
	})
	c.Start()

	App.Run(port)
}

func init() {
	logMiddleware = logging.New(defaultLogDir)
}
