
// 自动生成模板ComputeNode
package compute
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 算力节点 结构体  ComputeNode
type ComputeNode struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:节点名称;column:name;" binding:"required"`  //名字
  Region  *string `json:"region" form:"region" gorm:"comment:区域;column:region;"`  //区域
  CPU  *int64 `json:"cpu" form:"cpu" gorm:"comment:CPU核心数;column:cpu;"`  //CPU
  Memory  *int64 `json:"memory" form:"memory" gorm:"comment:内存(GB);column:memory;"`  //内存
  SystemDisk  *int64 `json:"systemDisk" form:"systemDisk" gorm:"comment:系统盘容量(GB);column:system_disk;"`  //系统盘容量
  DataDisk  *int64 `json:"dataDisk" form:"dataDisk" gorm:"comment:数据盘容量(GB);column:data_disk;"`  //数据盘容量
  PublicIP  *string `json:"publicIP" form:"publicIP" gorm:"comment:公网IP;column:public_ip;" binding:"required"`  //ip地址公网
  PrivateIP  *string `json:"privateIP" form:"privateIP" gorm:"comment:内网IP;column:private_ip;" binding:"required"`  //ip地址内网
  SshPort  *int64 `json:"sshPort" form:"sshPort" gorm:"default:22;comment:SSH端口;column:ssh_port;" binding:"required"`  //ssh端口
  Username  *string `json:"username" form:"username" gorm:"comment:SSH用户名;column:username;"`  //用户名
  Password  *string `json:"password" form:"password" gorm:"comment:SSH密码;column:password;"`  //密码
  GpuName  *string `json:"gpuName" form:"gpuName" gorm:"comment:GPU名称;column:gpu_name;"`  //显卡名称
  GpuCount  *int64 `json:"gpuCount" form:"gpuCount" gorm:"comment:GPU数量;column:gpu_count;"`  //显卡数量
  DockerEndpoint  *string `json:"dockerEndpoint" form:"dockerEndpoint" gorm:"comment:Docker连接地址;column:docker_endpoint;"`  //docker连接地址
  UseTLS  *bool `json:"useTLS" form:"useTLS" gorm:"default:true;comment:是否使用TLS;column:use_tls;"`  //使用TLS
  CACert  *string `json:"caCert" form:"caCert" gorm:"comment:CA证书;column:ca_cert;type:text;"`  //CA证书
  ClientCert  *string `json:"clientCert" form:"clientCert" gorm:"comment:客户端证书;column:client_cert;type:text;"`  //客户端证书
  ClientKey  *string `json:"clientKey" form:"clientKey" gorm:"comment:客户端私钥;column:client_key;type:text;"`  //客户端私钥
  ShelfStatus  *bool `json:"shelfStatus" form:"shelfStatus" gorm:"default:true;comment:是否上架;column:shelf_status;" binding:"required"`  //是否上架
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`  //备注
}


// TableName 算力节点 ComputeNode自定义表名 compute_nodes
func (ComputeNode) TableName() string {
    return "compute_nodes"
}





