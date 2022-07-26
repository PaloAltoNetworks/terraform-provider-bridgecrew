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
	Provider       string         `json:"provider"`
	ID             int            `json:"id,omitempty"`
	Title          string         `json:"title"`
	Severity       string         `json:"severity"`
	Category       string         `json:"category"`
	Guidelines     string         `json:"guidelines"`
	ConditionQuery ConditionQuery `json:"conditions,omitempty"`
	Benchmarks     Benchmark      `json:"benchmarks,omitempty"`
	Frameworks     []string       `json:"frameworks"`
}

type simplePolicy struct {
	Provider   string     `json:"provider"`
	ID         int        `json:"id,omitempty"`
	Title      string     `json:"title"`
	Severity   string     `json:"severity"`
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

// Amount is a sub-object of Account
type Amount struct {
	CLOSED     int
	DELETED    int
	OPEN       int
	REMEDIATED int
	SUPPRESSED int
}

// ConditionQuery is the construct for the complex query screen
type ConditionQuery struct {
	Ands []Conditions `json:"and,omitempty"`
}

// Conditions is part of the simple query
type Conditions struct {
	Attribute     string   `json:"attribute,omitempty"`
	CondType      string   `json:"cond_type,omitempty"`
	Operator      string   `json:"operator,omitempty"`
	ResourceTypes []string `json:"resource_types,omitempty"`
	Value         string   `json:"value,omitempty"`
	Or            []Or     `json:"or,omitempty"`
}

// Or Is the Condition query to construct an Or block in tf
type Or struct {
	Attribute     string   `json:"attribute,omitempty"`
	CondType      string   `json:"cond_type,omitempty"`
	Operator      string   `json:"operator,omitempty"`
	ResourceTypes []string `json:"resource_types,omitempty"`
	Value         string   `json:"value,omitempty"`
}

// Result is for parsing return messages from the platform
type Result struct {
	Policy string
	ID     string
}

// User is for eventually managing users
type User struct {
	Role         string   `json:"role"`
	Email        string   `json:"email"`
	Accounts     []string `json:"accounts"`
	LastModified int64    `json:"last_modified"`
	CustomerName string   `json:"customer_name"`
}

// Branch contains a branches used in CICD
type Branch struct {
	Name          string `json:"name"`
	CreationDate  string `json:"creationdate"`
	DefaultBranch bool   `json:"defaultbranch"`
}

// Repositories For CICD
type Repositories struct {
	Source   string   `json:"source"`
	Branches []Branch `json:"branches"`
}

// Tag is a structure for writing tag rules for Yor
type Tag struct {
	Name          string     `json:"name"`
	Description   string     `json:"description,omitempty"`
	TagRuleOOTBId string     `json:"source,omitempty"`
	IsEnabled     bool       `json:"isenabled,omitempty"`
	Repositories  []string   `json:"repositories"`
	Definition    Definition `json:"ruleDefinition"`
}

// Definition for Yor tags
type Definition struct {
	TagGroups []TagGroup `json:"tag_groups"`
}

// TagGroup for yor definitions
type TagGroup struct {
	Name string `json:"name"`
	Tags []Tags `json:"tags"`
}

// Tags struct to part of tag group
type Tags struct {
	Name  string                 `json:"name"`
	Value map[string]interface{} `json:"value"`
}

// Value part of Tags struct
type Value struct {
	Default string `json:"default"`
}

// Rule datatype for enforcement rule
type Rule struct {
	Name           string         `json:"name"`
	CodeCategories CodeCategories `json:"codeCategories"`
	Repositories   []Repo         `json:"repositories,omitempty"`
}

// CodeCategories is a data type for enforcement rules
type CodeCategories struct {
	OpenSource  Category `json:"OPEN_SOURCE"`
	Images      Category `json:"IMAGES"`
	IAC         Category `json:"IAC"`
	Secrets     Category `json:"SECRETS"`
	SupplyChain Category `json:"SUPPLY_CHAIN"`
}

// Category set of enforcement rules
type Category struct {
	HardFailThreshold    string `json:"hardFailThreshold"`
	SoftFailThreshold    string `json:"softFailThreshold"`
	CommentsBotThreshold string `json:"commentsBotThreshold"`
}

// Repo strucuture for repo/accounts
type Repo struct {
	AccountID   string `json:"accountId"`
	AccountName string `json:"accountName"`
}
