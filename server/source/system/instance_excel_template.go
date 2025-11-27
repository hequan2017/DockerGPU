package system

import (
    "context"
    sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    sysService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
    "github.com/pkg/errors"
    "gorm.io/gorm"
)

type initInstanceExcelTemplate struct{}

const initOrderInstanceExcelTemplate = initOrderExcelTemplate + 1

func init() {
    sysService.RegisterInit(initOrderInstanceExcelTemplate, &initInstanceExcelTemplate{})
}

func (i *initInstanceExcelTemplate) InitializerName() string {
    return "instance_export_template"
}

func (i *initInstanceExcelTemplate) MigrateTable(ctx context.Context) (context.Context, error) {
    db, ok := ctx.Value("db").(*gorm.DB)
    if !ok {
        return ctx, sysService.ErrMissingDBContext
    }
    return ctx, db.AutoMigrate(&sysModel.SysExportTemplate{}, &sysModel.Condition{}, &sysModel.JoinTemplate{})
}

func (i *initInstanceExcelTemplate) TableCreated(ctx context.Context) bool {
    db, ok := ctx.Value("db").(*gorm.DB)
    if !ok {
        return false
    }
    return db.Migrator().HasTable(&sysModel.SysExportTemplate{})
}

func (i *initInstanceExcelTemplate) InitializeData(ctx context.Context) (context.Context, error) {
    db, ok := ctx.Value("db").(*gorm.DB)
    if !ok {
        return ctx, sysService.ErrMissingDBContext
    }
    var count int64
    db.Model(&sysModel.SysExportTemplate{}).Where("template_id = ?", "instance_Instance").Count(&count)
    if count > 0 {
        return ctx, nil
    }
    entity := sysModel.SysExportTemplate{
        Name:       "instance_Instance",
        TableName:  "instances",
        TemplateID: "instance_Instance",
        TemplateInfo: `{"id":"ID","user_id":"用户ID","server_id":"节点ID","template_id":"模版ID","image_id":"镜像ID","container_id":"容器ID","name":"实例名称","status":"状态","remark":"备注","created_at":"创建时间","updated_at":"更新时间"}`,
    }
    if err := db.Create(&entity).Error; err != nil {
        return ctx, errors.Wrap(err, "sys_export_templates")
    }
    next := context.WithValue(ctx, i.InitializerName(), entity)
    return next, nil
}

func (i *initInstanceExcelTemplate) DataInserted(ctx context.Context) bool {
    db, ok := ctx.Value("db").(*gorm.DB)
    if !ok {
        return false
    }
    var count int64
    db.Model(&sysModel.SysExportTemplate{}).Where("template_id = ?", "instance_Instance").Count(&count)
    return count > 0
}

