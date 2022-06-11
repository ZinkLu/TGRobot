package vmshell

type Config struct {
	Username string `configKey:"username"`
	Password string `configKey:"password"`
	ServerId string `configKey:"serverId"`
}
