# DockerGPU
> 基于Docker的GPU 算力租赁系统

## 已开发功能模块

- 镜像库
  - 字段：名称、地址、描述、类型、是否上架(bool，默认上架)、备注
  - 后端：`server/model/registry/registry.go`；自动迁移与路由已注册
  - 前端：`web/src/view/registry/registry/registry.vue`，列表/表单，`是否上架`为开关控件
  - 兼容性修正：移除数据库枚举写法，采用布尔字段，避免 MySQL 语法问题

- 算力节点
  - 字段：名称(必填)、区域、CPU、内存、系统盘容量、数据盘容量、公网IP(必填)、内网IP(必填)、SSH端口(必填，默认22)、用户名、密码、显卡名称、显卡数量、Docker连接地址、使用TLS(bool，默认启用)、CA证书(TEXT)、客户端证书(TEXT)、客户端私钥(TEXT)、是否上架(bool，默认上架)、备注
  - 后端：`server/model/compute/computenode.go`；自动迁移与路由已注册
  - 前端：`web/src/view/compute/computeNode/computeNode.vue`，布尔项使用开关控件；证书为富文本编辑/查看

- 产品规格
  - 字段：名称(必填)、显卡型号(必填)、显卡数量、CPU核心数、内存(GB)、系统盘容量(GB)、数据盘容量(GB)、价格/小时(float64)、是否上架(bool，默认上架)、备注
  - 后端：`server/model/product/productSpec.go`；自动迁移与路由已注册
  - 前端：`web/src/view/product/productSpec/productSpec.vue`，价格使用数值输入，`是否上架`为开关控件

- 实例管理
  - 字段：所属用户ID(后端创建时自动填写)、来源服务器ID(算力节点)、来源模版ID(产品规格)、来源镜像ID(镜像库)、Docker容器ID(后端创建后回填)、实例名称、状态(创建中/运行中/停止/异常)、备注
  - 状态字典：`instance_status`（创建中/运行中/停止/异常）
  - 后端：`server/model/instance/instance.go`；创建接口自动从认证信息写入用户ID；状态使用字符串并默认 `creating`
  - 前端：`web/src/view/instance/instance/instance.vue`，来源ID下拉，状态字典选择；显示标签统一为“用户/节点/模版/镜像”

## 路由与迁移

- 路由组注册：`server/initialize/router_biz.go`
- 自动迁移：`server/initialize/gorm_biz.go`

## 前端运行

- 开发模式：`web/package.json` 中 `dev` 脚本 `vite --host --mode development`

## 说明

- 布尔字段统一以 `bool` 类型持久化，前端以开关控件展示
- 证书内容使用 `TEXT` 列，前端提供富文本编辑/查看组件
- 实例创建时，所属用户ID自动填充为当前登录用户；列表/搜索仍支持按用户筛选
