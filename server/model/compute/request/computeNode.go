
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ComputeNodeSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      Name  *string `json:"name" form:"name"` 
      Region  *string `json:"region" form:"region"` 
      CPU  *int `json:"cpu" form:"cpu"` 
      Memory  *int `json:"memory" form:"memory"` 
      SystemDisk  *int `json:"systemDisk" form:"systemDisk"` 
      DataDisk  *int `json:"dataDisk" form:"dataDisk"` 
      PublicIP  *string `json:"publicIP" form:"publicIP"` 
      PrivateIP  *string `json:"privateIP" form:"privateIP"` 
      SshPort  *int `json:"sshPort" form:"sshPort"` 
      Username  *string `json:"username" form:"username"` 
      Password  *string `json:"password" form:"password"` 
      GpuName  *string `json:"gpuName" form:"gpuName"` 
      GpuCount  *int `json:"gpuCount" form:"gpuCount"` 
      DockerEndpoint  *string `json:"dockerEndpoint" form:"dockerEndpoint"` 
      UseTLS  *bool `json:"useTLS" form:"useTLS"` 
      CACert  string `json:"caCert" form:"caCert"` 
      ClientCert  string `json:"clientCert" form:"clientCert"` 
      ClientKey  string `json:"clientKey" form:"clientKey"` 
      ShelfStatus  *bool `json:"shelfStatus" form:"shelfStatus"` 
      Remark  *string `json:"remark" form:"remark"` 
    request.PageInfo
}
