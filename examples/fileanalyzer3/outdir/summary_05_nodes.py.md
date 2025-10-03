# nodes.py - File Summary

### one_line  
A collection of PocketFlow `Node` classes that fetch codebases, identify core abstractions, analyze relationships, and order educational chapters for onboarding.

---

### Purpose  
This file implements a pipeline for **automatically analyzing codebases** (from GitHub or local directories) to generate beginner-friendly educational content. It identifies key abstractions, explains how they relate, and orders them into a coherent learning sequence — ideal for documentation, onboarding, or codebase exploration tools.

---

### Major Functions/Classes

| Class | Role |
|------|------|
| `FetchRepo` | Fetches files from a GitHub repo or local directory using configurable filters (patterns, size). |
| `IdentifyAbstractions` | Uses an LLM to extract top core abstractions (with names, descriptions, and file references) from code. |
| `AnalyzeRelationships` | Analyzes how abstractions interact and generates a project summary + relationship graph. |
| `OrderChapters` | Determines a logical learning order (chapter sequence) for abstractions based on their relationships. |
| `get_content_for_indices()` *(helper)* | Extracts file content by index for context injection into LLM prompts. |

---

### Key Technical Details & TODOs

- **LLM Integration**: Uses `call_llm()` to process natural language prompts and parse structured YAML output.
- **YAML Parsing & Validation**: All LLM outputs are parsed as YAML and rigorously validated for required fields, types, and index bounds.
- **Language Support**: Non-English abstractions/summaries/labels are supported via `language` parameter (e.g., Spanish, Japanese).
- **Caching**: LLM calls use `use_cache` flag (default: `True`) and only cache first attempt to avoid stale retries.
- **File Handling**: 
  - Converts fetched files to list of `(path, content)` tuples.
  - Uses **index-based referencing** in LLM prompts to avoid path duplication.
- **Error Handling**: Raises `ValueError` on missing files, invalid YAML, or malformed LLM output.
- **Extensibility**: Designed to plug into a larger PocketFlow pipeline; uses `shared` context for data passing.

#### ✅ TODOs / Improvements (inferred):
- [ ] Add **token budgeting** for large codebases in LLM context creation.
- [ ] Support **chunking or summarization** of very large files before LLM ingestion.
- [ ] Allow **custom prompt templates** per language or project type.
- [ ] Add **retry logic with exponential backoff** for LLM calls.
- [ ] Consider **abstraction deduplication** or merging for overlapping concepts.
- [ ] Validate that **all abstractions are included** in chapter order (currently not enforced).

---

### Short Usage Example

```python
from pocketflow import Flow
from nodes import FetchRepo, IdentifyAbstractions, AnalyzeRelationships, OrderChapters

# Shared context
shared = {
    "repo_url": "https://github.com/user/example-project",
    "include_patterns": ["*.py", "*.js"],
    "exclude_patterns": ["test*", "*.md"],
    "max_file_size": 50000,
    "language": "english",
    "use_cache": True,
    "max_abstraction_num": 8
}

# Build flow
flow = Flow(
    FetchRepo() >>
    IdentifyAbstractions() >>
    AnalyzeRelationships() >>
    OrderChapters()
)

# Run
flow.run(shared)

# Access results
print("Project Summary:", shared["relationships"]["summary"])
print("Abstractions:", shared["abstractions"])
print("Relationships:", shared["relationships"]["details"])
print("Chapter Order:", shared["chapters"])  # List of abstraction indices in order
```

> 📝 This pipeline is ideal for generating **interactive documentation**, **onboarding guides**, or **codebase maps** automatically.