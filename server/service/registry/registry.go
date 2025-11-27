package registry

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/registry"
	registryReq "github.com/flipped-aurora/gin-vue-admin/server/model/registry/request"
)

type RegistryService struct{}

// CreateRegistry 创建镜像库记录
// Author [yourname](https://github.com/yourname)
func (regService *RegistryService) CreateRegistry(ctx context.Context, reg *registry.Registry) (err error) {
	err = global.GVA_DB.Create(reg).Error
	return err
}

// DeleteRegistry 删除镜像库记录
// Author [yourname](https://github.com/yourname)
func (regService *RegistryService) DeleteRegistry(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&registry.Registry{}, "id = ?", ID).Error
	return err
}

// DeleteRegistryByIds 批量删除镜像库记录
// Author [yourname](https://github.com/yourname)
func (regService *RegistryService) DeleteRegistryByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]registry.Registry{}, "id in ?", IDs).Error
	return err
}

// UpdateRegistry 更新镜像库记录
// Author [yourname](https://github.com/yourname)
func (regService *RegistryService) UpdateRegistry(ctx context.Context, reg registry.Registry) (err error) {
	err = global.GVA_DB.Model(&registry.Registry{}).Where("id = ?", reg.ID).Select("*").Updates(&reg).Error
	return err
}

// GetRegistry 根据ID获取镜像库记录
// Author [yourname](https://github.com/yourname)
func (regService *RegistryService) GetRegistry(ctx context.Context, ID string) (reg registry.Registry, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&reg).Error
	return
}

// GetRegistryInfoList 分页获取镜像库记录
// Author [yourname](https://github.com/yourname)
func (regService *RegistryService) GetRegistryInfoList(ctx context.Context, info registryReq.RegistrySearch) (list []registry.Registry, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&registry.Registry{})
	var regs []registry.Registry
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Name != nil && *info.Name != "" {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Address != nil && *info.Address != "" {
		db = db.Where("address LIKE ?", "%"+*info.Address+"%")
	}
	if info.Description != nil && *info.Description != "" {
		db = db.Where("description LIKE ?", "%"+*info.Description+"%")
	}
	if info.Type != "" {
		db = db.Where("type = ?", info.Type)
	}
	if info.ShelfStatus != nil {
		db = db.Where("shelf_status = ?", *info.ShelfStatus)
	}
	if info.Remark != nil && *info.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+*info.Remark+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&regs).Error
	return regs, total, err
}
func (regService *RegistryService) GetRegistryPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
