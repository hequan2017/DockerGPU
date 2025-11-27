
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ProductSpecSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      Name  *string `json:"name" form:"name"` 
      GpuModel  *string `json:"gpuModel" form:"gpuModel"` 
      GpuCount  *int `json:"gpuCount" form:"gpuCount"` 
      CpuCores  *int `json:"cpuCores" form:"cpuCores"` 
      MemoryGB  *int `json:"memoryGB" form:"memoryGB"` 
      SystemDiskGB  *int `json:"systemDiskGB" form:"systemDiskGB"` 
      DataDiskGB  *int `json:"dataDiskGB" form:"dataDiskGB"` 
      PricePerHour  *float64 `json:"pricePerHour" form:"pricePerHour"` 
      ShelfStatus  *bool `json:"shelfStatus" form:"shelfStatus"` 
      Remark  *string `json:"remark" form:"remark"` 
    request.PageInfo
}
