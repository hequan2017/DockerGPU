package request

// IDRequest 通用ID请求结构体
type IDRequest struct {
	ID string `json:"id" form:"id" uri:"id" binding:"required"`
}

// RestartContainerRequest 重启容器请求结构体
type RestartContainerRequest struct {
	IDRequest
}

// StopContainerRequest 关闭容器请求结构体
type StopContainerRequest struct {
	IDRequest
}

// GetContainerLogsRequest 查看容器日志请求结构体
type GetContainerLogsRequest struct {
	ID     string `json:"id" form:"id" uri:"id" binding:"required"`
	Tail   int    `json:"tail" form:"tail" binding:"min=1,max=10000"`
	Follow bool   `json:"follow" form:"follow"`
}

// ExecContainerCmdRequest 执行容器命令请求结构体
type ExecContainerCmdRequest struct {
	ID   string   `json:"id" form:"id" uri:"id" binding:"required"`
	Cmd  []string `json:"cmd" form:"cmd" binding:"required"`
	Tty  bool     `json:"tty" form:"tty"`
}

// UpdateInstanceStatusRequest 更新实例状态请求结构体
type UpdateInstanceStatusRequest struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=running stopped error creating"`
}
