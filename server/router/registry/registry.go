package registry

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type RegistryRouter struct {}

// InitRegistryRouter 初始化 镜像库 路由信息
func (s *RegistryRouter) InitRegistryRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	regRouter := Router.Group("reg").Use(middleware.OperationRecord())
	regRouterWithoutRecord := Router.Group("reg")
	regRouterWithoutAuth := PublicRouter.Group("reg")
	{
		regRouter.POST("createRegistry", regApi.CreateRegistry)   // 新建镜像库
		regRouter.DELETE("deleteRegistry", regApi.DeleteRegistry) // 删除镜像库
		regRouter.DELETE("deleteRegistryByIds", regApi.DeleteRegistryByIds) // 批量删除镜像库
		regRouter.PUT("updateRegistry", regApi.UpdateRegistry)    // 更新镜像库
	}
	{
		regRouterWithoutRecord.GET("findRegistry", regApi.FindRegistry)        // 根据ID获取镜像库
		regRouterWithoutRecord.GET("getRegistryList", regApi.GetRegistryList)  // 获取镜像库列表
	}
	{
	    regRouterWithoutAuth.GET("getRegistryPublic", regApi.GetRegistryPublic)  // 镜像库开放接口
	}
}
