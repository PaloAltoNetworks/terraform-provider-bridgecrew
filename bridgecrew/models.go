package bridgecrew

// The Policy record
type Policy struct {
	Provider   string    `json:"provider"`
	ID         int       `json:"id,omitempty"`
	Benchmarks Benchmark `json:"benchmarks,omitempty"`
	Code       string    `json:"code,omitempty"`
	Frameworks []string  `json:"frameworks"`
}

type complexPolicy struct {
	Provider string `json:"provider"`
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	//PRISMASeverity string         `json:"pcseverity"`
	Category       string         `json:"category"`
	Guidelines     string         `json:"guidelines"`
	ConditionQuery ConditionQuery `json:"conditions,omitempty"`
	Benchmarks     Benchmark      `json:"benchmarks,omitempty"`
	Frameworks     []string       `json:"frameworks"`
}

type simplePolicy struct {
	Provider string `json:"provider"`
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	//PRISMASeverity string     `json:"pcseverity"`
	Category   string     `json:"category"`
	Guidelines string     `json:"guidelines"`
	Conditions Conditions `json:"conditions,omitempty"`
	Benchmarks Benchmark  `json:"benchmarks,omitempty"`
	Frameworks []string   `json:"frameworks"`
}

// Benchmark is child object to Policy
type Benchmark struct {
	Cisawsv12        []string `json:"CIS AWS V1.2,omitempty"`
	Cisawsv13        []string `json:"CIS AWS V1.3,omitempty"`
	Cisazurev11      []string `json:"CIS AZURE V11,omitempty"`
	Cisazurev12      []string `json:"CIS AZURE V12,omitempty"`
	Cisazurev13      []string `json:"CIS AZURE V13,omitempty"`
	Ciskubernetesv15 []string `json:"CIS KUBERNETES V15,omitempty"`
	Ciskubernetesv16 []string `json:"CIS KUBERNETES V16,omitempty"`
	Cisgcpv11        []string `json:"CIS GCP V11,omitempty"`
	Cisgkev11        []string `json:"CIS GKE V11,omitempty"`
	Cisdockerv11     []string `json:"CIS DOCKER V11,omitempty"`
	Ciseksv11        []string `json:"CIS EKS V11,omitempty"`
}

//// Account is a child object to Policy
//type Account struct {
//	repository     string
//	amounts        Amount
//	lastupdatedate string
//}

// Amount is a sub-object of Account
type Amount struct {
	CLOSED     int
	DELETED    int
	OPEN       int
	REMEDIATED int
	SUPPRESSED int
}

//ConditionQuery is the construct for the complex query screen
type ConditionQuery struct {
	Ands []Conditions `json:"and,omitempty"`
}

//Conditions is part of the simple query
type Conditions struct {
	Attribute     string   `json:"attribute,omitempty"`
	CondType      string   `json:"cond_type,omitempty"`
	Operator      string   `json:"operator,omitempty"`
	ResourceTypes []string `json:"resource_types,omitempty"`
	Value         string   `json:"value,omitempty"`
	Or            []Or     `json:"or,omitempty"`
}

//Or Is the Condition query to construct an Or block in tf
type Or struct {
	Attribute     string   `json:"attribute,omitempty"`
	CondType      string   `json:"cond_type,omitempty"`
	Operator      string   `json:"operator,omitempty"`
	ResourceTypes []string `json:"resource_types,omitempty"`
	Value         string   `json:"value,omitempty"`
}

//Result is for parsing return messages from the platform
type Result struct {
	Policy string
}

//User is for eventually managing users
type User struct {
	Role         string   `json:"role"`
	Email        string   `json:"email"`
	Accounts     []string `json:"accounts"`
	LastModified int64    `json:"last_modified"`
	CustomerName string   `json:"customer_name"`
}
