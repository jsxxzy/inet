package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/jsxxzy/inet"
)

// ========

var (
	DEFAULT_SECTION     = "default"
	DEFAULT_COMMENT     = []byte{'#'}
	DEFAULT_COMMENT_SEM = []byte{';'}
)

type ConfigInterface interface {
	String(key string) string
	Strings(key string) []string
	Bool(key string) (bool, error)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Float64(key string) (float64, error)
	Set(key string, value string) error
}

type Config struct {
	// map is not safe.
	sync.RWMutex
	// Section:key=value
	data map[string]map[string]string
}

// NewConfig create an empty configuration representation.
func NewConfig(confName string) (ConfigInterface, error) {
	c := &Config{
		data: make(map[string]map[string]string),
	}
	err := c.parse(confName)
	return c, err
}

// AddConfig adds a new section->key:value to the configuration.
func (c *Config) AddConfig(section string, option string, value string) bool {
	if section == "" {
		section = DEFAULT_SECTION
	}

	if _, ok := c.data[section]; !ok {
		c.data[section] = make(map[string]string)
	}

	_, ok := c.data[section][option]
	c.data[section][option] = value

	return !ok
}

func (c *Config) parse(fname string) (err error) {
	c.Lock()
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer c.Unlock()
	defer f.Close()

	buf := bufio.NewReader(f)

	var section string
	var lineNum int

	for {
		lineNum++
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if bytes.Equal(line, []byte{}) {
			continue
		} else if err != nil {
			return err
		}

		line = bytes.TrimSpace(line)
		switch {
		case bytes.HasPrefix(line, DEFAULT_COMMENT):
			continue
		case bytes.HasPrefix(line, DEFAULT_COMMENT_SEM):
			continue
		case bytes.HasPrefix(line, []byte{'['}) && bytes.HasSuffix(line, []byte{']'}):
			section = string(line[1 : len(line)-1])
		default:
			optionVal := bytes.SplitN(line, []byte{'='}, 2)
			if len(optionVal) != 2 {
				return fmt.Errorf("parse %s the content error : line %d , %s = ? ", fname, lineNum, optionVal[0])
			}
			option := bytes.TrimSpace(optionVal[0])
			value := bytes.TrimSpace(optionVal[1])
			c.AddConfig(section, strings.ToLower(string(option)), string(value))
		}
	}

	return nil
}

func (c *Config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.get(key))
}

func (c *Config) Int(key string) (int, error) {
	return strconv.Atoi(c.get(key))
}

func (c *Config) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.get(key), 10, 64)
}

func (c *Config) Float64(key string) (float64, error) {
	return strconv.ParseFloat(c.get(key), 64)
}

func (c *Config) String(key string) string {
	return c.get(key)
}

func (c *Config) Strings(key string) []string {
	v := c.get(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ",")
}

func (c *Config) Set(key string, value string) error {
	c.Lock()
	defer c.Unlock()
	if len(key) == 0 {
		return errors.New("key is empty.")
	}

	var (
		section string
		option  string
	)

	keys := strings.Split(strings.ToLower(key), "::")
	if len(keys) >= 2 {
		section = keys[0]
		option = keys[1]
	} else {
		option = keys[0]
	}

	c.AddConfig(section, option, value)
	return nil
}

// section.key or key
func (c *Config) get(key string) string {
	var (
		section string
		option  string
	)

	keys := strings.Split(strings.ToLower(key), "::")

	if len(keys) >= 2 {
		section = keys[0]
		option = keys[1]
	} else {
		section = DEFAULT_SECTION
		option = keys[0]
	}

	if value, ok := c.data[section][option]; ok {
		return value
	}

	return ""
}

// =========

// ConfigFile 配置文件
var ConfigFile = ""

var configFileName = ".inet.conf"

// 流量本地持久化文件名..
var flowFileName = ".inetflow"

// Auth 鉴权
type Auth struct {
	// Username 账号
	Username string
	// password 密码
	Password string
}

// 获取用户`home`目录
func getHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// 初始化配置文件
func initConfigFile() error {
	if exists(ConfigFile) {
		return nil
	}
	var initStr = `
# 请填入账号和密码即可
# https://github.com/jsxxzy/inettray

username = 
password = `
	return ioutil.WriteFile(ConfigFile, []byte(initStr), 0777)
}

// check file/dir exists
//
// https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetConfigFile 获取 `config` 配置文件
func GetConfigFile() (Auth, error) {
	conf, err := NewConfig(ConfigFile) // string(byteData))
	if err != nil {
		return Auth{}, err
	}
	u := conf.String("username")
	p := conf.String("password")
	return Auth{
		Username: u,
		Password: p,
	}, nil
}

func init() {
	homeDir, err := getHomeDir()
	if err != nil {
		panic(err)
	}
	ConfigFile = filepath.Join(homeDir, configFileName)
	initConfigFile()
}

// OpenConfig 打开配置文件
func OpenConfig() error {
	if runtime.GOOS == "windows" {
		tmpRun := exec.Command("notepad", ConfigFile)
		return tmpRun.Run()
	}
	return errors.New("不支持的操作系统")
}

func main() {
	var auth, err = GetConfigFile()
	if err != nil {
		OpenConfig()
		return
	}
	if !inet.HasLogin() {
		inet.Login(auth.Username, auth.Password)
	}
}
