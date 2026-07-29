This prototype addresses the challenge of managing and orchestrating multiple AI agents, focusing on defining their interactions, memory strategies, and overall control. It provides a simple, yet extensible, framework for developers to specify "who manages the agents" and how they should operate cohesively.

**Problem:** Managing and Orchestrating AI Agents
The agentic AI community faces significant hurdles in effectively coordinating multiple AI agents. Key issues include:
1.  **Interaction Definition:** How do agents communicate and collaborate?
2.  **Memory Management:** How do agents maintain state and context across interactions?
3.  **Overall Control:** Who dictates the flow and behavior of a group of agents?

This project proposes a lightweight, event-driven orchestration system written in Go. Go was chosen for its concurrency primitives (goroutines, channels) which are ideal for managing parallel agent processes and their communications, and its strong typing and performance suitable for infrastructure tools. The goal is to provide a clear, concise, and performant foundation for building more complex agent orchestration logic.

**Setup and Usage:**

1.  **Prerequisites:** Go (version 1.18 or higher)
2.  **Clone the repository:**
    ```bash
    git clone <repository_url>
    cd <repository_name>
    ```
3.  **Run the orchestrator:**
    ```bash
    go run main.go
    ```

The default demonstration will run entirely with fixture data (mock agent responses and predefined interaction patterns). No API keys are required.

**LLM Integration (Optional):**
This prototype is designed to be LLM-agnostic. The `Agent` interface can be implemented by any agent, whether it's backed by an LLM or not. For demonstration purposes, agents in this prototype are "mock" agents that simply return predefined responses.

To integrate a real LLM, you would implement the `Agent` interface with a concrete struct that makes API calls to your chosen LLM provider (e.g., Google Gemini, OpenAI, Anthropic). This would involve:
1.  Creating a new struct (e.g., `GeminiAgent`) that embeds or wraps a client for the LLM provider.
2.  Implementing the `Execute` method to send the `prompt` to the LLM and return its response.
3.  Updating the `main.go` file to instantiate and register your `GeminiAgent` instead of or alongside the `MockAgent`.

An optional `GeminiAgent` is NOT provided in this initial prototype, as the core problem is orchestration, not LLM invocation itself. The existing `MockAgent` and fixture data fully demonstrate the orchestration capabilities.
