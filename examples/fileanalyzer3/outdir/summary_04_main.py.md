# main.py - File Summary

**one_line**:  
CLI entry point for generating multilingual tutorials from GitHub repos or local directories using a modular flow pipeline.

**Purpose of the file**:  
This script serves as the command-line interface (CLI) for a tutorial generation system that analyzes codebases (either from a GitHub URL or a local directory), extracts key abstractions and relationships, and generates structured, language-specific tutorials. It parses user inputs, configures the processing pipeline, and orchestrates the flow execution.

**Major functions/classes**:  
- `main()`:  
  - Sets up argument parsing using `argparse`.
  - Validates and prepares inputs (source, patterns, token, language, etc.).
  - Initializes the `shared` state dictionary passed through the flow.
  - Instantiates and runs the tutorial generation flow via `create_tutorial_flow()`.
- **Imported dependency**: `create_tutorial_flow()` from `flow` module — defines the core processing steps (not shown here, but assumed to be a DAG of nodes).

**Key technical details & TODOs**:  
- **Source flexibility**: Supports both `--repo` (public GitHub URL) and `--dir` (local path) as mutually exclusive sources.
- **Pattern filtering**: Uses glob-style `include`/`exclude` patterns (defaults defined for common code files and build/test assets).
- **GitHub token handling**: Falls back to `GITHUB_TOKEN` environment variable if not provided; warns on missing token for repo mode.
- **LLM caching control**: `--no-cache` flag disables caching for reproducibility or debugging.
- **Language support**: `--language` parameter allows multilingual output (default: English).
- **Abstraction limit**: `--max-abstractions` caps the number of high-level concepts extracted (default: 10).
- **File size limit**: Skips files larger than `max-size` (default: 100KB) to avoid processing huge files.
- **Extensible shared state**: The `shared` dictionary is the central data bus for the flow, containing inputs, outputs, and metadata.
- **TODOs (implied)**:  
  - No error handling for invalid directories or unreachable repos (relies on downstream flow nodes).
  - No validation of language codes (assumes valid input).
  - No progress tracking or verbose mode (could benefit from `--verbose` flag).

**Short usage example**:  
```bash
# Generate tutorial from GitHub repo in French, with custom patterns
python main.py \
  --repo https://github.com/user/repo \
  --token ghp_... \
  --output ./tutorials/repo \
  --include "*.py" "*.md" \
  --exclude "tests/*" \
  --language french \
  --max-abstractions 15 \
  --no-cache

# Generate tutorial from local directory
python main.py \
  --dir ./my_project \
  --name "My Project" \
  --output ./docs/tutorial \
  --language spanish
```