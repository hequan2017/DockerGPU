package product

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProductSpecRouter struct {}

// InitProductSpecRouter 初始化 产品规格 路由信息
func (s *ProductSpecRouter) InitProductSpecRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	prodSpecRouter := Router.Group("prodSpec").Use(middleware.OperationRecord())
	prodSpecRouterWithoutRecord := Router.Group("prodSpec")
	prodSpecRouterWithoutAuth := PublicRouter.Group("prodSpec")
	{
		prodSpecRouter.POST("createProductSpec", prodSpecApi.CreateProductSpec)   // 新建产品规格
		prodSpecRouter.DELETE("deleteProductSpec", prodSpecApi.DeleteProductSpec) // 删除产品规格
		prodSpecRouter.DELETE("deleteProductSpecByIds", prodSpecApi.DeleteProductSpecByIds) // 批量删除产品规格
		prodSpecRouter.PUT("updateProductSpec", prodSpecApi.UpdateProductSpec)    // 更新产品规格
	}
	{
		prodSpecRouterWithoutRecord.GET("findProductSpec", prodSpecApi.FindProductSpec)        // 根据ID获取产品规格
		prodSpecRouterWithoutRecord.GET("getProductSpecList", prodSpecApi.GetProductSpecList)  // 获取产品规格列表
	}
	{
	    prodSpecRouterWithoutAuth.GET("getProductSpecPublic", prodSpecApi.GetProductSpecPublic)  // 产品规格开放接口
	}
}
