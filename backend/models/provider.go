package models

// CloudProvider 表示云提供商
type CloudProvider struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Regions []Region `json:"regions"`
}

// Region 表示区域
type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	AZs  []AZ   `json:"azs"`
}

// AZ 表示可用区
type AZ struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// VPC 表示VPC配置
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
	VpcIndex            int    `json:"vpcIndex"`
	AZ                  string `json:"az,omitempty"`
}

// BucketLifecycleRule 表示存储桶生命周期规则
type BucketLifecycleRule struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	ExpirationDays int    `json:"expirationDays"`
	TransitionDays int    `json:"transitionDays"`
}

// TransitGatewayConfig 表示Transit Gateway配置
type TransitGatewayConfig struct {
	RouteTableName    string `json:"routeTableName"`
	DefaultRouteTable bool   `json:"defaultRouteTable"`
	SubnetIds         string `json:"subnetIds"`
	DnsSupport        bool   `json:"dnsSupport"`
	Ipv6Support       bool   `json:"ipv6Support"`
	AttachmentName    string `json:"attachmentName"`
}

// ComponentConfig 表示云组件配置
type ComponentConfig struct {
	BucketName          string              `json:"bucketName"`
	BucketPolicyType     string              `json:"bucketPolicyType"`
	CustomBucketPolicy   string              `json:"customBucketPolicy"`
	EnableLifecycleRules bool                `json:"enableLifecycleRules"`
	LifecycleRule        BucketLifecycleRule `json:"lifecycleRule"`
	EnableRouteTables    bool                `json:"enableRouteTables"`
	EnableVpcAttachment  bool                `json:"enableVpcAttachment"`
	TransitGatewayConfig TransitGatewayConfig `json:"transitGatewayConfig"`
	TransitGatewayName   string              `json:"transitGatewayName"`
}

// DeploymentConfig 表示部署配置
type DeploymentConfig struct {
	CloudProvider       string                      `json:"cloudProvider"`
	Region              string                      `json:"region"`
	AZ                  string                      `json:"az"`
	VPC                 VPC                         `json:"vpc"`
	Subnet              Subnet                      `json:"subnet"`
	AllVpcs             []VPC                       `json:"allVpcs"`
	AllSubnets          []Subnet                    `json:"allSubnets"`
	Components          []string                    `json:"components"`
	ComponentProperties map[string]interface{}      `json:"componentProperties"`
	ComponentConfig     ComponentConfig             `json:"componentConfig"`
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
