package agents

import (
	"fmt"
	"log"
	"time"
)

// MockAgent is a simple agent that returns a predefined response based on its ID.
type MockAgent struct {
	agentID string
	// For demonstration, we'll use a simple map for "responses"
	// In a real system, this would involve LLM calls or complex logic.
	responses map[string]string
}

// NewMockAgent creates a new MockAgent.
func NewMockAgent(id string, responses map[string]string) *MockAgent {
	return &MockAgent{
		agentID:   id,
		responses: responses,
	}
}

func (ma *MockAgent) ID() string {
	return ma.agentID
}

func (ma *MockAgent) Execute(prompt string, memory []string) (string, error) {
	log.Printf("[%s] Executing with prompt: '%s', memory: %v", ma.agentID, prompt, memory)
	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Simple logic: if prompt matches a key, return its value. Otherwise, default.
	if resp, ok := ma.responses[prompt]; ok {
		return resp, nil
	}
	return fmt.Sprintf("MockAgent %s processed: '%s'", ma.agentID, prompt), nil
}
