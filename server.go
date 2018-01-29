package main

import (
	"account-auth-service/api"
	"account-auth-service/logic"
	"flag"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"gitlab.ucloudadmin.com/peter/ugo/common"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var (
	confFile = flag.String("c", "", "configuration file,json format")
)

func main() {
	runtime.GOMAXPROCS(0)
	// 解释命令行选项
	common.ProcessOptions()

	if err := common.LoadConfigFromFile(*confFile); err != nil {
		fmt.Println("Load Config File fail,", err)
		return
	}
	common.DumpConfigContent()
	//初始化日志
	dir, _ := common.GetConfigByKey("log.dir")
	name, _ := common.GetConfigByKey("log.name")
	//size, _ := common.GetConfigByKey("log.maxSize")
	//files, _ := common.GetConfigByKey("log.maxFiles")
	level, _ := common.GetConfigByKey("log.level")

	logger.SetRollingDaily(dir.(string), name.(string))
	//logger.SetRollingFile(dir.(string), name.(string), int32(files.(float64)), int64(size.(float64)), logger.MB)
	logger.SetLevel(logger.LEVEL(level.(float64)))

	//初始化数据库
	err := logic.NewDBRegist()
	if err != nil {
		fmt.Println("init DB error:", err)
		return
	}

	// 获取监听ip
	ip, err := common.GetConfigByKey("addr.ip")
	if err != nil {
		fmt.Println("can not get listen ip:", err)
		return
	}
	ipAddr, _ := ip.(string)
	// 获取监听端口
	port, err := common.GetConfigByKey("addr.port")
	if err != nil {
		fmt.Println("can not get listen port")
		return
	}

	fmt.Println(time.Now().Unix())
	addr := string(ipAddr) + ":" + strconv.Itoa(int(port.(float64)))

	fmt.Println(addr)
	logger.Debug("This is a test")

	err = http.ListenAndServe(addr, api.WsContainer)
	if err != nil {
		fmt.Println(err)
	}
}
