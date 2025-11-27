package product

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ ProductSpecRouter }

var prodSpecApi = api.ApiGroupApp.ProductApiGroup.ProductSpecApi
