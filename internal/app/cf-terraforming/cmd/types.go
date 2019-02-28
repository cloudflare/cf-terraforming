package cmd

type StateResponse struct {
	TFStateScaffold
}

type TFStateScaffold struct {
	Version          int      `json:"version"`
	TerraformVersion string   `json:"terraform_version"`
	Serial           int      `json:"serial"`
	Lineage          string   `json:"lineage"`
	Modules          []Module `json:"modules"`
}

type Module struct {
	Path      []string               `json:"path"`
	DependsOn []string               `json:"depends_on"`
	Outputs   interface{}            `json:"outputs"`
	Resource  map[string]interface{} `json:"resources"`
}

type Resource struct {
	Primary   Primary  `json:"primary"`
	DependsOn []string `json:"depends_on"`
	Deposed   []string `json:"deposed"`
	Provider  string   `json:"provider"`
	Type      string   `json:"type"`
}

type Primary struct {
	Id         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
	Meta       interface{} `json:"meta"`
	Tainted    bool        `json:"tainted"`
}

type Attributes struct {
	Enabled     string `json:"enabled"`
	Id          string `json:"id"`
	MultiScript bool   `json:"multi_script"`
	Pattern     string `json:"pattern"`
	Zone        string `json:"zone"`
	ZoneId      string `json:"zone_id"`
}
