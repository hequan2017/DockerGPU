package compute

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/compute"
	computeReq "github.com/flipped-aurora/gin-vue-admin/server/model/compute/request"
)

type ComputeNodeService struct{}

// CreateComputeNode 创建算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) CreateComputeNode(ctx context.Context, cmpNode *compute.ComputeNode) (err error) {
	err = global.GVA_DB.Create(cmpNode).Error
	return err
}

// DeleteComputeNode 删除算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) DeleteComputeNode(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&compute.ComputeNode{}, "id = ?", ID).Error
	return err
}

// DeleteComputeNodeByIds 批量删除算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) DeleteComputeNodeByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]compute.ComputeNode{}, "id in ?", IDs).Error
	return err
}

// UpdateComputeNode 更新算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) UpdateComputeNode(ctx context.Context, cmpNode compute.ComputeNode) (err error) {
	err = global.GVA_DB.Model(&compute.ComputeNode{}).Where("id = ?", cmpNode.ID).Updates(&cmpNode).Error
	return err
}

// GetComputeNode 根据ID获取算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) GetComputeNode(ctx context.Context, ID string) (cmpNode compute.ComputeNode, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&cmpNode).Error
	return
}

// GetComputeNodeInfoList 分页获取算力节点记录
// Author [yourname](https://github.com/yourname)
func (cmpNodeService *ComputeNodeService) GetComputeNodeInfoList(ctx context.Context, info computeReq.ComputeNodeSearch) (list []compute.ComputeNode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&compute.ComputeNode{})
	var cmpNodes []compute.ComputeNode
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Name != nil && *info.Name != "" {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Region != nil && *info.Region != "" {
		db = db.Where("region LIKE ?", "%"+*info.Region+"%")
	}
	if info.CPU != nil {
		db = db.Where("cpu >= ?", *info.CPU)
	}
	if info.Memory != nil {
		db = db.Where("memory >= ?", *info.Memory)
	}
	if info.SystemDisk != nil {
		db = db.Where("system_disk >= ?", *info.SystemDisk)
	}
	if info.DataDisk != nil {
		db = db.Where("data_disk >= ?", *info.DataDisk)
	}
	if info.PublicIP != nil && *info.PublicIP != "" {
		db = db.Where("public_ip LIKE ?", "%"+*info.PublicIP+"%")
	}
	if info.PrivateIP != nil && *info.PrivateIP != "" {
		db = db.Where("private_ip LIKE ?", "%"+*info.PrivateIP+"%")
	}
	if info.SshPort != nil {
		db = db.Where("ssh_port = ?", *info.SshPort)
	}
	if info.Username != nil && *info.Username != "" {
		db = db.Where("username LIKE ?", "%"+*info.Username+"%")
	}
	if info.Password != nil && *info.Password != "" {
		db = db.Where("password LIKE ?", "%"+*info.Password+"%")
	}
	if info.GpuName != nil && *info.GpuName != "" {
		db = db.Where("gpu_name LIKE ?", "%"+*info.GpuName+"%")
	}
	if info.GpuCount != nil {
		db = db.Where("gpu_count >= ?", *info.GpuCount)
	}
	if info.DockerEndpoint != nil && *info.DockerEndpoint != "" {
		db = db.Where("docker_endpoint LIKE ?", "%"+*info.DockerEndpoint+"%")
	}
	if info.UseTLS != nil {
		db = db.Where("use_tls = ?", *info.UseTLS)
	}
	if info.CACert != "" {
		// TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
	}
	if info.ClientCert != "" {
		// TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
	}
	if info.ClientKey != "" {
		// TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
	}
	if info.ShelfStatus != nil {
		db = db.Where("shelf_status = ?", *info.ShelfStatus)
	}
	if info.Remark != nil && *info.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+*info.Remark+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&cmpNodes).Error
	return cmpNodes, total, err
}
func (cmpNodeService *ComputeNodeService) GetComputeNodePublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

func (cmpNodeService *ComputeNodeService) CheckDockerTLS(ctx context.Context, ID string) (map[string]any, error) {
	var node compute.ComputeNode
	if err := global.GVA_DB.Where("id = ?", ID).First(&node).Error; err != nil {
		return map[string]any{"ok": false, "message": "node_not_found"}, err
	}
	endpoint := valueOrString(node.DockerEndpoint)
	useTLS := valueOrBool(node.UseTLS)
	ca := valueOrString(node.CACert)
	cert := valueOrString(node.ClientCert)
	key := valueOrString(node.ClientKey)
	if useTLS && (cert == "" || key == "") {
		return map[string]any{"ok": false, "message": "client_cert_required"}, nil
	}
	client := buildHTTPClient(useTLS, ca, cert, key)
	base := normalizeDockerEndpoint(endpoint, useTLS)
	req, _ := http.NewRequest("GET", base+"/_ping", nil)
	resp, err := client.Do(req)
	if err != nil {
		return map[string]any{"ok": false, "message": err.Error()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return map[string]any{"ok": false, "message": resp.Status}, nil
	}
	return map[string]any{"ok": true, "message": "OK"}, nil
}

func (cmpNodeService *ComputeNodeService) CheckAllDockerTLS(ctx context.Context) ([]map[string]any, error) {
	nodes := make([]compute.ComputeNode, 0)
	if err := global.GVA_DB.Model(&compute.ComputeNode{}).Where("deleted_at IS NULL AND shelf_status = ?", true).Find(&nodes).Error; err != nil {
		return nil, err
	}
	res := make([]map[string]any, 0, len(nodes))
	for _, n := range nodes {
		id := fmt.Sprintf("%d", n.ID)
		r, err := cmpNodeService.CheckDockerTLS(ctx, id)
		name := ""
		if n.Name != nil {
			name = *n.Name
		}
		if err != nil {
			res = append(res, map[string]any{"id": n.ID, "name": name, "ok": false, "message": err.Error()})
			continue
		}
		r["id"] = n.ID
		r["name"] = name
		res = append(res, r)
	}
	return res, nil
}

func buildHTTPClient(useTLS bool, ca, cert, key string) *http.Client {
	if !useTLS {
		return &http.Client{}
	}
	cfg := &tls.Config{MinVersion: tls.VersionTLS12}
	if ca != "" {
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(ca)); ok {
			cfg.RootCAs = pool
		}
	} else {
		cfg.InsecureSkipVerify = true
	}
	if cert != "" && key != "" {
		if pair, err := tls.X509KeyPair([]byte(cert), []byte(key)); err == nil {
			cfg.Certificates = []tls.Certificate{pair}
		}
	}
	return &http.Client{Transport: &http.Transport{TLSClientConfig: cfg}}
}

func normalizeDockerEndpoint(endpoint string, useTLS bool) string {
	if strings.HasPrefix(endpoint, "tcp://") {
		if useTLS {
			return "https://" + strings.TrimPrefix(endpoint, "tcp://")
		}
		return "http://" + strings.TrimPrefix(endpoint, "tcp://")
	}
	return endpoint
}

func valueOrBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}
func valueOrString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
