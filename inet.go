// Author: d1y<chenhonzhou@gmail.com>
// 职院校园网客户端
// 编写时间: 2020-10-29

package inet

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	js "github.com/dop251/goja"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	// ErrorNoAuth 未登录
	ErrorNoAuth = errors.New("未登录")
	// ErrorUserAuthFail 账号密码错误
	ErrorUserAuthFail = errors.New("账号密码错误")
	// ErrorMultipleDevices 多台设备同时在线
	ErrorMultipleDevices = errors.New("多台设备同时在线")
)

const (
	// ErrorNoAuthCode 未登录
	ErrorNoAuthCode = 0
	// ErrorLoginAuthCode 已登录
	ErrorLoginAuthCode = 2
	// LoginSuccess 登录成功
	LoginSuccess = 66
	// ErrorUserAuthFailCode 账号密码错误
	ErrorUserAuthFailCode = 1
	// ErrorMultipleDevicesCode 多台设备同时在线
	ErrorMultipleDevicesCode = 5
)

// AuthBaseURL 验证基础地址
//
//
const AuthBaseURL string = "http://210.22.55.58"

// LogoutAPI 注销接口
//
// ======
//
// 后台过来挨打来, `F.htm` 是什么意思??
//
// ======
var LogoutAPI string = createURL("/F.htm")

// LoginAPI 登录接口
var LoginAPI string = createURL("/0.htm")

// QueryInfoData 查询返回的数据
type QueryInfoData struct {

	// code 值
	code int

	// Portalname 名称
	Portalname string

	// Time 时间
	Time string

	// Flow 流量
	Flow float64

	// Xip 外网映射地址
	Xip string

	// UID 用户名(`id`)
	UID string

	// V4ip `ipv4` 地址
	V4ip string

	// V6ip `ipv6` 地址
	V6ip string
}

// LoginInfo 登录返回的结果
type LoginInfo struct {
	code int
	msg  string
}

func (L LoginInfo) Error() error {
	switch L.code {
	case ErrorUserAuthFailCode:
		return ErrorUserAuthFail
	case ErrorMultipleDevicesCode:
		return ErrorMultipleDevices
	}
	return nil
}

// GetMsg 获取消息
func (L LoginInfo) GetMsg() string {
	return L.msg
}

func (Qdata QueryInfoData) Error() error {
	switch Qdata.code {
	case ErrorNoAuthCode: // 未登录
		return ErrorNoAuth
	}
	return nil
}

// GetCode 返回代码
func (Qdata QueryInfoData) GetCode() int {
	return Qdata.code
}

// 拼接字符串
func createURL(p string) string {
	u, _ := url.Parse(AuthBaseURL)
	u.Path = p
	return u.String()
}

// 创建`md5`
func easyMD5(p string) string {
	var a = []byte(p)
	var b = fmt.Sprintf("%x", md5.Sum(a))
	return b
}

// =======================
// =======================
// => 作者提示:这里不应该写死!
// =======================
// =======================
var (
	pid  = "2"
	calg = "12345678"
	r1   = 0
	r2   = 1
)

// 创建密码, 逆向自: http://210.22.55.58/a41.js
//
// 生成的不严谨, 可能随时都会过期
func createPassword(p string) string {
	var p1 = pid + p + calg
	var token = easyMD5(p1) + calg + pid
	// fmt.Println("token", token)
	return token
}

// 创建绑定地址的`URL`
// !!!并没有什么用, 不要用这个函数🙅
//
// 注意!!可能是并不准确的实现方式
func createBindDeviceURL() string {
	return fmt.Sprintf("%v:9002/In0", AuthBaseURL)
}

// Login 登录
//
// 逆向源地址: http://210.22.55.58/a41.js
//
// 		var username, password = "用户名", "密码"
// 		info, _ := inet.Login(username, password)
// 		fmt.Println(info.GetMsg())
//
func Login(username string, password string) (LoginInfo, error) {

	var p = createPassword(password)

	var postData = url.Values{
		"DDDDD":  {username},
		"upass":  {p},
		"R1":     {"0"},
		"R2":     {"1"},
		"para":   {"00"},
		"0MKKey": {"123456"},
		"v6ip":   {""},
	}

	// link: https://stackoverflow.com/a/40006833
	// mockHeader := map[string][]string{
	// 	"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	// 	"Accept-Encoding":           {"gzip, deflate"},
	// 	"Accept-Language":           {"zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7"},
	// 	"Cache-Control":             {"max-age=0"},
	// 	"Content-Type":              {"application/x-www-form-urlencoded"},
	// 	"Origin":                    {"http://210.22.55.58"},
	// 	"Referer":                   {"http://210.22.55.58/0.htm"},
	// 	"Upgrade-Insecure-Requests": {"1"},
	// 	"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"},
	// }

	// var _, e = http.Get(createBindDeviceURL())

	// if e != nil {
	// 	fmt.Println("绑定地址失败", e)
	// 	return e
	// }

	var resp, err = http.PostForm(LoginAPI, postData)

	if err != nil {
		return LoginInfo{}, errors.New("请求登录失败")
	}

	defer resp.Body.Close()

	var body, _ = ioutil.ReadAll(resp.Body)
	var s, _ = gbkToUtf8(body)

	// ===
	// ioutil.WriteFile("dev.html", s, 0644)
	// ===

	info := calljsLoginInfo(s)

	return info, nil

}

