package vmshell

type Config struct {
	Username  string   `configKey:"username"`
	Password  string   `configKey:"password"`
	ServerIds []string `configKey:"serverIds"`
}
