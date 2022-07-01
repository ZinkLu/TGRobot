package tgo

type Config struct {
	Addr string `configKey:"api_addr"`
	Port int    `configKey:"api_port"`
}
