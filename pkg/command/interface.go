package command

// Factory is an interface for executing commands.
type Factory interface {
	Create(nameAndArgs []string) Command
}

// Command is an interface of executable and configurable command object.
type Command interface {
	SetDir(dir string) Command
	AddEnv(key, value string) Command
	ConnectIO() Command
	Exec() ([]byte, error)
}
