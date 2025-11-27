
package instance

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/instance"
    instanceReq "github.com/flipped-aurora/gin-vue-admin/server/model/instance/request"
)

type InstanceService struct {}
// CreateInstance 创建实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) CreateInstance(ctx context.Context, inst *instance.Instance) (err error) {
	err = global.GVA_DB.Create(inst).Error
	return err
}

// DeleteInstance 删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService)DeleteInstance(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&instance.Instance{},"id = ?",ID).Error
	return err
}

// DeleteInstanceByIds 批量删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService)DeleteInstanceByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]instance.Instance{},"id in ?",IDs).Error
	return err
}

// UpdateInstance 更新实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService)UpdateInstance(ctx context.Context, inst instance.Instance) (err error) {
	err = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?",inst.ID).Updates(&inst).Error
	return err
}

// GetInstance 根据ID获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService)GetInstance(ctx context.Context, ID string) (inst instance.Instance, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inst).Error
	return
}
// GetInstanceInfoList 分页获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService)GetInstanceInfoList(ctx context.Context, info instanceReq.InstanceSearch) (list []instance.Instance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&instance.Instance{})
    var insts []instance.Instance
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.UserID != nil {
        db = db.Where("user_id = ?", *info.UserID)
    }
    if info.ServerID != nil {
        db = db.Where("server_id = ?", *info.ServerID)
    }
    if info.TemplateID != nil {
        db = db.Where("template_id = ?", *info.TemplateID)
    }
    if info.ImageID != nil {
        db = db.Where("image_id = ?", *info.ImageID)
    }
    if info.ContainerID != nil && *info.ContainerID != "" {
        db = db.Where("container_id LIKE ?", "%"+ *info.ContainerID+"%")
    }
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+ *info.Name+"%")
    }
    if info.Status != "" {
        db = db.Where("status = ?", info.Status)
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

	err = db.Find(&insts).Error
	return  insts, total, err
}
func (instService *InstanceService)GetInstanceDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   imageId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("registries").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&imageId)
	   res["imageId"] = imageId
	   serverId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("compute_nodes").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&serverId)
	   res["serverId"] = serverId
	   templateId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("product_specs").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&templateId)
	   res["templateId"] = templateId
	   userId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("sys_users").Where("deleted_at IS NULL").Select("username as label,id as value").Scan(&userId)
	   res["userId"] = userId
	return
}
func (instService *InstanceService)GetInstancePublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
