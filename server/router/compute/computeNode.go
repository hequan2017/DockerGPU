package compute

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ComputeNodeRouter struct{}

// InitComputeNodeRouter 初始化 算力节点 路由信息
func (s *ComputeNodeRouter) InitComputeNodeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	cmpNodeRouter := Router.Group("cmpNode").Use(middleware.OperationRecord())
	cmpNodeRouterWithoutRecord := Router.Group("cmpNode")
	cmpNodeRouterWithoutAuth := PublicRouter.Group("cmpNode")
	{
		cmpNodeRouter.POST("createComputeNode", cmpNodeApi.CreateComputeNode)             // 新建算力节点
		cmpNodeRouter.DELETE("deleteComputeNode", cmpNodeApi.DeleteComputeNode)           // 删除算力节点
		cmpNodeRouter.DELETE("deleteComputeNodeByIds", cmpNodeApi.DeleteComputeNodeByIds) // 批量删除算力节点
		cmpNodeRouter.PUT("updateComputeNode", cmpNodeApi.UpdateComputeNode)              // 更新算力节点
	}
	{
		cmpNodeRouterWithoutRecord.GET("findComputeNode", cmpNodeApi.FindComputeNode)       // 根据ID获取算力节点
		cmpNodeRouterWithoutRecord.GET("getComputeNodeList", cmpNodeApi.GetComputeNodeList) // 获取算力节点列表
		cmpNodeRouterWithoutRecord.GET("checkDockerTLS", cmpNodeApi.CheckDockerTLS)         // 校验Docker端点TLS
		cmpNodeRouterWithoutRecord.GET("checkAllDockerTLS", cmpNodeApi.CheckAllDockerTLS)   // 批量校验所有上架节点
	}
	{
		cmpNodeRouterWithoutAuth.GET("getComputeNodePublic", cmpNodeApi.GetComputeNodePublic) // 算力节点开放接口
	}
}