// 返回登录结果
func calljsLoginInfo(htmlCodeBytes []byte) LoginInfo {
	jQuery, _ := goquery.NewDocumentFromReader(bytes.NewReader(htmlCodeBytes))
	var script = getJsCode(jQuery)
	VM := js.New()
	VM.RunString(script)
	var m1 = VM.Get("Msg")
	var m2 = VM.Get("msga")

	// 作者注解: 如果找不到全局变量就证明登录成功了/狗头保命
	//
	if m1 == nil && m2 == nil {
		return LoginInfo{
			code: LoginSuccess,
			msg:  "登录成功",
		}
	}
	var msg = strings.TrimSpace(m1.String())
	var msga = strings.TrimSpace(m2.String())
	code, _ := strconv.Atoi(msg)
	var message = "未知"
	switch msga {
	case "5":
		message = "多台设备在线"
		code = 5
		break
	case "1":
		message = "账号密码错误"
		code = 1
		break
	}
	return LoginInfo{
		code: code,    // msg,
		msg:  message, // msga,
	}
}

// HasLogin 判断是否登录
func HasLogin() bool {
	data, err := QueryInfo()
	if err != nil {
		return false
	}
	return !(data.GetCode() == ErrorNoAuthCode)
}

// Logout 注销
func Logout() error {
	_, err := http.Get(LogoutAPI)
	return err
}

// 拿到的`html`格式为`gbk`, 需要转为`utf-8`
//
// 参考: http://mengqi.info/html/2015/201507071345-using-golang-to-convert-text-between-gbk-and-utf-8.html
func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 运营商的js写的不标准, 此函数试图去除一些不标准
//
//
func jsCodeRemoveCommit(data string) string {
	var code = strings.TrimSpace(data)
	var left, right = "<!--", "// -->"
	if strings.HasPrefix(code, left) {
		code = code[len(left):]
	}
	if strings.HasSuffix(code, right) {
		var index = len(code) - len(right)
		code = code[0:index]
	}
	return code
}

// 拿到正确的 `js-code`
//
//
func getJsCode(jQuery *goquery.Document) string {
	var boom = jQuery.Find("SCRIPT")
	var jsCode = boom.Text()
	return jsCodeRemoveCommit(jsCode)
}

// 执行`js`拿到`data
//
//
func calljsGetInfo(jsCode string) (QueryInfoData, error) {
	var code = jsCodeRemoveCommit(jsCode)
	vm := js.New() // 创建engine实例
	var utf8code, _ = gbkToUtf8([]byte(code))
	_, err := vm.RunString(string(utf8code))

	// ===
	// ioutil.WriteFile("x.js", utf8code, 0644)
	// ===

	if err != nil {
		return QueryInfoData{}, errors.New("运行js错误: ")
	}

	// 名称
	var portalname = vm.Get("portalname").String()

	// 未知字段
	// var carrier = vm.Get("carrier").String()

	// 未知字段
	// var portalver = vm.Get("portalver").String()

	// 未知字段
	// var portalid = vm.Get("portalid").String()

	// 已使用时间(分钟)
	var time = strings.TrimSpace(vm.Get("time").String())

	// 流量(mb)
	var flow float64 = 0

	v, e := vm.RunString("flow1/1024+flow3+flow0/1024")

	if e != nil {
		return QueryInfoData{}, errors.New("将流量值转为`int`失败")
	}

	num := v.Export().(string)
	flow, _ = strconv.ParseFloat(num, 64)

	// 未知字段
	// var fsele = vm.Get("fsele").String()

	// 未知字段
	// var fee = vm.Get("fee").String()

	// 未知字段
	// var cvid = vm.Get("cvid").String()

	// 外网映射地址
	var xip = vm.Get("xip").String()

	// 未知字段
	// var pvid = vm.Get("pvid").String()

	// 用户名id
	var uid = vm.Get("uid").String()

	// ipv4 静态ip, 估计是用来鉴权的
	var v4ip = vm.Get("v4ip").String()

	// ipv6 v6ip, 一般没有
	var v6ip = vm.Get("v6ip").String()

	return QueryInfoData{
		code:       ErrorLoginAuthCode,
		Portalname: portalname,
		Time:       time,
		Flow:       flow,
		Xip:        xip,
		UID:        uid,
		V4ip:       v4ip,
		V6ip:       v6ip,
	}, nil
}

// QueryInfo 查询当前信息
func QueryInfo() (QueryInfoData, error) {

	// =======================
	//
	// 直接访问后台地址, 不需要任何鉴权
	//
	// =======================
	jQuery, err := goquery.NewDocument(AuthBaseURL)
	if err != nil {
		return QueryInfoData{}, errors.New("请求登录网管失败")
	}

	// =======================
	//
	// 未登录将会自动跳转到 `/0.htm` 登录界面, 但判断条件为 `title` 字符串为空
	//
	// =======================
	var title = jQuery.Find("title").Text()

	if len(title) == 0 {
		return QueryInfoData{
			code: ErrorNoAuthCode,
		}, nil
	}

	// =======================
	//
	// 我怀疑后台是不是脑阔有问题, 脚本标签居然写成大写的(无意冒犯...)
	//
	// =======================
	var boom = jQuery.Find("script[language=\"JavaScript\"]")

	var jsCode = boom.Text()
	return calljsGetInfo(jsCode)
}
