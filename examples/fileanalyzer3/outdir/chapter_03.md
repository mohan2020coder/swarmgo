# **Chapter 3: Setting Up Your Environment**

**Objective**: Prepare your system to run and experiment with the AI-powered codebase tutorial generator.

By the end of this chapter, you'll have a fully configured development environment capable of analyzing GitHub repositories or local codebases and generating beginner-friendly, multilingual tutorials — all powered by Large Language Models (LLMs) and the lightweight **PocketFlow** agent framework.

---

## **1. Prerequisites**

Before you begin, ensure the following tools are installed on your machine:

| Tool | Purpose | Installation |
|------|--------|--------------|
| **Python 3.10+** | Core runtime for the application | [python.org](https://www.python.org/downloads/) |
| **Git** | Clone repositories and manage source control | [git-scm.com](https://git-scm.com/) |
| **Docker** (optional) | Run the app in an isolated container | [docker.com](https://www.docker.com/products/docker-desktop/) |
| **LLM API Key** | Required for Gemini, Claude, or OpenAI | See [API Key Setup](#3-api-key-setup) |

> 💡 **Note**: While Docker simplifies setup, you can run everything directly with Python if preferred.

---

## **2. Clone the Project**

Start by cloning the repository:

```bash
git clone https://github.com/your-org/pocketflow-tutorial-generator.git
cd pocketflow-tutorial-generator
```

> 🔄 Replace the URL with the actual repo if different. The structure should include:
> - `main.py` – CLI entrypoint
> - `flow.py` – Pipeline orchestration
> - `nodes.py` – Core logic modules
> - `utils/` – Helper tools (`call_llm.py`, `crawl_github_files.py`, etc.)
> - `Dockerfile` – Container definition
> - `requirements.txt` – Python dependencies

---

## **3. API Key Setup**

The system uses LLMs (like **Google Gemini 2.5 Pro**, **Anthropic Claude 3.7**, or **OpenAI O1**) to analyze code and generate tutorials.

### ✅ **Recommended: Use Environment Variables**

Set your API key securely using environment variables:

#### For **Google Gemini** (default provider):
```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

#### For **Anthropic Claude**:
```bash
export ANTHROPIC_API_KEY="your-anthropic-api-key-here"
```

#### For **OpenAI / Azure / OpenRouter**:
```bash
export OPENAI_API_KEY="..."          # For OpenAI
# OR
export OPENROUTER_API_KEY="..."      # For OpenRouter
```

> 🔐 **Best Practice**: Never commit API keys to version control. Use `.env` files or shell profiles (e.g., `.zshrc`, `.bash_profile`) to store them.

> 🛠️ **Switching Providers**: Open `utils/call_llm.py` and comment/uncomment the desired `call_llm` implementation. The rest of the code remains unchanged.

---

## **4. Install Dependencies**

### Option A: **Native Python (Recommended for Development)**

Install required packages:

```bash
pip install -r requirements.txt
```

> 📦 Key dependencies include:
> - `google-generativeai`, `openai`, `anthropic` – LLM clients
> - `requests` – GitHub API access
> - `PyYAML` – Parse LLM responses
> - `pathspec` – `.gitignore` support
> - `PocketFlow` – Lightweight workflow engine

### Option B: **Docker (Recommended for Reproducibility & Production)**

Build the image using the provided `Dockerfile`:

```bash
docker build -t pocketflow-tutorial-gen .
```

> ✅ The Dockerfile:
> - Uses `python:3.10-slim` for minimal footprint
> - Installs `git` for GitHub cloning
> - Copies code and installs Python packages
> - Sets `main.py` as the entrypoint

> ⚠️ **Security Note**: By default, the container runs as root. For production, consider adding a non-root user:
> ```Dockerfile
> RUN useradd -m -s /bin/bash appuser
> USER appuser
> ```

> 📝 **Tip**: Add a `.dockerignore` file to exclude logs, cache, and IDE files:
> ```
> __pycache__
> *.pyc
> .git
> .vscode
> logs/
> llm_cache.json
> ```

---

## **5. Configure Input & Output**

The