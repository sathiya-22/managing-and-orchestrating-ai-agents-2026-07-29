package agents

// Agent represents an AI agent that can execute a task and return a response.
type Agent interface {
	ID() string
	Execute(prompt string, memory []string) (string, error)
}
