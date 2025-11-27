package instance

import (
	"fmt"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/instance"
	instanceReq "github.com/flipped-aurora/gin-vue-admin/server/model/instance/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type InstanceApi struct{}

// CreateInstance 创建实例管理
// @Tags Instance
// @Summary 创建实例管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body instance.Instance true "创建实例管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /inst/createInstance [post]
func (instApi *InstanceApi) CreateInstance(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var inst instance.Instance
	err := c.ShouldBindJSON(&inst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	uid64 := int64(uid)
	inst.UserID = &uid64
	err = instService.CreateInstance(ctx, &inst)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteInstance 删除实例管理
// @Tags Instance
// @Summary 删除实例管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body instance.Instance true "删除实例管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /inst/deleteInstance [delete]
func (instApi *InstanceApi) DeleteInstance(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	// 非管理员只能删除自己的实例
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		reinst, gerr := instService.GetInstance(ctx, ID)
		if gerr != nil {
			global.GVA_LOG.Error("查询失败!", zap.Error(gerr))
			response.FailWithMessage("删除失败:"+gerr.Error(), c)
			return
		}
		uid := utils.GetUserID(c)
		own := false
		if reinst.UserID != nil {
			own = uint(*reinst.UserID) == uid
		}
		if !own {
			response.FailWithMessage("无权限删除他人实例", c)
			return
		}
	}
	err := instService.DeleteInstance(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteInstanceByIds 批量删除实例管理
// @Tags Instance
// @Summary 批量删除实例管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /inst/deleteInstanceByIds [delete]
func (instApi *InstanceApi) DeleteInstanceByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		// 过滤仅属于当前用户的ID
		uid := utils.GetUserID(c)
		ownIDs := make([]string, 0, len(IDs))
		for _, id := range IDs {
			reinst, _ := instService.GetInstance(ctx, id)
			if reinst.UserID != nil && uint(*reinst.UserID) == uid {
				ownIDs = append(ownIDs, id)
			}
		}
		IDs = ownIDs
		if len(IDs) == 0 {
			response.FailWithMessage("无权限删除所选实例", c)
			return
		}
	}
	err := instService.DeleteInstanceByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// RestartContainer 重启容器
// @Tags Instance
// @Summary 重启实例容器
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "实例ID"
// @Success 200 {object} response.Response{msg=string} "重启成功"
// @Router /inst/restartContainer [post]
func (instApi *InstanceApi) RestartContainer(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("缺少实例ID", c)
		return
	}
	// 权限：非管理员仅可操作自己实例
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		reinst, gerr := instService.GetInstance(ctx, ID)
		if gerr != nil {
			response.FailWithMessage(gerr.Error(), c)
			return
		}
		uid := utils.GetUserID(c)
		if reinst.UserID == nil || uint(*reinst.UserID) != uid {
			response.FailWithMessage("无权限操作他人实例", c)
			return
		}
	}
	if err := instService.RestartContainer(ctx, ID); err != nil {
		global.GVA_LOG.Error("重启失败!", zap.Error(err))
		response.FailWithMessage("重启失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("重启成功", c)
}

// StopContainer 关闭容器
// @Tags Instance
// @Summary 关闭实例容器
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "实例ID"
// @Success 200 {object} response.Response{msg=string} "关闭成功"
// @Router /inst/stopContainer [post]
func (instApi *InstanceApi) StopContainer(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("缺少实例ID", c)
		return
	}
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		reinst, gerr := instService.GetInstance(ctx, ID)
		if gerr != nil {
			response.FailWithMessage(gerr.Error(), c)
			return
		}
		uid := utils.GetUserID(c)
		if reinst.UserID == nil || uint(*reinst.UserID) != uid {
			response.FailWithMessage("无权限操作他人实例", c)
			return
		}
	}
	if err := instService.StopContainer(ctx, ID); err != nil {
		global.GVA_LOG.Error("关闭失败!", zap.Error(err))
		response.FailWithMessage("关闭失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("关闭成功", c)
}

// GetContainerLogs 查看容器日志
// @Tags Instance
// @Summary 查看容器日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "实例ID"
// @Param tail query int false "尾部行数，默认200"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /inst/getContainerLogs [get]
func (instApi *InstanceApi) GetContainerLogs(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	tail := 0
	if t := c.Query("tail"); t != "" {
		_, _ = fmt.Sscanf(t, "%d", &tail)
	}
	if ID == "" {
		response.FailWithMessage("缺少实例ID", c)
		return
	}
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		reinst, gerr := instService.GetInstance(ctx, ID)
		if gerr != nil {
			response.FailWithMessage(gerr.Error(), c)
			return
		}
		uid := utils.GetUserID(c)
		if reinst.UserID == nil || uint(*reinst.UserID) != uid {
			response.FailWithMessage("无权限操作他人实例", c)
			return
		}
	}
	logs, err := instService.GetContainerLogs(ctx, ID, tail)
	if err != nil {
		global.GVA_LOG.Error("获取日志失败!", zap.Error(err))
		response.FailWithMessage("获取日志失败:"+err.Error(), c)
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(logs))
}

// ExecContainerCmd 容器内执行命令
// @Tags Instance
// @Summary 容器内执行命令
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "{"ID":"实例ID","Cmd":["sh","-c","echo hi"]}"
// @Success 200 {object} response.Response{data=object,msg=string} "执行成功"
// @Router /inst/execContainerCmd [post]
func (instApi *InstanceApi) ExecContainerCmd(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		ID  string   `json:"ID"`
		Cmd []string `json:"Cmd"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == "" {
		response.FailWithMessage("缺少实例ID", c)
		return
	}
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		reinst, gerr := instService.GetInstance(ctx, req.ID)
		if gerr != nil {
			response.FailWithMessage(gerr.Error(), c)
			return
		}
		uid := utils.GetUserID(c)
		if reinst.UserID == nil || uint(*reinst.UserID) != uid {
			response.FailWithMessage("无权限操作他人实例", c)
			return
		}
	}
	out, err := instService.ExecContainerCmd(ctx, req.ID, req.Cmd)
	if err != nil {
		global.GVA_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"output": out}, c)
}

// TerminalWS 交互式终端（简版）：通过WebSocket接收命令并返回输出
// @Tags Instance
// @Summary 交互式终端（WebSocket）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "实例ID"
// @Router /inst/terminal [get]
func (instApi *InstanceApi) TerminalWS(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	// 支持通过query传入token，便于浏览器WebSocket握手
	if tok := c.Query("token"); tok != "" {
		j := utils.NewJWT()
		if cl, err := j.ParseToken(tok); err == nil {
			c.Set("claims", &cl)
			utils.SetToken(c, tok, int(cl.ExpiresAt.Unix()-time.Now().Unix()))
		}
	}
	if ID == "" {
		c.String(http.StatusBadRequest, "缺少实例ID")
		return
	}
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		reinst, gerr := instService.GetInstance(ctx, ID)
		if gerr != nil {
			c.String(http.StatusBadRequest, gerr.Error())
			return
		}
		uid := utils.GetUserID(c)
		if reinst.UserID == nil || uint(*reinst.UserID) != uid {
			c.String(http.StatusForbidden, "无权限操作他人实例")
			return
		}
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	hr, shell, closer, err := instService.OpenTerminal(ctx, ID)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte("ERROR: "+err.Error()+"\n"))
		return
	}
	_ = conn.WriteMessage(websocket.TextMessage, []byte("Connected with "+shell+". Type exit to quit.\n"))
	done := make(chan struct{})
	go func() {
		defer func() { close(done) }()
		buf := make([]byte, 8192)
		for {
			n, er := hr.Reader.Read(buf)
			if n > 0 {
				_ = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
			if er != nil {
				return
			}
		}
	}()
	for {
		mt, msg, er := conn.ReadMessage()
		if er != nil {
			break
		}
		if mt == websocket.TextMessage || mt == websocket.BinaryMessage {
			if string(msg) == "exit" {
				break
			}
			_, _ = hr.Conn.Write(msg)
		}
	}
	closer()
}

// UpdateInstance 更新实例管理
// @Tags Instance
// @Summary 更新实例管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body instance.Instance true "更新实例管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /inst/updateInstance [put]
func (instApi *InstanceApi) UpdateInstance(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var inst instance.Instance
	err := c.ShouldBindJSON(&inst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = instService.UpdateInstance(ctx, inst)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindInstance 用id查询实例管理
// @Tags Instance
// @Summary 用id查询实例管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询实例管理"
// @Success 200 {object} response.Response{data=instance.Instance,msg=string} "查询成功"
// @Router /inst/findInstance [get]
func (instApi *InstanceApi) FindInstance(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reinst, err := instService.GetInstance(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reinst, c)
}

// GetInstanceList 分页获取实例管理列表
// @Tags Instance
// @Summary 分页获取实例管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query instanceReq.InstanceSearch true "分页获取实例管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /inst/getInstanceList [get]
func (instApi *InstanceApi) GetInstanceList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo instanceReq.InstanceSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 非管理员仅可查看自身数据
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		uid := utils.GetUserID(c)
		uidInt := int(uid)
		pageInfo.UserID = &uidInt
	}
	list, total, err := instService.GetInstanceInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetInstanceDataSource 获取Instance的数据源
// @Tags Instance
// @Summary 获取Instance的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /inst/getInstanceDataSource [get]
func (instApi *InstanceApi) GetInstanceDataSource(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口为获取数据源定义的数据
	dataSource, err := instService.GetInstanceDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetMatchedComputeNodes 根据规格匹配可用算力节点
// @Tags Instance
// @Summary 根据规格匹配可用算力节点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param specId query uint true "产品规格ID"
// @Success 200 {object} response.Response{data=[]map[string]any,msg=string} "获取成功"
// @Router /inst/getMatchedComputeNodes [get]
func (instApi *InstanceApi) GetMatchedComputeNodes(c *gin.Context) {
	ctx := c.Request.Context()
	specID := c.Query("specId")
	if specID == "" {
		response.FailWithMessage("缺少specId", c)
		return
	}
	options, err := instService.GetMatchedComputeNodes(ctx, specID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithData(options, c)
}

// GetInstancePublic 不需要鉴权的实例管理接口
// @Tags Instance
// @Summary 不需要鉴权的实例管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /inst/getInstancePublic [get]
func (instApi *InstanceApi) GetInstancePublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	instService.GetInstancePublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的实例管理接口信息",
	}, "获取成功", c)
}
