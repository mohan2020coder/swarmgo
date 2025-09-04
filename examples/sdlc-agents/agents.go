package main

import swarmgo "github.com/mohan2020coder/swarmgo"

// Define reusable agents as variables instead of functions
var (
	Planner = &swarmgo.Agent{
		Name:         "Planner",
		Instructions: "You are the Project Planner. Break the user's SDLC task into clear, short steps.",
		Model:        ModelName,
	}

	Architect = &swarmgo.Agent{
		Name:         "Architect",
		Instructions: "You are the Software Architect. Suggest design, structure, and edge cases.",
		Model:        ModelName,
	}

	Coder = &swarmgo.Agent{
		Name: "Coder",
		Instructions: `You are the Implementer. Provide the minimal working code.
Return code in ONE fenced block.`,
		Model: ModelName,
	}

	Reviewer = &swarmgo.Agent{
		Name:         "Reviewer",
		Instructions: `You are the Reviewer. Approve or suggest minimal fixes. Always output a corrected code block if needed.`,
		Model:        ModelName,
	}
)
