# **Chapter 4: Deep Dive – Fetching the Codebase**  
**Objective**: Understand how the system ingests and filters source code.

Before any AI can analyze a codebase, it must first **fetch** the right files — not too many, not too few — in a way that's efficient, safe, and reproducible. This chapter takes you behind the scenes of how the system retrieves code from **GitHub repositories** or **local directories**, applies intelligent filtering, and prepares it for analysis.

We’ll explore:
- How the system fetches code from GitHub or your local machine
- How it filters files using patterns and `.gitignore`
- How file size and structure are managed to avoid overload
- The role of caching and resilience in the pipeline

Let’s dive into the **core ingestion engine** of the tutorial generator.

---

## 🔍 **1. The Entry Point: `FetchRepo` Node**

At the start of every tutorial generation flow is the `FetchRepo` class in `nodes.py`. This is the **gatekeeper** of the pipeline — responsible for retrieving source files and passing them downstream.

```python
class FetchRepo(Node):
    def prep(self, shared):
        # Get configuration from shared context
        repo_url = shared.get("repo_url")
        local_dir = shared.get("dir")
        include_patterns = shared.get("include_patterns", ["*.py", "*.js", "*.ts", "*.go", "*.rs", "*.md", "*.txt", "*.html", "*.css", "*.json", "*.yaml", "*.yml"])
        exclude_patterns = shared.get("exclude_patterns", ["*test*", "*spec*", "*.log", "__pycache__/*", "node_modules/*", ".git/*", ".venv/*", "dist/*", "build/*"])
        max_file_size = shared.get("max_file_size", 100 * 1024)  # 100 KB default
        use_relative_paths = True

        # Choose source: GitHub or local
        if repo_url:
            return crawl_github_files(
                repo_url=repo_url,
                token=shared.get("github_token"),
                max_file_size=max_file_size,
                use_relative_paths=use_relative_paths,
                include_patterns=set(include_patterns),
                exclude_patterns=set(exclude_patterns)
            )
        else:
            return crawl_local_files(
                directory=local_dir,
                include_patterns=set(include_patterns),
                exclude_patterns=set(exclude_patterns),
                max_file_size=max_file_size,
                use_relative_paths=use_relative_paths
            )
```

> ✅ **Key Insight**: `FetchRepo` is **source-agnostic**. Whether the code lives on GitHub or your laptop, the same filtering logic applies.

---

## 🌐 **2. Fetching from GitHub: Smart Crawling with API or Git**

When a GitHub URL is provided (e.g., `--repo https://github.com/encode/fastapi`), the system uses `crawl_github_files.py` to retrieve files.

### 🔧 **Two Modes of Operation**

| Mode | When Used | How It Works |
|------|-----------|-------------|
| **GitHub API Mode (HTTPS URL)** | Default for public/private repos | Uses GitHub’s `/contents` REST API to list and download files recursively. Supports branches, commits, and subdirectories (e.g., `tree/main/src`). |
| **Git Clone Mode (SSH URL)** | Fallback for SSH or complex cases | Clones the repo into a temporary directory using `git`, then reads files locally. Less efficient but more reliable for private repos with SSH access. |

### 🛡️ **Security & Rate Limiting**
- Uses `GITHUB_TOKEN` (from env or CLI) to access private repos and avoid rate limits.
- Automatically pauses and retries when hitting GitHub API rate limits (403 responses).
- Validates file size **before** downloading to prevent memory bloat.

### 🎯 **Pattern Filtering**
Files are filtered using **glob-style patterns**:
```python
include_patterns = {"*.py", "*.md"}
exclude_patterns = {"*test*", "*.log"}
```
This ensures only relevant source and documentation files are included — skipping tests, logs, binaries, and build artifacts.

> 💡 **Pro Tip**: You can customize these patterns via CLI flags:
> ```bash
> --include "*.py" "*.md" --exclude "tests/*" "*.min.js"
> ```

---

## 💾 **3. Fetching from Local Directories: `.gitignore`-Aware Crawling**

When analyzing a local project