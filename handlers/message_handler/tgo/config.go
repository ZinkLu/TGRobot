package tgo

type Config struct {
	Addr     string   `configKey:"api_addr"`
	Port     int      `configKey:"api_port"`
	UserCert bool     `configKey:"use_cert"`
	CaPaths  []string `configKey:"ca_paths"`
	CertPath string   `configKey:"cert_path"`
	CertKey  string   `configKey:"cert_key"`
	Verify   bool     `configKey:"verify"`
	SNI      string   `configKey:"sni"`
}
