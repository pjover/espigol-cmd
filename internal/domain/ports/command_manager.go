package ports

type CommandManager interface {
	AddCommand(interface{})
	Execute() []string
}
