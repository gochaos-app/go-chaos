package cmd

type GenConfig struct {
	App         string      `hcl:"app"`
	Description string      `hcl:"description"`
	Job         []JobConfig `hcl:"job,block"`
}

type JobConfig struct {
	Region    string      `hcl:"region,optional"`
	Namespace string      `hcl:"namespace,optional"`
	Cloud     string      `hcl:"cloud,label"`
	Service   string      `hcl:"service,label"`
	Chaos     ChaosConfig `hcl:"config,block"`
}

type ChaosConfig struct {
	Config string   `hcl:"chaos,label"`
	Tags   []string `hcl:"tags"`
	Chaos  string   `hcl:"chaos"`
	Count  int      `hcl:"count"`
}
