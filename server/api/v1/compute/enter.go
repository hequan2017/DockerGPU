package compute

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ ComputeNodeApi }

var cmpNodeService = service.ServiceGroupApp.ComputeServiceGroup.ComputeNodeService
