package server

const (
	// RecommendedHomeDir defines the default directory used to place all  service configurations.
	RecommendedHomeDir = ".voyage"

	// RecommendedEnvPrefix defines the ENV prefix used by all voyage service.
	RecommendedEnvPrefix = "VOYAGE"
)

// SecureServingInfo holds configuration of the TLS server.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}
