
// 自动生成模板ProductSpec
package product
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 产品规格 结构体  ProductSpec
type ProductSpec struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:规格名称;column:name;" binding:"required"`  //名称
  GpuModel  *string `json:"gpuModel" form:"gpuModel" gorm:"comment:显卡型号;column:gpu_model;" binding:"required"`  //显卡型号
  GpuCount  *int64 `json:"gpuCount" form:"gpuCount" gorm:"comment:显卡数量;column:gpu_count;"`  //显卡数量
  CpuCores  *int64 `json:"cpuCores" form:"cpuCores" gorm:"comment:CPU核心数;column:cpu_cores;"`  //CPU核心数
  MemoryGB  *int64 `json:"memoryGB" form:"memoryGB" gorm:"comment:内存(GB);column:memory_gb;"`  //内存(GB)
  SystemDiskGB  *int64 `json:"systemDiskGB" form:"systemDiskGB" gorm:"comment:系统盘容量(GB);column:system_disk_gb;"`  //系统盘容量(GB)
  DataDiskGB  *int64 `json:"dataDiskGB" form:"dataDiskGB" gorm:"comment:数据盘容量(GB);column:data_disk_gb;"`  //数据盘容量(GB)
  PricePerHour  *float64 `json:"pricePerHour" form:"pricePerHour" gorm:"comment:价格/小时;column:price_per_hour;"`  //价格/小时
  ShelfStatus  *bool `json:"shelfStatus" form:"shelfStatus" gorm:"default:true;comment:是否上架;column:shelf_status;" binding:"required"`  //是否上架
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`  //备注
}


// TableName 产品规格 ProductSpec自定义表名 product_specs
func (ProductSpec) TableName() string {
    return "product_specs"
}





