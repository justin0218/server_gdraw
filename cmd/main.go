package main

import (
	"fmt"
	"server_gdraw/api"
	"server_gdraw/configs"
	"server_gdraw/internal/routers"
	"server_gdraw/internal/services"
)

func main() {
	//	var str = `西瓜、水果，2个字
	//香蕉、水果，2个字
	//凳子、家具，2个字
	//菠萝、水果，2个字
	//橘子、水果，2个字
	//电话机、家用电器，3个字
	//草莓、水果，2个字
	//袜子、服饰，2个字
	//苹果、水果，2个字
	//牙膏、个人清洁用品，2个字
	//梨、水果，1个字
	//樱桃、水果，2个字
	//猕猴桃、水果，3个字
	//桃子、水果，2个字
	//葡萄、水果，2个字
	//芒果、水果，2个字
	//梳子、个人用品，2个字
	//桂圆、水果，2个字
	//荔枝、水果，2个字
	//榴莲、水果，2个字
	//火车、交通工具，2个字
	//脸盆、生活用品，2个字`
	//	for _, val := range strings.Split(str, "\n") {
	//		in := models.GdrawQuestion{Question: strings.Split(val, "、")[0], QuestionTips: strings.Split(val, "、")[1]}
	//		api.Mysql.Get().Table("gdraw_questions").Omit("create_time", "update_time").Create(&in)
	//	}
	//	return
	api.Log.Get().Debug("starting...")
	jobs := services.NewJob()
	go jobs.GdrawData()
	err := api.Rds.Get().Ping().Err()
	if err != nil {
		panic(err)
	}
	api.Mysql.Get()
	fmt.Println("server started!!!")
	err = routers.Init().Run(fmt.Sprintf(":%d", configs.Dft.Get().Http.Port))
	if err != nil {
		panic(err)
	}
}
