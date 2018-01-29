package logic

import (
	"account-auth-service/model"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.ucloudadmin.com/peter/ugo/common"
)

type Req struct {
	Action     string      `json:"Action"`
	Account_id int         `json:"account_id"`
	Data       interface{} `json:"Data"`
}

type Resp struct {
	Ret         int32       `json:"RetCode"`
	Message     string      `json:"Message"`
	Total_count uint32      `json:"TotalCount"`
	Data        interface{} `json:"Data"`
}

//注册数据库表
var TABLES = []interface{}{
	&model.User_auth{},
}

// type MysqlConfig struct {
// 	Id       string
// 	Type     interface{}
// 	Selector interface{}
// 	Data     []interface{}
// 	Limit    int
// 	Offset   int
// }

func NewDBRegist() error {
	db, err := ConnectMysql()
	if err != nil {
		fmt.Println("connect db error :")
		return err
	}
	defer db.Close()

	err = db.AutoMigrate(TABLES...).Error
	if err != nil {
		fmt.Println("connect db error :")
		return err
	}

	return nil

}

func ConnectMysql() (db *gorm.DB, err error) {
	mysqlAddr, err := common.GetConfigByKey("mysql.addr")
	if err != nil {
		fmt.Println(err)
		return
	}

	mysqlUser, err := common.GetConfigByKey("mysql.user")
	if err != nil {
		fmt.Println(err)
		return
	}

	mysqlPwd, err := common.GetConfigByKey("mysql.pwd")
	if err != nil {
		fmt.Println(err)
		return
	}

	mysqldb, err := common.GetConfigByKey("mysql.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	mysqlurl := mysqlUser.(string) + ":" + mysqlPwd.(string) + "@tcp(" + mysqlAddr.(string) + ")/" + mysqldb.(string) + "?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", mysqlurl)

	if err != nil {
		fmt.Println("connect db error")
		fmt.Println(err)
	}
	return
}

func ConnectRedis() *redis.Client {
	redisaddr, err := common.GetConfigByKey("redis.addr")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisaddr.(string),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}
