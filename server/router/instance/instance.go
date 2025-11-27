package instance

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type InstanceRouter struct {}

// InitInstanceRouter 初始化 实例管理 路由信息
func (s *InstanceRouter) InitInstanceRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	instRouter := Router.Group("inst").Use(middleware.OperationRecord())
	instRouterWithoutRecord := Router.Group("inst")
	instRouterWithoutAuth := PublicRouter.Group("inst")
	{
		instRouter.POST("createInstance", instApi.CreateInstance)   // 新建实例管理
		instRouter.DELETE("deleteInstance", instApi.DeleteInstance) // 删除实例管理
		instRouter.DELETE("deleteInstanceByIds", instApi.DeleteInstanceByIds) // 批量删除实例管理
		instRouter.PUT("updateInstance", instApi.UpdateInstance)    // 更新实例管理
	}
	{
		instRouterWithoutRecord.GET("findInstance", instApi.FindInstance)        // 根据ID获取实例管理
		instRouterWithoutRecord.GET("getInstanceList", instApi.GetInstanceList)  // 获取实例管理列表
	}
	{
	    instRouterWithoutAuth.GET("getInstanceDataSource", instApi.GetInstanceDataSource)  // 获取实例管理数据源
	    instRouterWithoutAuth.GET("getInstancePublic", instApi.GetInstancePublic)  // 实例管理开放接口
	}
}
