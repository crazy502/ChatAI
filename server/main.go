package main

import (
	"fmt"
	"log"

	"server/common/mysql"
	"server/common/rabbitmq"
	"server/common/redis"
	"server/config"
	"server/router"
)

func StartServer(addr string, port int) error {
	r := router.InitRouter()
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

func main() {
	conf := config.GetConfig()
	host := conf.Host
	port := conf.Port

	if err := mysql.InitMysql(); err != nil {
		log.Println("InitMysql error , " + err.Error())
		return
	}

	redis.Init()
	log.Println("redis init success  ")

	if err := rabbitmq.InitRabbitMQ(); err != nil {
		log.Println("rabbitmq init degraded mode:", err)
	} else {
		log.Println("rabbitmq init success  ")
	}

	err := StartServer(host, port)
	if err != nil {
		panic(err)
	}
}
