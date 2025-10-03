# crawl_local_files.py - File Summary

### one_line
Crawls local directories to collect file contents with pattern-based inclusion/exclusion, `.gitignore` support, size limits, and progress tracking.

---

### Purpose
This utility mimics a GitHub file crawler for **local filesystems**, enabling selective reading of files based on include/exclude patterns, `.gitignore` rules, and file size constraints. It is designed for use in code analysis, automation, or AI training data preparation tools where controlled file ingestion is required.

---

### Major functions/classes

- **`crawl_local_files()`**  
  Main function that:
  - Walks a directory tree (`os.walk`)
  - Applies `.gitignore` rules (if present)
  - Filters files using `include_patterns` and `exclude_patterns` via `fnmatch`
  - Enforces a maximum file size
  - Reads and returns file contents as a dictionary
  - Shows real-time progress with colored output

---

### Key technical details & TODOs

#### ✅ **Features**
- **`.gitignore` Support**: Uses `pathspec.PathSpec` to respect `.gitignore` rules (both files and directories).
- **Pattern Matching**: Uses `fnmatch` for glob-style patterns (e.g., `*.py`, `tests/*`).
- **Early Directory Pruning**: Skips entire directories during `os.walk` if excluded (improves performance).
- **Progress Feedback**: Prints real-time progress (file count, percentage, status) in green (`\033[92m`).
- **Relative Paths**: Optionally returns paths relative to input directory.
- **Robust Encoding**: Uses `utf-8-sig` to handle BOM in files.
- **Error Handling**: Gracefully skips unreadable files and logs warnings.

#### ⚠️ **Limitations / TODOs**
- **No Symlink Handling**: Follows symlinks by default (potential cycles or unintended reads).
- **No Binary Detection**: Attempts to read all files as text — may fail on binaries.
- **Progress Overhead**: Frequent `print()` calls may slow down large crawls; consider optional verbosity flag.
- **Hardcoded Color Codes**: ANSI colors may not render well in all terminals.
- **No Async Support**: Synchronous I/O; not suitable for huge repositories without optimization.

> 🔧 **Suggested Improvements (TODOs):**
> - Add `verbose`/`quiet` mode flag to control output.
> - Add option to detect and skip binary files.
> - Support for custom `.gitignore`-like files (e.g., `.crawlignore`).
> - Return metadata (size, mtime) alongside content.
> - Allow custom file read handlers (e.g., for line filtering).

---

### Short usage example

```python
# Crawl current directory, include only Python files, exclude test and cache dirs
result = crawl_local_files(
    directory=".",
    include_patterns={"*.py", "*.md"},
    exclude_patterns={"tests/*", "__pycache__/*", "*.log"},
    max_file_size=1024 * 1024,  # 1 MB limit
    use_relative_paths=True
)

# Print file paths
for path in result["files"]:
    print(path)

# Use file content
for path, content in result["files"].items():
    print(f"--- {path} ---\n{content[:200]}...\n")
```

> 💡 **Tip**: Run as a script (`python crawl_local_files.py`) to crawl the parent directory with default exclusions (e.g., `.git`, `.venv`).