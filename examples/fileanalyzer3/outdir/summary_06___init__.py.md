# __init__.py - File Summary

**one_line**:  
Utility module initializer that imports and exposes key functions for shared use across the project.

**Purpose**:  
This `__init__.py` file serves as the entry point for the `utils` package, centralizing and exposing commonly used utility functions and classes to streamline imports throughout the codebase (e.g., `from utils import log, retry, ...`).

**Major functions/classes**:  
- Exposes the following (imported from submodules):
  - `log` – Standardized logging function with timestamps and levels.
  - `retry` – Decorator to retry a function on failure with exponential backoff.
  - `read_json`, `write_json` – Safe JSON file I/O with error handling.
  - `hash_string` – Utility to generate deterministic hash (e.g., for caching keys).
  - `Timer` – Context manager for measuring execution time.

**Key technical details & TODOs**:  
- Uses `from .logging import log`, `from .decorators import retry`, etc., to keep internal structure modular.
- Designed to avoid circular imports by lazy-loading heavy dependencies inside functions where possible.
- All exposed utilities are stateless and thread-safe.
- **TODO**: Add type hints to all public functions in submodules.
- **TODO**: Consider adding a `__all__` list to explicitly control what gets exported on `from utils import *`.

**Short usage example**:  
```python
from utils import log, retry, read_json, Timer

@retry(max_attempts=3)
def fetch_data():
    log("Fetching data...")
    with Timer("Data load"):
        return read_json("data.json")

data = fetch_data()
```