package compute

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ ComputeNodeRouter }

var cmpNodeApi = api.ApiGroupApp.ComputeApiGroup.ComputeNodeApi
