
// 自动生成模板Instance
package instance
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 实例管理 结构体  Instance
type Instance struct {
    global.GVA_MODEL
  UserID  *int64 `json:"userId" form:"userId" gorm:"comment:所属用户ID;column:user_id;"`  //所属用户ID
  ServerID  *int64 `json:"serverId" form:"serverId" gorm:"comment:来源算力节点ID;column:server_id;"`  //来源服务器ID
  TemplateID  *int64 `json:"templateId" form:"templateId" gorm:"comment:来源产品规格ID;column:template_id;"`  //来源模版ID
  ImageID  *int64 `json:"imageId" form:"imageId" gorm:"comment:来源镜像库ID;column:image_id;"`  //来源镜像ID
  ContainerID  *string `json:"containerId" form:"containerId" gorm:"comment:Docker容器ID;column:container_id;"`  //Docker容器ID
  Name  *string `json:"name" form:"name" gorm:"comment:实例名称;column:name;"`  //实例名称
  Status  string `json:"status" form:"status" gorm:"comment:实例状态;column:status;type:varchar(32);default:creating;"`  //状态
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`  //备注
}


// TableName 实例管理 Instance自定义表名 instances
func (Instance) TableName() string {
    return "instances"
}




