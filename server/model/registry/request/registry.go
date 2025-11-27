
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type RegistrySearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      Name  *string `json:"name" form:"name"` 
      Address  *string `json:"address" form:"address"` 
      Description  *string `json:"description" form:"description"` 
      Type  string `json:"type" form:"type"` 
      ShelfStatus  *bool `json:"shelfStatus" form:"shelfStatus"` 
      Remark  *string `json:"remark" form:"remark"` 
    request.PageInfo
}
