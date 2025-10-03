# call_llm.py - File Summary

**one_line**:  
Unified LLM API caller with caching, logging, and multi-provider support (Gemini, Azure, Anthropic, OpenAI, OpenRouter).

**Purpose**:  
Provides a single, consistent interface to call various Large Language Model (LLM) APIs, with built-in prompt/response logging, disk-based caching, and environment-based configuration. Designed for use in larger workflows (e.g., PocketFlow) to reduce costs and latency via caching and improve observability via structured logging.

**Major functions/classes**:
- `call_llm(prompt: str, use_cache: bool = True) -> str`  
  The main function to send a prompt to an LLM and get a response. Supports caching, logging, and currently defaults to **Google Gemini 2.5 Pro** via API key.
- **Commented alternative implementations** for:
  - Azure OpenAI
  - Anthropic Claude 3.7 Sonnet (with extended thinking)
  - OpenAI o1
  - OpenRouter (supports any model via OpenRouter API)
- **Logging setup**: All prompts and responses are logged to a daily rotating file in the `logs/` directory.
- **Caching mechanism**: Uses `llm_cache.json` to store prompt-response pairs, avoiding redundant API calls.

**Key technical details & TODOs**:
- ✅ **Default provider**: Google Gemini via `genai.Client(api_key=...)` using `GEMINI_API_KEY` env var.
- ✅ **Caching**: Simple JSON file cache (`llm_cache.json`) keyed by prompt. Thread-unsafe but sufficient for single-threaded use.
- ✅ **Logging**: Logs to `logs/llm_calls_YYYYMMDD.log` with timestamps, levels, and full prompt/response.
- 🔁 **Multi-provider support**: Easily switch providers by uncommenting the desired implementation and setting relevant environment variables.
- 🛑 **Security**: API keys are read from environment variables (recommended), but fallbacks to hardcoded values (e.g., `"your-project-id"`) — **TODO: Remove or warn about insecure defaults**.
- 🧪 **Testing**: Includes a simple `__main__` block to test the function.
- ⚠️ **Cache race condition**: Cache is reloaded before write, but concurrent access could still cause issues.
- 📦 **Dependencies**: Requires `google-generativeai`, and optionally `openai`, `anthropic`, or `requests` for other providers.
- 🔧 **Configurable via env vars**:
  - `LOG_DIR`: Log directory (default: `logs`)
  - `GEMINI_API_KEY`, `GEMINI_MODEL`, `GEMINI_PROJECT_ID`, `GEMINI_LOCATION`
  - `OPENROUTER_API_KEY`, `OPENROUTER_MODEL`
  - `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, etc.
- 🧹 **TODO**: Add cache TTL, size limit, or hash-based keys to avoid JSON size issues with long prompts.
- 🧹 **TODO**: Add error handling and retry logic for API failures.

**Short usage example**:
```python
from utils.call_llm import call_llm

# Set environment variables first (e.g., GEMINI_API_KEY)
response = call_llm("Explain quantum computing in simple terms")
print(response)

# Disable cache for fresh call
response_fresh = call_llm("Hello", use_cache=False)
```

> 💡 **Tip**: Switch providers by commenting out the current `call_llm` and uncommenting another — ensure required env vars and packages are set up.