
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type InstanceSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      UserID  *int `json:"userId" form:"userId"` 
      ServerID  *int `json:"serverId" form:"serverId"` 
      TemplateID  *int `json:"templateId" form:"templateId"` 
      ImageID  *int `json:"imageId" form:"imageId"` 
      ContainerID  *string `json:"containerId" form:"containerId"` 
      Name  *string `json:"name" form:"name"` 
      Status  string `json:"status" form:"status"` 
      Remark  *string `json:"remark" form:"remark"` 
    request.PageInfo
}
