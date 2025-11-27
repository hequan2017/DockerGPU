package registry

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ RegistryRouter }

var regApi = api.ApiGroupApp.RegistryApiGroup.RegistryApi
