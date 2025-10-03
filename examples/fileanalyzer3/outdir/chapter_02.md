# **Chapter 2: System Overview & High-Level Architecture**

## **Objective**
In this chapter, you’ll gain a comprehensive understanding of how the **AI-Powered Codebase Tutorial Generator** works end-to-end. We’ll walk through the high-level architecture, explain the core components, and show how they interact to transform any codebase — whether from GitHub or your local machine — into a beginner-friendly, multilingual tutorial with visualizable structure.

By the end of this chapter, you’ll understand:
- The **overall data flow** from code input to final tutorial output.
- The **modular design** powered by *PocketFlow* agents.
- The **key architectural abstractions** and their responsibilities.
- How **LLMs are used** to reason about code, not just generate text.
- How to **extend or customize** the system for your own use cases.

---

## **2.1 End-to-End Flow: From Code to Tutorial**

Let’s start with the big picture. Here’s what happens when you run:

```bash
python main.py --repo https://github.com/encode/fastapi --language english
```

### 🔁 **Step-by-Step Pipeline**

| Step | Component | What It Does |
|------|---------|-------------|
| 1️⃣ | `FetchRepo` | Downloads or reads the codebase using pattern filters and size limits. |
| 2️⃣ | `IdentifyAbstractions` | Uses an LLM to detect core concepts (e.g., `FastAPI`, `APIRouter`, `Dependency`) in the code. |
| 3️⃣ | `AnalyzeRelationships` | Asks the LLM how these abstractions relate (e.g., "`APIRouter` is used by `FastAPI`"). |
| 4️⃣ | `OrderChapters` | Determines a logical learning order (e.g., start with `FastAPI`, then `APIRouter`). |
| 5️⃣ | `WriteChapters` *(Batch)* | Generates detailed, beginner-friendly explanations for each abstraction — in parallel. |
| 6️⃣ | `CombineTutorial` | Merges all chapters into a single structured document (Markdown/HTML), including summary and metadata. |

> 💡 **All steps are orchestrated by `TutorialFlow`**, a PocketFlow-based DAG (Directed Acyclic Graph) defined in `flow.py`.

---

## **2.2 High-Level Architecture Diagram**

```
+------------------+
|   User Input     |
| (CLI: repo/dir)  |
+--------+---------+
         |
         v
+------------------+
|  main.py         | <--- Entrypoint: Parses args, sets up shared context
+--------+---------+
         |
         v
+------------------+
|  TutorialFlow    | <--- PocketFlow orchestrator (flow.py)
+--------+---------+
         |
         +------------------> [FetchRepo]
         |                         |
         |                         v
         |                [IdentifyAbstractions] → LLM → YAML
         |                         |
         |                         v
         |                [AnalyzeRelationships] → LLM → Summary + Graph
         |                         |
         |                         v
         |                [OrderChapters] → LLM → Chapter Order
         |                         |
         |                         v
         |                [WriteChapters] → LLM (Batch) → Per-chapter content
         |                         |
         |                         v
         +----------------> [CombineTutorial] → Final Tutorial (MD/HTML)
                                   |
                                   v
                         +------------------+
                         |  ./output/       | ← Generated tutorial
                         |  (by language)   |
                         +------------------+
```

> 🔗 **Data flows via a shared dictionary** (`shared`) passed between nodes. This includes:
> - `files`: List of `(path, content)` tuples
> - `abstractions`: Extracted concepts
> - `relationships`: How they interact
> - `chapters`: Ordered list of abstraction indices
> - `chapter_contents`: Generated explanations
> - `project_summary`, `language`, `repo_url`, etc.

---

## **2.3 Core Architectural Abstractions**

Let’s dive into the **7 key components** that make this system powerful, modular, and extensible.

### **1. `TutorialFlow` – The Orchestrator**
- **Defined in**: `flow.py`, `main.py`
- **Framework**: [PocketFlow](https://github.com/The-Pocket/PocketFlow) (lightweight agentic workflow engine)
- **Role**: Chains nodes together using `>>` syntax and manages execution order.
- **Features**:
  -