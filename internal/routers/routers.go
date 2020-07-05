package routers

import (
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"server_gdraw/api"
	"server_gdraw/internal/middleware"
	"server_gdraw/internal/routers/v1/client"
	"server_gdraw/internal/routers/v1/ws"
	"server_gdraw/internal/services"
	"server_gdraw/pkg/resp"
)

func Init() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type,Uid",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.Group("/").Use(middleware.Ws.HttpUpgrader()).GET("/", ws.Ws.Gdraw)

	apiV1OpenApiGdraw := r.Group("v1/gdraw/openapi") //.Use(middleware.NewMiddleware().VerifyToken())
	{
		apiV1OpenApiGdraw.POST("/user/auth", client.GdrawUserController.UserAuth)
		apiV1OpenApiGdraw.POST("/user/regist", client.GdrawUserController.UserRegist)
		apiV1OpenApiGdraw.POST("/user/login", client.GdrawUserController.Login)
		apiV1OpenApiGdraw.GET("/vcode/get", client.GdrawVcodeController.Get)
		apiV1OpenApiGdraw.GET("/token/verify", client.GdrawUserController.TokenVerify)

		apiV1OpenApiGdraw.GET("/question/refresh", func(context *gin.Context) {
			rk := api.Rds.GetQuestionIdsKey()
			api.Rds.Get().Del(rk)
			ret, err := services.GdrawQuestioner.GetQuestionIds()
			if err != nil {
				resp.RespInternalErr(context, err.Error())
				return
			}
			resp.RespOk(context, ret)
		})
	}

	apiV1ApiGdraw := r.Group("v1/gdraw/api").Use(services.GdrawUserer.VerifyToken())
	{
		apiV1ApiGdraw.GET("/room/list", client.GdrawController.GetRoomList)
		apiV1ApiGdraw.GET("/index", client.GdrawController.GetIndexData)
		apiV1ApiGdraw.GET("/rank/list", client.GdrawController.GetRankData)
		apiV1ApiGdraw.POST("/feedback", client.GdrawController.CreateFeedback)
		apiV1ApiGdraw.POST("/share/create", client.GdrawController.CreateShare)
	}
	return r
}
