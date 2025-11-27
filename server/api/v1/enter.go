package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/compute"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/instance"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/product"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/registry"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	RegistryApiGroup registry.ApiGroup
	ComputeApiGroup  compute.ApiGroup
	ProductApiGroup  product.ApiGroup
	InstanceApiGroup instance.ApiGroup
}
