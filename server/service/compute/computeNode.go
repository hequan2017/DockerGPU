
package compute

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/compute"
    computeReq "github.com/flipped-aurora/gin-vue-admin/server/model/compute/request"
)

type ComputeNodeService struct {}
// CreateComputeNode 创建算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) CreateComputeNode(ctx context.Context, cmpNode *compute.ComputeNode) (err error) {
	err = global.GVA_DB.Create(cmpNode).Error
	return err
}

// DeleteComputeNode 删除算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService)DeleteComputeNode(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&compute.ComputeNode{},"id = ?",ID).Error
	return err
}

// DeleteComputeNodeByIds 批量删除算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService)DeleteComputeNodeByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]compute.ComputeNode{},"id in ?",IDs).Error
	return err
}

// UpdateComputeNode 更新算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService)UpdateComputeNode(ctx context.Context, cmpNode compute.ComputeNode) (err error) {
	err = global.GVA_DB.Model(&compute.ComputeNode{}).Where("id = ?",cmpNode.ID).Updates(&cmpNode).Error
	return err
}

// GetComputeNode 根据ID获取算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService)GetComputeNode(ctx context.Context, ID string) (cmpNode compute.ComputeNode, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&cmpNode).Error
	return
}
// GetComputeNodeInfoList 分页获取算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService)GetComputeNodeInfoList(ctx context.Context, info computeReq.ComputeNodeSearch) (list []compute.ComputeNode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&compute.ComputeNode{})
    var cmpNodes []compute.ComputeNode
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+ *info.Name+"%")
    }
    if info.Region != nil && *info.Region != "" {
        db = db.Where("region LIKE ?", "%"+ *info.Region+"%")
    }
    if info.CPU != nil {
        db = db.Where("cpu >= ?", *info.CPU)
    }
    if info.Memory != nil {
        db = db.Where("memory >= ?", *info.Memory)
    }
    if info.SystemDisk != nil {
        db = db.Where("system_disk >= ?", *info.SystemDisk)
    }
    if info.DataDisk != nil {
        db = db.Where("data_disk >= ?", *info.DataDisk)
    }
    if info.PublicIP != nil && *info.PublicIP != "" {
        db = db.Where("public_ip LIKE ?", "%"+ *info.PublicIP+"%")
    }
    if info.PrivateIP != nil && *info.PrivateIP != "" {
        db = db.Where("private_ip LIKE ?", "%"+ *info.PrivateIP+"%")
    }
    if info.SshPort != nil {
        db = db.Where("ssh_port = ?", *info.SshPort)
    }
    if info.Username != nil && *info.Username != "" {
        db = db.Where("username LIKE ?", "%"+ *info.Username+"%")
    }
    if info.Password != nil && *info.Password != "" {
        db = db.Where("password LIKE ?", "%"+ *info.Password+"%")
    }
    if info.GpuName != nil && *info.GpuName != "" {
        db = db.Where("gpu_name LIKE ?", "%"+ *info.GpuName+"%")
    }
    if info.GpuCount != nil {
        db = db.Where("gpu_count >= ?", *info.GpuCount)
    }
    if info.DockerEndpoint != nil && *info.DockerEndpoint != "" {
        db = db.Where("docker_endpoint LIKE ?", "%"+ *info.DockerEndpoint+"%")
    }
    if info.UseTLS != nil {
        db = db.Where("use_tls = ?", *info.UseTLS)
    }
    if info.CACert != "" {
        // TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
    }
    if info.ClientCert != "" {
        // TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
    }
    if info.ClientKey != "" {
        // TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
    }
    if info.ShelfStatus != nil {
        db = db.Where("shelf_status = ?", *info.ShelfStatus)
    }
    if info.Remark != nil && *info.Remark != "" {
        db = db.Where("remark LIKE ?", "%"+ *info.Remark+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&cmpNodes).Error
	return  cmpNodes, total, err
}
func (cmpNodeService *ComputeNodeService)GetComputeNodePublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
