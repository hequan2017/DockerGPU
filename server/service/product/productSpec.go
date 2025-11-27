
package product

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/product"
    productReq "github.com/flipped-aurora/gin-vue-admin/server/model/product/request"
)

type ProductSpecService struct {}
// CreateProductSpec 创建产品规格记录
// Author [yourname](https://github.com/yourname)
func (prodSpecService *ProductSpecService) CreateProductSpec(ctx context.Context, prodSpec *product.ProductSpec) (err error) {
	err = global.GVA_DB.Create(prodSpec).Error
	return err
}

// DeleteProductSpec 删除产品规格记录
// Author [yourname](https://github.com/yourname)
func (prodSpecService *ProductSpecService)DeleteProductSpec(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&product.ProductSpec{},"id = ?",ID).Error
	return err
}

// DeleteProductSpecByIds 批量删除产品规格记录
// Author [yourname](https://github.com/yourname)
func (prodSpecService *ProductSpecService)DeleteProductSpecByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]product.ProductSpec{},"id in ?",IDs).Error
	return err
}

// UpdateProductSpec 更新产品规格记录
// Author [yourname](https://github.com/yourname)
func (prodSpecService *ProductSpecService)UpdateProductSpec(ctx context.Context, prodSpec product.ProductSpec) (err error) {
	err = global.GVA_DB.Model(&product.ProductSpec{}).Where("id = ?",prodSpec.ID).Updates(&prodSpec).Error
	return err
}

// GetProductSpec 根据ID获取产品规格记录
// Author [yourname](https://github.com/yourname)
func (prodSpecService *ProductSpecService)GetProductSpec(ctx context.Context, ID string) (prodSpec product.ProductSpec, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&prodSpec).Error
	return
}
// GetProductSpecInfoList 分页获取产品规格记录
// Author [yourname](https://github.com/yourname)
func (prodSpecService *ProductSpecService)GetProductSpecInfoList(ctx context.Context, info productReq.ProductSpecSearch) (list []product.ProductSpec, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&product.ProductSpec{})
    var prodSpecs []product.ProductSpec
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+ *info.Name+"%")
    }
    if info.GpuModel != nil && *info.GpuModel != "" {
        db = db.Where("gpu_model LIKE ?", "%"+ *info.GpuModel+"%")
    }
    if info.GpuCount != nil {
        db = db.Where("gpu_count >= ?", *info.GpuCount)
    }
    if info.CpuCores != nil {
        db = db.Where("cpu_cores >= ?", *info.CpuCores)
    }
    if info.MemoryGB != nil {
        db = db.Where("memory_gb >= ?", *info.MemoryGB)
    }
    if info.SystemDiskGB != nil {
        db = db.Where("system_disk_gb >= ?", *info.SystemDiskGB)
    }
    if info.DataDiskGB != nil {
        db = db.Where("data_disk_gb >= ?", *info.DataDiskGB)
    }
    if info.PricePerHour != nil {
        db = db.Where("price_per_hour >= ?", *info.PricePerHour)
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

	err = db.Find(&prodSpecs).Error
	return  prodSpecs, total, err
}
func (prodSpecService *ProductSpecService)GetProductSpecPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
