package instance

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
    "time"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/compute"
    "github.com/flipped-aurora/gin-vue-admin/server/model/instance"
    instanceReq "github.com/flipped-aurora/gin-vue-admin/server/model/instance/request"
)

type InstanceService struct{}

// CreateInstance 创建实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) CreateInstance(ctx context.Context, inst *instance.Instance) (err error) {
	// 先创建记录，状态为 creating
	if inst.Status == "" {
		inst.Status = "creating"
	}
	if err = global.GVA_DB.Create(inst).Error; err != nil {
		return err
	}

	// 读取节点、规格、镜像信息
	var node compute.ComputeNode
	if inst.ServerID == nil {
		_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"status": "error"}).Error
		return fmt.Errorf("缺少节点信息，请选择节点后再提交")
	}
	if inst.TemplateID == nil {
		_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"status": "error"}).Error
		return fmt.Errorf("缺少产品规格信息，请选择规格后再提交")
	}
	if inst.ImageID == nil {
		_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"status": "error"}).Error
		return fmt.Errorf("缺少镜像信息，请选择镜像后再提交")
	}
	if err = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error; err != nil {
		// 更新状态为 error
		_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Update("status", "error").Error
		return err
	}

	type specRow struct {
		GpuModel     *string
		GpuCount     *int64
		CpuCores     *int64
		MemoryGB     *int64
		SystemDiskGB *int64
		DataDiskGB   *int64
	}
	var spec specRow
	if inst.TemplateID != nil {
		if err = global.GVA_DB.Table("product_specs").Where("id = ?", *inst.TemplateID).Select("gpu_model,gpu_count,cpu_cores,memory_gb,system_disk_gb,data_disk_gb").Scan(&spec).Error; err != nil {
			_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Update("status", "error").Error
			return err
		}
	}

	type registryRow struct{ Address *string }
	var reg registryRow
	if inst.ImageID != nil {
		if err = global.GVA_DB.Table("registries").Where("id = ?", *inst.ImageID).Select("address").Scan(&reg).Error; err != nil {
			_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Update("status", "error").Error
			return err
		}
	}

	// 构造 Docker API 请求参数
	image := ""
	if reg.Address != nil {
		image = *reg.Address
	}
	name := ""
	if inst.Name != nil {
		name = *inst.Name
	} else {
		name = fmt.Sprintf("inst-%d", inst.ID)
	}

	nanoCPUs := int64(0)
	if spec.CpuCores != nil && *spec.CpuCores > 0 {
		nanoCPUs = *spec.CpuCores * 1_000_000_000
	}
	memBytes := int64(0)
	if spec.MemoryGB != nil && *spec.MemoryGB > 0 {
		memBytes = *spec.MemoryGB * 1024 * 1024 * 1024
	}
	gpuCount := int64(0)
	gpuModel := ""
	if spec.GpuCount != nil {
		gpuCount = *spec.GpuCount
	}
	if spec.GpuModel != nil {
		gpuModel = *spec.GpuModel
	}

	// 生成 HostConfig 及 DeviceRequests（NVIDIA GPU）
	hostConfig := map[string]any{
		"NanoCPUs": nanoCPUs,
		"Memory":   memBytes,
	}
	if spec.SystemDiskGB != nil && *spec.SystemDiskGB > 0 {
		hostConfig["StorageOpt"] = map[string]string{
			"size": fmt.Sprintf("%dG", *spec.SystemDiskGB),
		}
	}
	if gpuCount > 0 {
		hostConfig["DeviceRequests"] = []map[string]any{
			{
				"Count":        gpuCount,
				"Driver":       "nvidia",
				"Capabilities": [][]string{{"gpu"}},
			},
		}
	}
	// 磁盘映射：使用命名卷并挂载到容器，容量信息通过Labels记录（Docker不原生限制卷容量）
	mounts := []map[string]any{}
	if spec.SystemDiskGB != nil && *spec.SystemDiskGB > 0 {
		mounts = append(mounts, map[string]any{
			"Type":   "volume",
			"Target": "/system",
			"Source": fmt.Sprintf("inst-%d-system", inst.ID),
		})
	}
	if spec.DataDiskGB != nil && *spec.DataDiskGB > 0 {
		mounts = append(mounts, map[string]any{
			"Type":   "volume",
			"Target": "/data",
			"Source": fmt.Sprintf("inst-%d-data", inst.ID),
		})
	}
	if len(mounts) > 0 {
		hostConfig["Mounts"] = mounts
	}

	body := map[string]any{
		"Image": image,
		"Env": []string{
			fmt.Sprintf("GPU_MODEL=%s", gpuModel),
			fmt.Sprintf("GPU_COUNT=%d", gpuCount),
		},
		"HostConfig": hostConfig,
		"Labels": map[string]string{
			"spec.cpu_cores":      fmt.Sprintf("%d", nanoCPUs/1_000_000_000),
			"spec.memory_gb":      fmt.Sprintf("%d", memBytes/1024/1024/1024),
			"spec.system_disk_gb": fmt.Sprintf("%d", valueOrZero(spec.SystemDiskGB)),
			"spec.data_disk_gb":   fmt.Sprintf("%d", valueOrZero(spec.DataDiskGB)),
			"spec.gpu_model":      gpuModel,
			"spec.gpu_count":      fmt.Sprintf("%d", gpuCount),
		},
	}

	// 调用 Docker API 创建并启动容器
	endpoint := valueOrString(node.DockerEndpoint)
	useTLS := valueOrBool(node.UseTLS)
	ca := valueOrString(node.CACert)
	cert := valueOrString(node.ClientCert)
	key := valueOrString(node.ClientKey)
	if useTLS && (cert == "" || key == "") {
		_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"status": "error"}).Error
		return fmt.Errorf("Docker端点要求TLS客户端证书，请在算力节点配置填写ClientCert/ClientKey")
	}
	if useTLS {
		if _, perr := tls.X509KeyPair([]byte(cert), []byte(key)); perr != nil {
			_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"status": "error"}).Error
			return fmt.Errorf("TLS客户端证书/私钥解析失败: %v", perr)
		}
	}
	containerID, derr := dockerCreateAndStart(endpoint, useTLS, ca, cert, key, name, body)
	if derr != nil {
		_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"status": "error"}).Error
		return derr
	}

	// 回填 containerId 与状态
	cid := containerID
	_ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(map[string]any{"container_id": cid, "status": "running"}).Error
	return nil
}

