package cmd

type GenConfig struct {
	App           string        `hcl:"app"`
	Description   string        `hcl:"description"`
	Function      string        `hcl:"function,optional"`
	Job           []JobConfig   `hcl:"job,block"`
	Hypothesis    *Hypothesis   `hcl:"hypothesis,block"`
	Notifications []NotifConfig `hcl:"notification,block"`
}

type Hypothesis struct {
	Name        string `hcl:"name"`
	Description string `hcl:"description"`
	Pings       string `hcl:"workers"`
	Url         string `hcl:"url"`
	Report      string `hcl:"report"`
}

type NotifConfig struct {
	Type string   `hcl:"type,label"`
	From string   `hcl:"from,optional"` // Only used for gmail notification
	To   []string `hcl:"to"`            // Several emails or channels can be specified here
	Body string   `hcl:"body"`
}

type JobConfig struct {
	Region    string      `hcl:"region,optional"`
	Namespace string      `hcl:"namespace,optional"`
	Project   string      `hcl:"project,optional"`
	Cloud     string      `hcl:"cloud,label"`
	Service   string      `hcl:"service,label"`
	Chaos     ChaosConfig `hcl:"config,block"`
}

type ChaosConfig struct {
	Tag   string `hcl:"tag"`
	Chaos string `hcl:"chaos"`
	Count int    `hcl:"count"`
}

/*
type ScriptConfig struct {
	Executor string `hcl:"executor"`
	Source   string `hcl:"source"`
}
*/
