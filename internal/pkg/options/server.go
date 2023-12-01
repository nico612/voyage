package options

const EnvModel = "ENV_MODEL"

const (
	DebugModel   = "debug"
	ReleaseModel = "release"
	TestModel    = "test"
)

type ServerRunOptions struct {
	Mode          string `json:"mode"        mapstructure:"mode"`
	Healthz       bool   `json:"healthz"     mapstructure:"healthz"`
	UseMultipoint bool   `json:"use-multipoint" mapstructure:"use-multipoint"` // 多点登录
	IplimitCount  int32  `json:"iplimit-count"        mapstructure:"iplimit-count"`
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Mode:          ReleaseModel,
		Healthz:       true,
		UseMultipoint: false,
		IplimitCount:  15000,
	}
}
