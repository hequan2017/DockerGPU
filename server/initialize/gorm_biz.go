package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/compute"
	"github.com/flipped-aurora/gin-vue-admin/server/model/instance"
	"github.com/flipped-aurora/gin-vue-admin/server/model/product"
	"github.com/flipped-aurora/gin-vue-admin/server/model/registry"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(registry.Registry{}, compute.ComputeNode{}, product.ProductSpec{}, instance.Instance{})
	if err != nil {
		return err
	}
	return nil
}
