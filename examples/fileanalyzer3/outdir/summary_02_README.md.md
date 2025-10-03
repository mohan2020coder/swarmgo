# README.md - File Summary

## one_line
A tutorial project that uses AI to automatically analyze GitHub repositories and generate beginner-friendly, well-visualized codebase tutorials.

## Purpose of the file
The `README.md` serves as the main documentation and landing page for a project that transforms any codebase into an accessible, AI-generated tutorial using Pocket Flow—a lightweight LLM agent framework. It explains how to set up, run, and customize the tool, showcases example outputs from popular open-source projects, and provides development context and usage instructions.

## Major functions/classes
- **Main entry point**: `main.py` – Orchestrates the entire flow: repo/dir crawling, code analysis, abstraction detection, and tutorial generation.
- **Codebase crawler**: Downloads GitHub repos or reads local directories, filters files based on include/exclude patterns and size limits.
- **Abstraction analyzer**: Uses LLM agents to detect core concepts, patterns, and interactions in the codebase.
- **Tutorial generator**: Transforms analysis into structured, beginner-friendly tutorials with visualizations (e.g., diagrams, flowcharts).
- **LLM interface**: `utils/call_llm.py` – Handles communication with Gemini, Claude, or other LLMs for reasoning and content generation.
- **Pocket Flow integration**: Lightweight framework enabling modular, agent-driven workflow execution (see [PocketFlow](https://github.com/The-Pocket/PocketFlow)).

## Key technical details & TODOs
- **LLM-powered analysis**: Uses modern LLMs (Gemini 2.5 Pro, Claude 3.7 with thinking, O1) to reason about code structure and abstractions.
- **Caching**: Enabled by default to avoid redundant LLM calls; can be disabled via `--no-cache`.
- **Language support**: Tutorials can be generated in any language (`--language` flag).
- **Configurable scope**: Filter files by extension (`--include`), exclude paths (`--exclude`), and limit file size (`--max-size`).
- **GitHub integration**: Supports public/private repos via `--repo` and `--token` or `GITHUB_TOKEN`.
- **Docker support**: Fully containerized with volume mounts for output and code input.
- **Extensible design**: Built using Pocket Flow's agentic coding paradigm—agents build agents—for modular, scalable workflows.
- **TODOs / Improvements (implied)**:
  - Add support for more visualization types (e.g., Mermaid diagrams in output).
  - Improve handling of large or complex codebases (e.g., better abstraction summarization).
  - Enhance caching strategy (e.g., persistent cache across runs).
  - Support multi-repo or monorepo analysis.

## Short usage example
```bash
# Analyze a GitHub repo and generate an English tutorial
python main.py --repo https://github.com/encode/fastapi --include "*.py" --exclude "tests/*" --max-size 100000

# Analyze a local directory and generate a Chinese tutorial
python main.py --dir ./my-project --include "*.js" "*.ts" --language "Chinese"

# Run in Docker (requires GEMINI_API_KEY)
docker run -it --rm \
  -e GEMINI_API_KEY="your_key" \
  -v "$(pwd)/output":/app/output \
  pocketflow-app --repo https://github.com/pallets/flask
```

> ✅ Output: A fully AI-generated tutorial (HTML or Markdown) saved in `./output`, explaining the codebase's architecture, key components, and how everything fits together—perfect for onboarding or learning.