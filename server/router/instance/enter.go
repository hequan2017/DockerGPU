package instance

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ InstanceRouter }

var instApi = api.ApiGroupApp.InstanceApiGroup.InstanceApi
