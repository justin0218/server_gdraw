package routers

import (
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"server_gdraw/internal/controllers"
	"server_gdraw/internal/middleware"
	"server_gdraw/store"
)

func Init() *gin.Engine {
	r := gin.Default()
	gin.SetMode(new(store.Config).Get().Runmode)
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type,Uid",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.GET("/health", func(context *gin.Context) {
		context.JSON(200, map[string]string{"msg": "ok"})
		return
	})
	wsController := new(controllers.Ws)
	r.Group("/").Use(middleware.HttpUpgrader()).GET("/", wsController.Gdraw)
	apiV1OpenApiGdraw := r.Group("v1/gdraw/openapi") //.Use(middleware.NewMiddleware().VerifyToken())
	{
		userController := new(controllers.User)
		apiV1OpenApiGdraw.POST("/user/auth", userController.UserAuth)
		apiV1OpenApiGdraw.GET("/token/verify", userController.TokenVerify)

		vcodeController := new(controllers.Vcode)
		apiV1OpenApiGdraw.GET("/vcode/get", vcodeController.Get)

		//apiV1OpenApiGdraw.GET("/question/refresh", func(context *gin.Context) {
		//	rk := api.Rds.GetQuestionIdsKey()
		//	api.Rds.Get().Del(rk)
		//	ret, err := services.GdrawQuestioner.GetQuestionIds()
		//	if err != nil {
		//		resp.RespInternalErr(context, err.Error())
		//		return
		//	}
		//	resp.RespOk(context, ret)
		//})
	}

	apiV1ApiGdraw := r.Group("v1/gdraw/api").Use(middleware.VerifyToken())
	{
		roomController := new(controllers.Room)
		apiV1ApiGdraw.GET("/room/list", roomController.GetRoomList)
		apiV1ApiGdraw.GET("/index", roomController.GetIndexData)
		apiV1ApiGdraw.GET("/rank/list", roomController.GetRankData)
		apiV1ApiGdraw.POST("/share/create", roomController.CreateShare)
	}
	return r
}
