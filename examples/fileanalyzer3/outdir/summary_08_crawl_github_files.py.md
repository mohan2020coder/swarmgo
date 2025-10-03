# crawl_github_files.py - File Summary

### one_line  
A utility to crawl and download files from a GitHub repository (public or private) at a specific commit/path, with support for filtering, size limits, and relative path handling.

---

### Purpose  
This script enables programmatic retrieval of source code or documentation files from a GitHub repository via the GitHub API (or SSH clone), with fine-grained control over which files are downloaded based on patterns, file size, and path. It is useful for tools that need to analyze or process repository contents (e.g., code generation, documentation extraction, CI/CD automation).

---

### Major Functions/Classes

#### `crawl_github_files(...)`
Main function that orchestrates crawling:
- Parses GitHub URL (including commit/branch and subdirectory).
- Handles both **HTTPS** (via GitHub API) and **SSH** (via `git clone`) URLs.
- Downloads files recursively from a specified subdirectory.
- Applies **include/exclude glob patterns** to filter files.
- Enforces **maximum file size** to avoid memory issues.
- Returns a dictionary with file contents and crawl statistics.

#### `should_include_file(file_path, file_name)`
Helper to evaluate inclusion based on `include_patterns` and `exclude_patterns` using `fnmatch`.

#### `fetch_contents(path)` (nested)
Recursive function to traverse GitHub repository structure using the GitHub REST API (`/contents` endpoint), downloading files and skipping directories based on filters.

#### `fetch_branches()` and `check_tree()`
Used to validate and extract commit/branch reference from the URL when not explicitly provided.

---

### Key Technical Details & TODOs

#### ✅ **Features**
- **Supports both public and private repositories** via GitHub token (`GITHUB_TOKEN` env or argument).
- **Handles rate limits** by pausing and retrying when hitting GitHub API limits.
- **Two URL modes**:
  - HTTPS: Uses GitHub API (supports branch/commit + path).
  - SSH: Falls back to `git clone` into temp directory (no branch parsing in URL; uses default).
- **Flexible filtering**:
  - `include_patterns`: e.g., `{"*.py", "*.md"}`
  - `exclude_patterns`: e.g., `{"*test*", "*.log"}`
- **Relative paths option**: Strips base directory prefix if `use_relative_paths=True`.
- **Robust error handling** for 403, 404, encoding, file I/O, and size checks.

#### ⚠️ **Technical Notes**
- **SSH URLs do not support branch/commit parsing** from the URL — relies on default branch.
- **File size is checked twice**: once from metadata, once from `download_url` headers.
- **Base64 decoding** used when `download_url` is missing (fallback).
- **Recursive directory traversal** with early pruning using exclude patterns (optimized in new implementation).
- Uses **`tempfile.TemporaryDirectory`** for SSH cloning (auto-cleaned).

#### 🔧 **TODOs / Improvements**
- [ ] **Add support for SSH URLs with branch/commit** (e.g., `git@github.com:owner/repo.git#commit/path`).
- [ ] **Add progress bar or logging** for large repositories.
- [ ] **Support for shallow clones** in SSH mode to reduce bandwidth.
- [ ] **Caching mechanism** to avoid re-downloading unchanged files.
- [ ] **Better handling of symlinks and binary files** (currently skipped silently).
- [ ] **Validate `include_patterns`/`exclude_patterns` syntax** before use.
- [ ] **Add retry logic** for transient network failures.

---

### Short Usage Example

```python
import os
from crawl_github_files import crawl_github_files

# Optional: Set token for private repos or to avoid rate limits
os.environ["GITHUB_TOKEN"] = "your_token_here"

result = crawl_github_files(
    repo_url="https://github.com/microsoft/autogen/tree/main/python/packages/autogen-core",
    token=None,  # uses GITHUB_TOKEN env
    max_file_size=500 * 1024,  # 500 KB
    use_relative_paths=True,
    include_patterns={"*.py", "*.md"},
    exclude_patterns={"*test*"}
)

print(f"Downloaded {result['stats']['downloaded_count']} files")
for path, content in result["files"].items():
    print(f"{path}: {len(content)} chars")
```

> 💡 **Tip**: Use `include_patterns` and `exclude_patterns` to focus on relevant files (e.g., source code, docs) and skip binaries, logs, or tests.