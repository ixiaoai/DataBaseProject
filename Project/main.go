package main

import (
	"project/router"
	"project/utils"
)

func main() {

	utils.InitConfig() //初始化配置文件
	utils.InitMysql()  //初始化数据库

	r := router.Router()
	err := r.Run(":8080")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
