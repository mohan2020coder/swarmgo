```markdown
# Core Architectural Abstractions & Components

Below are up to 15 core abstractions/components identified from the codebase, based on file summaries and their relationships.

---

## 1. **TutorialFlow**
**Description**:  
Orchestrates the end-to-end pipeline for generating AI-powered codebase tutorials using PocketFlow.

**Files**:  
- `flow.py`, `main.py`, `nodes.py`

**Dependencies**:  
- `FetchRepo`, `IdentifyAbstractions`, `AnalyzeRelationships`, `OrderChapters`, `WriteChapters`, `CombineTutorial`, `call_llm`, `crawl_github_files`, `crawl_local_files`

**Responsibilities**:  
- Defines and executes the sequential workflow (DAG) for tutorial generation.
- Manages data flow between nodes via the `shared` context.
- Enables retry logic and batch processing for robustness and performance.

---

## 2. **FetchRepo**
**Description**:  
Retrieves source code from a GitHub repository or local directory with filtering and size constraints.

**Files**:  
- `nodes.py` (class: `FetchRepo`), `utils/crawl_github_files.py`, `utils/crawl_local_files.py`

**Dependencies**:  
- `crawl_github_files`, `crawl_local_files`, `pathspec`, `requests`, `tempfile`

**Responsibilities**:  
- Parses repo URL or local path input.
- Downloads and filters files using glob patterns and `.gitignore`.
- Enforces maximum file size to prevent memory overflow.
- Returns structured list of `(path, content)` tuples for downstream processing.

---

## 3. **IdentifyAbstractions**
**Description**:  
Uses LLMs to extract high-level abstractions (e.g., classes, patterns) from raw code.

**Files**:  
- `nodes.py` (class: `IdentifyAbstractions`)

**Dependencies**:  
- `call_llm`, `utils.read_json`, `yaml`, `shared` context

**Responsibilities**:  
- Constructs prompts with code context (by file index) for LLM.
- Parses LLM-generated YAML to extract abstraction names, descriptions, and file references.
- Validates output structure and enforces limits (e.g., max 10 abstractions).
- Stores results in `shared["abstractions"]`.

---

## 4. **AnalyzeRelationships**
**Description**:  
Analyzes how identified abstractions interact and generates a project summary and relationship map.

**Files**:  
- `nodes.py` (class: `AnalyzeRelationships`)

**Dependencies**:  
- `call_llm`, `get_content_for_indices`, `shared` context

**Responsibilities**:  
- Builds LLM prompt with abstractions and relevant file snippets.
- Extracts project summary and pairwise relationships (e.g., "X depends on Y").
- Outputs structured YAML into `shared["relationships"]` for tutorial structuring.

---

## 5. **OrderChapters**
**Description**:  
Determines a pedagogical order for tutorial chapters based on abstraction relationships.

**Files**:  
- `nodes.py` (class: `OrderChapters`)

**Dependencies**:  
- `call_llm`, `shared` context (abstractions, relationships)

**Responsibilities**:  
- Uses LLM to reason about dependencies and learning sequence.
- Returns a list of abstraction indices in logical order (e.g., from foundational to advanced).
- Stores result in `shared["chapters"]`.

---

## 6. **WriteChapters**
**Description**:  
Generates written content for each chapter (abstraction) using LLM, in a target language.

**Files**:  
- `nodes.py` (class: `WriteChapters`, implemented as `BatchNode`)

**Dependencies**:  
- `call_llm`, `get_content_for_indices`, `shared` context

**Responsibilities**:  
- Processes multiple abstractions in parallel (batch mode).
- For each chapter, generates beginner-friendly explanations with code examples.
- Supports multilingual output via language parameter.
- Stores chapter content in `shared["chapter_contents"]`.

---

## 7. **CombineTutorial**
**Description**:  
Merges all generated chapters into a single, structured tutorial document.

**Files**:  
- `nodes.py` (class: `CombineTutorial`)

**Dependencies**:  
- `shared` context (`chapter_contents`, `project_summary`, `abstractions`)

**Responsibilities**:  
- Combines chapter content in the ordered sequence.
- Adds metadata (title, language, summary).
- Outputs final tutorial (Markdown/HTML) to specified output directory.

---

## 8. **CallLLM**
**Description**:  
Unified interface to call various LLMs (Gemini, Claude