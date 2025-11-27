import service from '@/utils/request'
// @Tags Registry
// @Summary 创建镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Registry true "创建镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /reg/createRegistry [post]
export const createRegistry = (data) => {
  return service({
    url: '/reg/createRegistry',
    method: 'post',
    data
  })
}

// @Tags Registry
// @Summary 删除镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Registry true "删除镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /reg/deleteRegistry [delete]
export const deleteRegistry = (params) => {
  return service({
    url: '/reg/deleteRegistry',
    method: 'delete',
    params
  })
}

// @Tags Registry
// @Summary 批量删除镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /reg/deleteRegistry [delete]
export const deleteRegistryByIds = (params) => {
  return service({
    url: '/reg/deleteRegistryByIds',
    method: 'delete',
    params
  })
}

// @Tags Registry
// @Summary 更新镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Registry true "更新镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /reg/updateRegistry [put]
export const updateRegistry = (data) => {
  return service({
    url: '/reg/updateRegistry',
    method: 'put',
    data
  })
}

// @Tags Registry
// @Summary 用id查询镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Registry true "用id查询镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /reg/findRegistry [get]
export const findRegistry = (params) => {
  return service({
    url: '/reg/findRegistry',
    method: 'get',
    params
  })
}

// @Tags Registry
// @Summary 分页获取镜像库列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取镜像库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /reg/getRegistryList [get]
export const getRegistryList = (params) => {
  return service({
    url: '/reg/getRegistryList',
    method: 'get',
    params
  })
}

// @Tags Registry
// @Summary 不需要鉴权的镜像库接口
// @Accept application/json
// @Produce application/json
// @Param data query registryReq.RegistrySearch true "分页获取镜像库列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /reg/getRegistryPublic [get]
export const getRegistryPublic = () => {
  return service({
    url: '/reg/getRegistryPublic',
    method: 'get',
  })
}
