// 自动生成模板Registry
package registry

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 镜像库 结构体  Registry
type Registry struct {
	global.GVA_MODEL
	Name        *string `json:"name" form:"name" gorm:"comment:镜像库名称;column:name;" binding:"required"`                 //名字
	Address     *string `json:"address" form:"address" gorm:"comment:镜像仓库地址;column:address;" binding:"required"`       //地址
	Description *string `json:"description" form:"description" gorm:"comment:描述信息;column:description;"`                //描述
	Type        string  `json:"type" form:"type" gorm:"comment:仓库类型;column:type;type:varchar(50);" binding:"required"` //类型
	ShelfStatus bool    `json:"shelfStatus" form:"shelfStatus" gorm:"default:true;comment:是否上架;column:shelf_status;"`  //是否上架
	Remark      *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`                                 //备注
}

// TableName 镜像库 Registry自定义表名 registries
func (Registry) TableName() string {
	return "registries"
}
