package product

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ ProductSpecApi }

var prodSpecService = service.ServiceGroupApp.ProductServiceGroup.ProductSpecService
