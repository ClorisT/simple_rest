package api

import (
	"simple_rest/api/controller/dbaccess"
	"simple_rest/api/controller/demo"
	"simple_rest/api/controller/fixme"
	"simple_rest/api/controller/user"
	"simple_rest/api/controller/user/pwd"

	"github.com/gin-gonic/gin"
)

// BindRouting : 綁定路由與操作方法
func BindRouting(router *gin.Engine) {

	// 在路由中區分版本號
	v1 := router.Group("/v1")
	{
		v1.GET("/get", demo.Getting)
		v1.POST("/post", demo.Postting)
		v1.GET("/user", dbaccess.GetUser)
		v1.GET("/deposit", fixme.GetDepostit)
		v1.POST("user/create", user.Creating)
		v1.POST("user/delete", user.Delete)
		v1.POST("user/pwd/change", pwd.Change)
		v1.GET("user/login", user.Login)
	}

}
