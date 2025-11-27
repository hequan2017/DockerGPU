package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/compute"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/instance"
	"github.com/flipped-aurora/gin-vue-admin/server/service/product"
	"github.com/flipped-aurora/gin-vue-admin/server/service/registry"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	RegistryServiceGroup registry.ServiceGroup
	ComputeServiceGroup  compute.ServiceGroup
	ProductServiceGroup  product.ServiceGroup
	InstanceServiceGroup instance.ServiceGroup
}
