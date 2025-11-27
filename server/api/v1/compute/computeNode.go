package compute

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/compute"
    computeReq "github.com/flipped-aurora/gin-vue-admin/server/model/compute/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type ComputeNodeApi struct {}



// CreateComputeNode 创建算力节点
// @Tags ComputeNode
// @Summary 创建算力节点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body compute.ComputeNode true "创建算力节点"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cmpNode/createComputeNode [post]
func (cmpNodeApi *ComputeNodeApi) CreateComputeNode(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var cmpNode compute.ComputeNode
	err := c.ShouldBindJSON(&cmpNode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cmpNodeService.CreateComputeNode(ctx,&cmpNode)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteComputeNode 删除算力节点
// @Tags ComputeNode
// @Summary 删除算力节点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body compute.ComputeNode true "删除算力节点"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cmpNode/deleteComputeNode [delete]
func (cmpNodeApi *ComputeNodeApi) DeleteComputeNode(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := cmpNodeService.DeleteComputeNode(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteComputeNodeByIds 批量删除算力节点
// @Tags ComputeNode
// @Summary 批量删除算力节点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /cmpNode/deleteComputeNodeByIds [delete]
func (cmpNodeApi *ComputeNodeApi) DeleteComputeNodeByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := cmpNodeService.DeleteComputeNodeByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateComputeNode 更新算力节点
// @Tags ComputeNode
// @Summary 更新算力节点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body compute.ComputeNode true "更新算力节点"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cmpNode/updateComputeNode [put]
func (cmpNodeApi *ComputeNodeApi) UpdateComputeNode(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var cmpNode compute.ComputeNode
	err := c.ShouldBindJSON(&cmpNode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cmpNodeService.UpdateComputeNode(ctx,cmpNode)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindComputeNode 用id查询算力节点
// @Tags ComputeNode
// @Summary 用id查询算力节点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询算力节点"
// @Success 200 {object} response.Response{data=compute.ComputeNode,msg=string} "查询成功"
// @Router /cmpNode/findComputeNode [get]
func (cmpNodeApi *ComputeNodeApi) FindComputeNode(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	recmpNode, err := cmpNodeService.GetComputeNode(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(recmpNode, c)
}
// GetComputeNodeList 分页获取算力节点列表
// @Tags ComputeNode
// @Summary 分页获取算力节点列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query computeReq.ComputeNodeSearch true "分页获取算力节点列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /cmpNode/getComputeNodeList [get]
func (cmpNodeApi *ComputeNodeApi) GetComputeNodeList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo computeReq.ComputeNodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := cmpNodeService.GetComputeNodeInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetComputeNodePublic 不需要鉴权的算力节点接口
// @Tags ComputeNode
// @Summary 不需要鉴权的算力节点接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cmpNode/getComputeNodePublic [get]
func (cmpNodeApi *ComputeNodeApi) GetComputeNodePublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    cmpNodeService.GetComputeNodePublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的算力节点接口信息",
    }, "获取成功", c)
}
