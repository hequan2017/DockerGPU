package registry

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ RegistryApi }

var regService = service.ServiceGroupApp.RegistryServiceGroup.RegistryService
