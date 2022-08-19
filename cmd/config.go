package cmd

type GenConfig struct {
	App         string        `hcl:"app"`
	Description string        `hcl:"description"`
	Job         []JobConfig   `hcl:"job,block"`
	Script      *ScriptConfig `hcl:"script,block"`
}

type JobConfig struct {
	Region    string      `hcl:"region,optional"`
	Namespace string      `hcl:"namespace,optional"`
	Cloud     string      `hcl:"cloud,label"`
	Service   string      `hcl:"service,label"`
	Chaos     ChaosConfig `hcl:"config,block"`
}

type ChaosConfig struct {
	Tag   string `hcl:"tag"`
	Chaos string `hcl:"chaos"`
	Count int    `hcl:"count"`
}

type ScriptConfig struct {
	Executor string `hcl:"executor"`
	Source   string `hcl:"source"`
}