func dockerCreateAndStart(endpoint string, useTLS bool, ca, cert, key, name string, body map[string]any) (string, error) {
	base := normalizeDockerEndpoint(endpoint, useTLS)
	// create
	createURL := fmt.Sprintf("%s/containers/create?name=%s", base, url.QueryEscape(name))
	payload, _ := json.Marshal(body)
	client := buildHTTPClient(useTLS, ca, cert, key)
	req, _ := http.NewRequest("POST", createURL, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		var emsg struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &emsg)
		// 当存储驱动不支持rootfs大小限制时，移除StorageOpt后重试一次
		if strings.Contains(emsg.Message, "supported only for overlay over xfs") || strings.Contains(string(b), "supported only for overlay over xfs") {
			if host, ok := body["HostConfig"].(map[string]any); ok {
				if _, has := host["StorageOpt"]; has {
					delete(host, "StorageOpt")
					body["HostConfig"] = host
					payload2, _ := json.Marshal(body)
					req2, _ := http.NewRequest("POST", createURL, bytes.NewReader(payload2))
					req2.Header.Set("Content-Type", "application/json")
					resp2, err2 := client.Do(req2)
					if err2 != nil {
						return "", err2
					}
					defer resp2.Body.Close()
					if resp2.StatusCode >= 300 {
						b2, _ := io.ReadAll(resp2.Body)
						var emsg2 struct {
							Message string `json:"message"`
						}
						_ = json.Unmarshal(b2, &emsg2)
						if emsg2.Message != "" {
							return "", fmt.Errorf("docker create failed: %s - %s", resp2.Status, emsg2.Message)
						}
						return "", fmt.Errorf("docker create failed: %s - %s", resp2.Status, string(b2))
					}
					var createRes2 struct {
						ID string `json:"Id"`
					}
					if err := json.NewDecoder(resp2.Body).Decode(&createRes2); err != nil {
						return "", err
					}
					// start after retry success
					startURL := fmt.Sprintf("%s/containers/%s/start", base, createRes2.ID)
					startReq, _ := http.NewRequest("POST", startURL, nil)
					startResp, err := client.Do(startReq)
					if err != nil {
						return "", err
					}
					defer startResp.Body.Close()
					if startResp.StatusCode >= 300 && startResp.StatusCode != 204 {
						b, _ := io.ReadAll(startResp.Body)
						var emsg struct {
							Message string `json:"message"`
						}
						_ = json.Unmarshal(b, &emsg)
						if emsg.Message != "" {
							return "", fmt.Errorf("docker start failed: %s - %s", startResp.Status, emsg.Message)
						}
						return "", fmt.Errorf("docker start failed: %s - %s", startResp.Status, string(b))
					}
					return createRes2.ID, nil
				}
			}
		}
		if emsg.Message != "" {
			return "", fmt.Errorf("docker create failed: %s - %s", resp.Status, emsg.Message)
		}
		return "", fmt.Errorf("docker create failed: %s - %s", resp.Status, string(b))
	}
	var createRes struct {
		ID string `json:"Id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createRes); err != nil {
		return "", err
	}

	// start
	startURL := fmt.Sprintf("%s/containers/%s/start", base, createRes.ID)
	startReq, _ := http.NewRequest("POST", startURL, nil)
	startResp, err := client.Do(startReq)
	if err != nil {
		return "", err
	}
	defer startResp.Body.Close()
	if startResp.StatusCode >= 300 && startResp.StatusCode != 204 {
		b, _ := io.ReadAll(startResp.Body)
		var emsg struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &emsg)
		if emsg.Message != "" {
			return "", fmt.Errorf("docker start failed: %s - %s", startResp.Status, emsg.Message)
		}
		return "", fmt.Errorf("docker start failed: %s - %s", startResp.Status, string(b))
	}
	return createRes.ID, nil
}

func dockerDeleteContainerAndVolumes(endpoint string, useTLS bool, ca, cert, key, containerID string, volumes []string) error {
	if endpoint == "" {
		return nil
	}
	base := normalizeDockerEndpoint(endpoint, useTLS)
	client := buildHTTPClient(useTLS, ca, cert, key)
	if containerID != "" {
		delURL := fmt.Sprintf("%s/containers/%s?v=true&force=true", base, url.PathEscape(containerID))
		req, _ := http.NewRequest("DELETE", delURL, nil)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 && resp.StatusCode != 404 {
			b, _ := io.ReadAll(resp.Body)
			var emsg struct {
				Message string `json:"message"`
			}
			_ = json.Unmarshal(b, &emsg)
			if emsg.Message != "" {
				return fmt.Errorf("docker delete container failed: %s - %s", resp.Status, emsg.Message)
			}
			return fmt.Errorf("docker delete container failed: %s - %s", resp.Status, string(b))
		}
	}
	for _, v := range volumes {
		if v == "" {
			continue
		}
		volURL := fmt.Sprintf("%s/volumes/%s", base, url.PathEscape(v))
		req, _ := http.NewRequest("DELETE", volURL, nil)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 && resp.StatusCode != 404 {
			b, _ := io.ReadAll(resp.Body)
			var emsg struct {
				Message string `json:"message"`
			}
			_ = json.Unmarshal(b, &emsg)
			if emsg.Message != "" {
				return fmt.Errorf("docker delete volume failed: %s - %s", resp.Status, emsg.Message)
			}
			return fmt.Errorf("docker delete volume failed: %s - %s", resp.Status, string(b))
		}
	}
	return nil
}

func dockerRestartContainer(endpoint string, useTLS bool, ca, cert, key, containerID string) error {
	if endpoint == "" || containerID == "" {
		return fmt.Errorf("缺少Docker端点或容器ID")
	}
	base := normalizeDockerEndpoint(endpoint, useTLS)
	client := buildHTTPClient(useTLS, ca, cert, key)
	url := fmt.Sprintf("%s/containers/%s/restart", base, url.PathEscape(containerID))
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		var emsg struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &emsg)
		if emsg.Message != "" {
			return fmt.Errorf("docker restart failed: %s - %s", resp.Status, emsg.Message)
		}
		return fmt.Errorf("docker restart failed: %s - %s", resp.Status, string(b))
	}
	return nil
}

func dockerStopContainer(endpoint string, useTLS bool, ca, cert, key, containerID string) error {
	if endpoint == "" || containerID == "" {
		return fmt.Errorf("缺少Docker端点或容器ID")
	}
	base := normalizeDockerEndpoint(endpoint, useTLS)
	client := buildHTTPClient(useTLS, ca, cert, key)
	url := fmt.Sprintf("%s/containers/%s/stop", base, url.PathEscape(containerID))
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != 304 { // 304: already stopped
		b, _ := io.ReadAll(resp.Body)
		var emsg struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &emsg)
		if emsg.Message != "" {
			return fmt.Errorf("docker stop failed: %s - %s", resp.Status, emsg.Message)
		}
		return fmt.Errorf("docker stop failed: %s - %s", resp.Status, string(b))
	}
	return nil
}

func dockerContainerLogs(endpoint string, useTLS bool, ca, cert, key, containerID string, tail int) (string, error) {
	if endpoint == "" || containerID == "" {
		return "", fmt.Errorf("缺少Docker端点或容器ID")
	}
	base := normalizeDockerEndpoint(endpoint, useTLS)
	client := buildHTTPClient(useTLS, ca, cert, key)
	if tail <= 0 {
		tail = 200
	}
	url := fmt.Sprintf("%s/containers/%s/logs?stdout=1&stderr=1&tail=%d", base, url.PathEscape(containerID), tail)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		var emsg struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &emsg)
		if emsg.Message != "" {
			return "", fmt.Errorf("docker logs failed: %s - %s", resp.Status, emsg.Message)
		}
		return "", fmt.Errorf("docker logs failed: %s - %s", resp.Status, string(b))
	}
	b, _ := io.ReadAll(resp.Body)
	return string(b), nil
}

func dockerExecRun(endpoint string, useTLS bool, ca, cert, key, containerID string, cmd []string) (string, error) {
	if endpoint == "" || containerID == "" {
		return "", fmt.Errorf("缺少Docker端点或容器ID")
	}
	if len(cmd) == 0 {
		cmd = []string{"/bin/sh", "-c", "echo 'no cmd'"}
	}
	base := normalizeDockerEndpoint(endpoint, useTLS)
	client := buildHTTPClient(useTLS, ca, cert, key)
	// create exec
	body := map[string]any{
		"AttachStdin":  false,
		"AttachStdout": true,
		"AttachStderr": true,
		"Tty":          true,
		"Cmd":          cmd,
	}
	payload, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/%s/exec", base, url.PathEscape(containerID)), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create exec failed: %s - %s", resp.Status, string(b))
	}
	var execRes struct {
		ID string `json:"Id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&execRes); err != nil {
		return "", err
	}
	// start exec (Tty true for raw stream)
	startBody := map[string]any{"Detach": false, "Tty": true}
	payload2, _ := json.Marshal(startBody)
	req2, _ := http.NewRequest("POST", fmt.Sprintf("%s/exec/%s/start", base, url.PathEscape(execRes.ID)), bytes.NewReader(payload2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode >= 300 {
		b, _ := io.ReadAll(resp2.Body)
		return "", fmt.Errorf("start exec failed: %s - %s", resp2.Status, string(b))
	}
	out, _ := io.ReadAll(resp2.Body)
	return string(out), nil
}

func dockerInspectStatus(endpoint string, useTLS bool, ca, cert, key, containerID string) (string, error) {
    if endpoint == "" || containerID == "" {
        return "", fmt.Errorf("缺少Docker端点或容器ID")
    }
    base := normalizeDockerEndpoint(endpoint, useTLS)
    client := buildHTTPClient(useTLS, ca, cert, key)
    url := fmt.Sprintf("%s/containers/%s/json", base, url.PathEscape(containerID))
    req, _ := http.NewRequest("GET", url, nil)
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 {
        b, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("docker inspect failed: %s - %s", resp.Status, string(b))
    }
    var obj struct {
        State struct {
            Status string `json:"Status"`
        } `json:"State"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
        return "", err
    }
    return obj.State.Status, nil
}

func newDockerCLI(endpoint string, useTLS bool, ca, cert, key string) (*client.Client, error) {
    base := normalizeDockerEndpoint(endpoint, useTLS)
    tr := &http.Transport{}
    if useTLS {
        pool := x509.NewCertPool()
        if ca != "" {
            _ = pool.AppendCertsFromPEM([]byte(ca))
        }
        var certs []tls.Certificate
        if cert != "" && key != "" {
            pair, err := tls.X509KeyPair([]byte(cert), []byte(key))
            if err != nil {
                return nil, err
            }
            certs = append(certs, pair)
        }
        tr.TLSClientConfig = &tls.Config{RootCAs: pool, Certificates: certs, MinVersion: tls.VersionTLS12}
    }
    httpClient := &http.Client{Transport: tr}
    return client.NewClientWithOpts(client.WithHost(base), client.WithHTTPClient(httpClient), client.WithAPIVersionNegotiation())
}

func dockerExecTTYAttach(ctx context.Context, cli *client.Client, containerID string, shell string) (*types.HijackedResponse, error) {
    cfg := types.ExecConfig{AttachStdin: true, AttachStdout: true, AttachStderr: true, Tty: true, Cmd: []string{shell}}
    created, err := cli.ContainerExecCreate(ctx, containerID, cfg)
    if err != nil {
        return nil, err
    }
    hr, err := cli.ContainerExecAttach(ctx, created.ID, types.ExecStartCheck{Tty: true})
    if err != nil {
        return nil, err
    }
    return &hr, nil
}

func (instService *InstanceService) OpenTerminal(ctx context.Context, ID string) (*types.HijackedResponse, string, func(), error) {
    var inst instance.Instance
    if err := global.GVA_DB.Where("id = ?", ID).First(&inst).Error; err != nil {
        return nil, "", nil, err
    }
    var node compute.ComputeNode
    if inst.ServerID != nil {
        _ = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error
    }
    cli, err := newDockerCLI(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey))
    if err != nil {
        return nil, "", nil, err
    }
    containerID := valueOrString(inst.ContainerID)
    shell := "/bin/bash"
    hr, err := dockerExecTTYAttach(ctx, cli, containerID, shell)
    if err != nil {
        shell = "/bin/sh"
        hr, err = dockerExecTTYAttach(ctx, cli, containerID, shell)
        if err != nil {
            _ = cli.Close()
            return nil, "", nil, err
        }
    }
    closer := func() { hr.Close(); _ = cli.Close() }
    return hr, shell, closer, nil
}

func buildHTTPClient(useTLS bool, ca, cert, key string) *http.Client {
	if !useTLS {
		return &http.Client{Timeout: 60 * time.Second}
	}
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	hasCA := false
	if ca != "" {
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM([]byte(ca)); ok {
			tlsConfig.RootCAs = roots
			hasCA = true
		}
	}
	if cert != "" && key != "" {
		if pair, err := tls.X509KeyPair([]byte(cert), []byte(key)); err == nil {
			tlsConfig.Certificates = []tls.Certificate{pair}
		}
	}
	// 当未提供CA时，允许跳过证书校验（自签证书场景）
	if !hasCA {
		tlsConfig.InsecureSkipVerify = true
	}
	return &http.Client{Timeout: 60 * time.Second, Transport: &http.Transport{TLSClientConfig: tlsConfig}}
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

func valueOrZero(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
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

// DeleteInstance 删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) DeleteInstance(ctx context.Context, ID string) (err error) {
	var inst instance.Instance
	if e := global.GVA_DB.Where("id = ?", ID).First(&inst).Error; e != nil {
		return e
	}
	var node compute.ComputeNode
	if inst.ServerID != nil {
		_ = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error
	}
	endpoint := valueOrString(node.DockerEndpoint)
	useTLS := valueOrBool(node.UseTLS)
	ca := valueOrString(node.CACert)
	cert := valueOrString(node.ClientCert)
	key := valueOrString(node.ClientKey)
	volumes := []string{
		fmt.Sprintf("inst-%d-system", inst.ID),
		fmt.Sprintf("inst-%d-data", inst.ID),
	}
	if derr := dockerDeleteContainerAndVolumes(endpoint, useTLS, ca, cert, key, valueOrString(inst.ContainerID), volumes); derr != nil {
		return derr
	}
	err = global.GVA_DB.Delete(&instance.Instance{}, "id = ?", ID).Error
	return err
}

// DeleteInstanceByIds 批量删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) DeleteInstanceByIds(ctx context.Context, IDs []string) (err error) {
	for _, id := range IDs {
		if e := instService.DeleteInstance(ctx, id); e != nil {
			return e
		}
	}
	return nil
}

// UpdateInstance 更新实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) UpdateInstance(ctx context.Context, inst instance.Instance) (err error) {
	err = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Updates(&inst).Error
	return err
}

// GetInstance 根据ID获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) GetInstance(ctx context.Context, ID string) (inst instance.Instance, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inst).Error
	return
}

// GetInstanceInfoList 分页获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) GetInstanceInfoList(ctx context.Context, info instanceReq.InstanceSearch) (list []instance.Instance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&instance.Instance{})
	var insts []instance.Instance
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.UserID != nil {
		db = db.Where("user_id = ?", *info.UserID)
	}
	if info.ServerID != nil {
		db = db.Where("server_id = ?", *info.ServerID)
	}
	if info.TemplateID != nil {
		db = db.Where("template_id = ?", *info.TemplateID)
	}
	if info.ImageID != nil {
		db = db.Where("image_id = ?", *info.ImageID)
	}
	if info.ContainerID != nil && *info.ContainerID != "" {
		db = db.Where("container_id LIKE ?", "%"+*info.ContainerID+"%")
	}
	if info.Name != nil && *info.Name != "" {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
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

    err = db.Find(&insts).Error
    if err != nil {
        return insts, total, err
    }
    for i := range insts {
        var node compute.ComputeNode
        if insts[i].ServerID == nil || insts[i].ContainerID == nil || valueOrString(insts[i].ContainerID) == "" {
            continue
        }
        _ = global.GVA_DB.Where("id = ?", *insts[i].ServerID).First(&node).Error
        status, serr := dockerInspectStatus(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(insts[i].ContainerID))
        if serr != nil {
            continue
        }
        if status != insts[i].Status {
            insts[i].Status = status
            _ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", insts[i].ID).Update("status", status).Error
        }
    }
    return insts, total, nil
}
func (instService *InstanceService) GetInstanceDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	imageId := make([]map[string]any, 0)
	global.GVA_DB.Table("registries").
		Where("deleted_at IS NULL AND shelf_status = ?", true).
		Select("name as label,id as value").
		Scan(&imageId)
	res["imageId"] = imageId

	serverId := make([]map[string]any, 0)
	global.GVA_DB.Table("compute_nodes").
		Where("deleted_at IS NULL").
		Select("name as label,id as value").
		Scan(&serverId)
	res["serverId"] = serverId

	type specRow struct {
		ID           int64   `gorm:"column:id"`
		Name         *string `gorm:"column:name"`
		GpuModel     *string `gorm:"column:gpu_model"`
		GpuCount     *int64  `gorm:"column:gpu_count"`
		CpuCores     *int64  `gorm:"column:cpu_cores"`
		MemoryGB     *int64  `gorm:"column:memory_gb"`
		SystemDiskGB *int64  `gorm:"column:system_disk_gb"`
		DataDiskGB   *int64  `gorm:"column:data_disk_gb"`
	}
	rows := make([]specRow, 0)
	err = global.GVA_DB.Table("product_specs").
		Where("deleted_at IS NULL AND shelf_status = ?", true).
		Select("id, name, gpu_model, gpu_count, cpu_cores, memory_gb, system_disk_gb, data_disk_gb").
		Scan(&rows).Error
	if err != nil {
		return
	}
	templateId := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		name := ""
		if r.Name != nil {
			name = *r.Name
		}
		gpuModel := ""
		if r.GpuModel != nil {
			gpuModel = *r.GpuModel
		}
		gpuCount := int64(0)
		if r.GpuCount != nil {
			gpuCount = *r.GpuCount
		}
		cpu := int64(0)
		if r.CpuCores != nil {
			cpu = *r.CpuCores
		}
		mem := int64(0)
		if r.MemoryGB != nil {
			mem = *r.MemoryGB
		}
		sysDisk := int64(0)
		if r.SystemDiskGB != nil {
			sysDisk = *r.SystemDiskGB
		}
		dataDisk := int64(0)
		if r.DataDiskGB != nil {
			dataDisk = *r.DataDiskGB
		}

		label := name
		if gpuModel != "" || gpuCount > 0 {
			label += " | " + gpuModel
			if gpuCount > 0 {
				label += " x " + fmt.Sprintf("%d", gpuCount)
			}
		}
		label += " | CPU " + fmt.Sprintf("%d", cpu) + "核"
		label += " | 内存 " + fmt.Sprintf("%d", mem) + "GB"
		label += " | 系统盘 " + fmt.Sprintf("%d", sysDisk) + "GB"
		label += " | 数据盘 " + fmt.Sprintf("%d", dataDisk) + "GB"

		templateId = append(templateId, map[string]any{
			"label": label,
			"value": r.ID,
		})
	}
	res["templateId"] = templateId

	userId := make([]map[string]any, 0)
	global.GVA_DB.Table("sys_users").
		Where("deleted_at IS NULL").
		Select("username as label,id as value").
		Scan(&userId)
	res["userId"] = userId
	return
}

// GetMatchedComputeNodes 根据产品规格匹配满足要求的已上架算力节点（扣除已分配实例占用）
func (instService *InstanceService) GetMatchedComputeNodes(ctx context.Context, specID string) (options []map[string]any, err error) {
	// 获取规格需求
	type specReq struct {
		ID           int64
		GpuModel     *string
		GpuCount     *int64
		CpuCores     *int64
		MemoryGB     *int64
		SystemDiskGB *int64
		DataDiskGB   *int64
	}
	var req specReq
	err = global.GVA_DB.Table("product_specs").
		Where("id = ? AND deleted_at IS NULL AND shelf_status = ?", specID, true).
		Select("id, gpu_model, gpu_count, cpu_cores, memory_gb, system_disk_gb, data_disk_gb").
		Scan(&req).Error
	if err != nil {
		return
	}

	// 取默认需求值（处理nil）
	gpuModel := ""
	if req.GpuModel != nil {
		gpuModel = *req.GpuModel
	}
	needGPU := int64(0)
	if req.GpuCount != nil {
		needGPU = *req.GpuCount
	}
	needCPU := int64(0)
	if req.CpuCores != nil {
		needCPU = *req.CpuCores
	}
	needMem := int64(0)
	if req.MemoryGB != nil {
		needMem = *req.MemoryGB
	}
	needSys := int64(0)
	if req.SystemDiskGB != nil {
		needSys = *req.SystemDiskGB
	}
	needData := int64(0)
	if req.DataDiskGB != nil {
		needData = *req.DataDiskGB
	}

	// 读取已上架节点容量
	nodes := make([]compute.ComputeNode, 0)
	err = global.GVA_DB.Model(&compute.ComputeNode{}).
		Where("deleted_at IS NULL AND shelf_status = ?", true).
		Find(&nodes).Error
	if err != nil {
		return
	}

	// 统计每个节点已占用资源（基于实例绑定的规格）
	type usageRow struct {
		ServerID uint
		UsedCPU  int64
		UsedMem  int64
		UsedSys  int64
		UsedData int64
		UsedGPU  int64
	}
	used := make([]usageRow, 0)
	err = global.GVA_DB.Table("instances AS i").
		Where("i.deleted_at IS NULL").
		Select("i.server_id as server_id, COALESCE(SUM(ps.cpu_cores),0) as used_cpu, COALESCE(SUM(ps.memory_gb),0) as used_mem, COALESCE(SUM(ps.system_disk_gb),0) as used_sys, COALESCE(SUM(ps.data_disk_gb),0) as used_data, COALESCE(SUM(ps.gpu_count),0) as used_gpu").
		Joins("LEFT JOIN product_specs ps ON ps.id = i.template_id").
		Group("i.server_id").
		Scan(&used).Error
	if err != nil {
		return
	}
	usedMap := make(map[uint]usageRow)
	for _, u := range used {
		usedMap[u.ServerID] = u
	}

	// 过滤满足需求的节点（剩余容量）且显卡型号匹配
	options = make([]map[string]any, 0)
	for _, n := range nodes {
		// 取节点容量
		cpu := int64(0)
		mem := int64(0)
		sys := int64(0)
		data := int64(0)
		gpu := int64(0)
		name := ""
		nodeGpuModel := ""
		if n.CPU != nil {
			cpu = *n.CPU
		}
		if n.Memory != nil {
			mem = *n.Memory
		}
		if n.SystemDisk != nil {
			sys = *n.SystemDisk
		}
		if n.DataDisk != nil {
			data = *n.DataDisk
		}
		if n.GpuCount != nil {
			gpu = *n.GpuCount
		}
		if n.Name != nil {
			name = *n.Name
		}
		if n.GpuName != nil {
			nodeGpuModel = *n.GpuName
		}

		u := usedMap[n.ID]
		availCPU := cpu - u.UsedCPU
		availMem := mem - u.UsedMem
		availSys := sys - u.UsedSys
		availData := data - u.UsedData
		availGPU := gpu - u.UsedGPU

		if nodeGpuModel != gpuModel {
			continue
		}
		if availCPU < needCPU || availMem < needMem || availSys < needSys || availData < needData || availGPU < needGPU {
			continue
		}

		options = append(options, map[string]any{
			"label": name,
			"value": n.ID,
		})
	}
	return
}
func (instService *InstanceService) GetInstancePublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// RestartContainer 重启实例对应容器
func (instService *InstanceService) RestartContainer(ctx context.Context, ID string) (err error) {
	var inst instance.Instance
	if err = global.GVA_DB.Where("id = ?", ID).First(&inst).Error; err != nil {
		return err
	}
	var node compute.ComputeNode
	if inst.ServerID != nil {
		_ = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error
	}
    if err = dockerRestartContainer(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(inst.ContainerID)); err != nil {
        return err
    }
    status, _ := dockerInspectStatus(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(inst.ContainerID))
    if status == "" {
        status = "running"
    }
    _ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Update("status", status).Error
    return nil
}

// StopContainer 关闭实例对应容器
func (instService *InstanceService) StopContainer(ctx context.Context, ID string) (err error) {
	var inst instance.Instance
	if err = global.GVA_DB.Where("id = ?", ID).First(&inst).Error; err != nil {
		return err
	}
	var node compute.ComputeNode
	if inst.ServerID != nil {
		_ = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error
	}
    if err = dockerStopContainer(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(inst.ContainerID)); err != nil {
        return err
    }
    status, _ := dockerInspectStatus(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(inst.ContainerID))
    if status == "" {
        status = "exited"
    }
    _ = global.GVA_DB.Model(&instance.Instance{}).Where("id = ?", inst.ID).Update("status", status).Error
    return nil
}

// GetContainerLogs 查看容器日志
func (instService *InstanceService) GetContainerLogs(ctx context.Context, ID string, tail int) (string, error) {
	var inst instance.Instance
	if err := global.GVA_DB.Where("id = ?", ID).First(&inst).Error; err != nil {
		return "", err
	}
	var node compute.ComputeNode
	if inst.ServerID != nil {
		_ = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error
	}
	return dockerContainerLogs(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(inst.ContainerID), tail)
}

// ExecContainerCmd 在容器中执行命令并返回输出
func (instService *InstanceService) ExecContainerCmd(ctx context.Context, ID string, cmd []string) (string, error) {
	var inst instance.Instance
	if err := global.GVA_DB.Where("id = ?", ID).First(&inst).Error; err != nil {
		return "", err
	}
	var node compute.ComputeNode
	if inst.ServerID != nil {
		_ = global.GVA_DB.Where("id = ?", *inst.ServerID).First(&node).Error
	}
	return dockerExecRun(valueOrString(node.DockerEndpoint), valueOrBool(node.UseTLS), valueOrString(node.CACert), valueOrString(node.ClientCert), valueOrString(node.ClientKey), valueOrString(inst.ContainerID), cmd)
}
