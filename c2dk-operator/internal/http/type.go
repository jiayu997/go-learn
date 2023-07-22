package http

type Resource struct {
	Memory string `json:"memory"`
	Cpu    string `json:"cpu"`
}

// 就绪探针、存活探针定义 类型为list
type ProbeInfo struct {
	//初始延迟
	InitialDelaySecond int32 `json:"initial_delay_second"`
	//允许失败次数
	FailureThreshold int32 `json:"failure_threshold"`
	//协议
	Scheme string `json:"scheme"`
	//探测成功次数
	SuccessThreshold int32 `json:"success_threshold"`
	// 超时
	TimeoutSecond int32 `json:"timeout_second"`
	// 间隔时间
	PeriodSecond int32 `json:"period_second"`
	// 端口
	Port int32 `json:"port"`
	// 模式 readiness liveness
	Mode string `json:"mode"`
	// path
	Path string `json:"path"`
}

// 制品信息
type ArtifactInfo struct {
	// 制品名称
	ArtifactName string `json:"artifact_name"`
	// 制品编码
	ArtifactCode string `json:"artifact_code"`
	// 制品描述
	Describe string `json:"describe"`
	// 镜像坐标
	Image string `json:"image"`
	// 镜像版本
	ImageVersion string `json:"image_version"`
}

// 依赖信息
type DepServiceInfo struct {
	// 被依赖的服务code
	MicroServiceCode string `json:"microservice_code"`
	// 名称
	AttrName string `json:"attr_name"`
	//类型（分为系统注入、自定义）
	AttrType string `json:"attr_type"`
	// 值
	AttrValue string `json:"attr_value"`
	// 描述
	Describe string `json:"describe"`
}

type ConnectInfo struct {
	// 名称
	AttrName string `json:"attr_name"`
	//类型（分为系统注入、自定义）
	AttrType string `json:"attr_type"`
	// 值
	AttrValue string `json:"attr_value"`
	// 描述
	Describe string `json:"describe"`
}
type VolumeInfo struct {
	// 存储卷名称
	VolumeName string `json:"volume_name"`
	// 挂载的绝对路径
	VolumePath string `json:"volume_path"`
	// 卷读写权限
	VolumeAccessMode string `json:"volume_access_mode"`
}

type ConfigMapContext struct {
	// volume-name/mountpath name generate by random
	Name string `json:"name,omitempty"`

	// key
	Key string `json:"key"`
	//挂载内容（适用于configmap类型）
	FileContent string `json:"file_content"`
	//挂载的绝对路径
	VolumePath string `json:"volume_path"`
	// 文件读写权限
	FileAccessMode string `json:"file_access_mode"`
	// subpath
	SubPath string `json:"subpath"`
}

type ConfigMapInfo struct {
	// configmap name
	VolumeName string `json:"volume_name"`

	// mount info for configmap
	ConfigMapContextList []ConfigMapContext `json:"configmap_context"`
}

// todo: add secret type
type StorageInfo struct {
	VolumeList    []VolumeInfo    `json:"volumes"`
	ConfigMapList []ConfigMapInfo `json:"configmaps"`
}

type PortInfo struct {
	// 端口名称
	PortName string `json:"port_name"`
	// 端口协议
	Protocol string `json:"protocol"`
	// 容器端口
	ContainerPort int32 `json:"container_port"`
	// 节点端口
	NodePort int32 `json:"node_port"`
	// 是否开放集群外端口
	IOuterService bool `json:"is_outer_service"`
	// 是否为集群内端口
	IsInnerService bool `json:"is_inner_service"`
}

type EnvInfo struct {
	AttrName  string `json:"attr_name"`
	AttrValue string `json:"attr_value"`
	AttrType  string `json:"attr_type"`
	Status    string `json:"status"`
}

type MicroServiceInfo struct {
	MicroServiceName      string            `json:"microservice_name"`
	MicroServiceCode      string            `json:"microservice_code"`
	MicroServiceNamespace string            `json:"micro_service_namespace"`
	MicroServiceTags      map[string]string `json:"microservice_tags"`
	Kind                  string            `json:"kind"`
	InstanceCount         int32             `json:"instance_count"`
	RequireResource       Resource          `json:"require_resource"`
	LimitResource         Resource          `json:"limit_resource"`
	Labels                map[string]string `json:"labels"`
	Annotations           map[string]string `json:"annotations"`
	ProbeList             []ProbeInfo       `json:"probes"`
	Artifact              ArtifactInfo      `json:"artifact_info"`
	DepServiceMapList     []DepServiceInfo  `json:"dep_service_map_list"`
	ConnectInfoMapList    []ConnectInfo     `json:"connect_info_map_list"`
	Storage               StorageInfo       `json:"storage"`
	PortMapList           []PortInfo        `json:"port_map_list"`
	EnvMapList            []EnvInfo         `json:"env_map_list"`
	Cmd                   []string          `json:"cmd"`
	Args                  []string          `json:"args"`
	Language              string            `json:"language"`
}

type MiddlewareInfo struct {
}

type Template struct {
	TemplateID       string             `json:"template_id"`
	TemplateVersion  string             `json:"template_version"`
	AppName          string             `json:"app_name"`
	AppCode          string             `json:"app_code"`
	AppDescribe      string             `json:"app_describe"`
	MicroServiceList []MicroServiceInfo `json:"microservice"`
	Middleware       MiddlewareInfo     `json:"middleware"`
}
