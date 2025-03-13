package models

// Provider 表示云服务提供商
type Provider struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Logo  string `json:"logo"`
}

// Region 表示云服务提供商的区域
type Region struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AvailabilityZone 表示区域内的可用区
type AvailabilityZone struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ComponentProperty 表示云组件的属性
type ComponentProperty struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
	Placeholder  string `json:"placeholder"`
	Description  string `json:"description"`
}

// Component 表示云组件
type Component struct {
	Name        string             `json:"name"`
	Value       string             `json:"value"`
	Description string             `json:"description"`
	Properties  []ComponentProperty `json:"properties"`
}

// VPC 表示虚拟私有云配置
type VPC struct {
	Name              string `json:"name"`
	CIDR              string `json:"cidr"`
	EnableDnsSupport  bool   `json:"enableDnsSupport,omitempty"`
	EnableDnsHostnames bool  `json:"enableDnsHostnames,omitempty"`
}

// Subnet 表示子网配置
type Subnet struct {
	Name                string `json:"name"`
	CIDR                string `json:"cidr"`
	MapPublicIpOnLaunch bool   `json:"mapPublicIpOnLaunch,omitempty"`
}

// DeploymentConfig 表示部署配置
type DeploymentConfig struct {
	CloudProvider      string                     `json:"cloudProvider"`
	Region             string                     `json:"region"`
	AZ                 string                     `json:"az"`
	VPC                VPC                        `json:"vpc"`
	Subnet             Subnet                     `json:"subnet"`
	Components         []Component                `json:"components"`
	ComponentProperties map[string]map[string]string `json:"componentProperties"`
}

// DeploymentStatus 表示部署状态
type DeploymentStatus struct {
	Status   string      `json:"status"` // idle, preparing, deploying, completed, failed
	Progress int         `json:"progress"`
	Message  string      `json:"message"`
	Logs     []string    `json:"logs"`
	Result   interface{} `json:"result"`
	Topology interface{} `json:"topology"`
}

// TopologyNode 表示拓扑图中的节点
type TopologyNode struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

// TopologyEdge 表示拓扑图中的连接
type TopologyEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

// Topology 表示资源拓扑图
type Topology struct {
	Nodes []TopologyNode `json:"nodes"`
	Edges []TopologyEdge `json:"edges"`
}
