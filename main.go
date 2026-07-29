package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Agent represents an AI agent that can execute a task and return a response.
type Agent interface {
	ID() string
	Execute(prompt string, memory []string) (string, error)
}

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

// Orchestrator manages the interaction and lifecycle of multiple agents.
type Orchestrator struct {
	agents map[string]Agent
	memory map[string][]string // Agent-specific memory
	mu     sync.Mutex
	eventBus chan AgentEvent
}

// AgentEvent represents an event in the orchestration system.
type AgentEvent struct {
	SourceAgentID string
	TargetAgentID string
	Payload       string
}

// NewOrchestrator creates a new Orchestrator.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		agents:   make(map[string]Agent),
		memory:   make(map[string][]string),
		eventBus: make(chan AgentEvent, 100), // Buffered channel
	}
}

// RegisterAgent adds an agent to the orchestrator.
func (o *Orchestrator) RegisterAgent(agent Agent) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.agents[agent.ID()] = agent
	o.memory[agent.ID()] = []string{} // Initialize memory for the agent
	log.Printf("Agent '%s' registered.", agent.ID())
}

// StartEventLoop processes events from the event bus.
func (o *Orchestrator) StartEventLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Orchestrator event loop started.")
	for event := range o.eventBus {
		log.Printf("Processing event: Source='%s', Target='%s', Payload='%s'", event.SourceAgentID, event.TargetAgentID, event.Payload)
		o.handleEvent(event)
	}
	log.Println("Orchestrator event loop stopped.")
}

// SendEvent pushes an event onto the event bus.
func (o *Orchestrator) SendEvent(event AgentEvent) {
	o.eventBus <- event
}

// handleEvent processes a single event. This is where orchestration logic lives.
func (o *Orchestrator) handleEvent(event AgentEvent) {
	o.mu.Lock()
	targetAgent, exists := o.agents[event.TargetAgentID]
	if !exists {
		log.Printf("Error: Target agent '%s' not found for event.", event.TargetAgentID)
		o.mu.Unlock()
		return
	}
	o.mu.Unlock() // Release lock before potentially long-running agent execution

	// Update memory for the target agent with the incoming payload
	o.mu.Lock()
	o.memory[event.TargetAgentID] = append(o.memory[event.TargetAgentID], fmt.Sprintf("From %s: %s", event.SourceAgentID, event.Payload))
	currentMemory := o.memory[event.TargetAgentID] // Get a copy for execution
	o.mu.Unlock()

	// Execute the target agent with the payload as prompt and its current memory
	response, err := targetAgent.Execute(event.Payload, currentMemory)
	if err != nil {
		log.Printf("Error executing agent '%s': %v", event.TargetAgentID, err)
		return
	}

	log.Printf("Agent '%s' responded: '%s'", event.TargetAgentID, response)

	// Decision logic: based on the response, decide next steps.
	// This is a simplified example; real systems would have more complex routing/decision.
	if event.TargetAgentID == "task_manager" {
		if response == "Task completed, notify user." {
			o.SendEvent(AgentEvent{
				SourceAgentID: "task_manager",
				TargetAgentID: "user_notifier",
				Payload:       "Task completed successfully.",
			})
		} else if response == "Need more info from data_analyzer." {
			o.SendEvent(AgentEvent{
				SourceAgentID: "task_manager",
				TargetAgentID: "data_analyzer",
				Payload:       "Analyze the provided data and report back.",
			})
		}
	} else if event.TargetAgentID == "data_analyzer" {
		// Data analyzer finished, report back to task manager
		o.SendEvent(AgentEvent{
			SourceAgentID: "data_analyzer",
			TargetAgentID: "task_manager",
			Payload:       "Data analysis complete: [Summary of analysis].",
		})
	}
	// Add the agent's response to its own memory and potentially the source agent's memory
	o.mu.Lock()
	o.memory[event.TargetAgentID] = append(o.memory[event.TargetAgentID], fmt.Sprintf("Self-response: %s", response))
	if event.SourceAgentID != "" { // Don't add to memory if it's an initial event
		o.memory[event.SourceAgentID] = append(o.memory[event.SourceAgentID], fmt.Sprintf("Response from %s: %s", event.TargetAgentID, response))
	}
	o.mu.Unlock()
}

func main() {
	orchestrator := NewOrchestrator()

	// === FIXTURE DATA: Register Mock Agents ===
	taskManagerResponses := map[string]string{
		"Start task processing.":     "Need more info from data_analyzer.",
		"Data analysis complete: [Summary of analysis].": "Task completed, notify user.",
	}
	dataAnalyzerResponses := map[string]string{
		"Analyze the provided data and report back.": "Data analysis complete: [Summary of analysis].",
	}
	userNotifierResponses := map[string]string{
		"Task completed successfully.": "User notified.",
	}

	orchestrator.RegisterAgent(NewMockAgent("task_manager", taskManagerResponses))
	orchestrator.RegisterAgent(NewMockAgent("data_analyzer", dataAnalyzerResponses))
	orchestrator.RegisterAgent(NewMockAgent("user_notifier", userNotifierResponses))

	var wg sync.WaitGroup
	wg.Add(1)
	go orchestrator.StartEventLoop(&wg)

	// === Simulate an initial event ===
	log.Println("Sending initial event to task_manager...")
	orchestrator.SendEvent(AgentEvent{
		SourceAgentID: "system", // 'system' can be a pseudo-agent for initial triggers
		TargetAgentID: "task_manager",
		Payload:       "Start task processing.",
	})

	// Give some time for events to propagate and agents to respond
	time.Sleep(2 * time.Second)

	// Close the event bus and wait for the event loop to finish
	close(orchestrator.eventBus)
	wg.Wait()

	fmt.Println("\n--- Final Agent Memory ---")
	orchestrator.mu.Lock()
	for agentID, mem := range orchestrator.memory {
		fmt.Printf("Agent '%s' Memory:\n", agentID)
		for i, entry := range mem {
			fmt.Printf("  %d: %s\n", i+1, entry)
		}
		fmt.Println("")
	}
	orchestrator.mu.Unlock()

	log.Println("Orchestration complete.")
}
