package common

type ReleaseCmdOpts struct {
	Version        string
	OutputFormat   string
	ConfigFilepath string
	Config         RunConfig
	Component      string
	ListVersions   bool
	ShowReport     bool
	Distribution   string
	OutputFilePath string
}

type RecommendCmdOpts struct {
	CurrentVersion string
	OutputFormat   string
	ConfigFilepath string
	Config         RunConfig
	Component      string
}

type VulnerabilityCmdOpts struct {
	CurrentVersion string
	OutputFormat   string
	ConfigFilepath string
	Config         RunConfig
}

const (
	GITHUB_API_KEY = "GITHUB_API_KEY"
	K8sRepoUrl     = "https://github.com/kubernetes/kubernetes"
	RKEReporURL    = "https://github.com/rancher/rke2"
)

type RunConfig struct {
	GitHub struct {
		APIKey string `yaml:"apiKey"`
	} `yaml:"githubConfig"`

	VulnerabilityScannner ScannerConfig `yaml:"vulnerabilityScannerConfig"`

	Database RunConfigDB `yaml:"db"`

	ConfigSpec struct {
		ConfigYAML ThirdPartyConfig `yaml:"config.yaml"`
	} `yaml:"thirdPartyComponentConfig"`
}

type ThirdPartyConfig struct {
	ThirdPartyComponents      []string              `yaml:"thirdPartyComponents"`
	ThirdPartyComponentPolicy []ThirdPartyComponent `yaml:"thirdPartyComponentPolicy"`
}

type ThirdPartyComponent struct {
	ComponentName string `yaml:"componentName"`
	GitHubSource  string `yaml:"githubSource"`
	TopK          int    `yaml:"topK"`
	Policies      []struct {
		K8sVersion string `yaml:"k8sVersion"`
		MinVersion string `yaml:"minVersion"`
	} `yaml:"policies"`
}

type ScannerConfig struct {
	Snyk struct {
		Endpoint  string `yaml:"endpoint"`
		AuthToken string `yaml:"authToken"`
	} `yaml:"snyk"`
}

type RunConfigDB struct {
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:dbIdx`
	} `yaml:"redis"`
}
