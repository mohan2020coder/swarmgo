# PocketFlow Tutorial: Building AI-Powered Codebase Tutorials

## Chapter 1: Getting Started with PocketFlow

### Introduction
PocketFlow helps developers quickly understand new codebases by generating structured tutorials. This chapter gets you up and running with the basic setup and first tutorial generation.

### Key Components

1. **README.md**: The project's documentation hub showing:
   - Supported languages (Python, JavaScript, etc.)
   - Example commands for GitHub and local analysis
   - Online service availability (code2tutorial.com)

2. **main.py**: The command-line interface that:
   - Accepts --repo or --dir arguments
   - Handles file pattern filtering
   - Manages output generation

3. **Dockerfile**: Provides containerization with:
   - Python 3.10 base image
   - Git support for repository operations
   - Clean image size optimization

### Code Walkthrough

```english
# Example Docker build and run
docker build -t pocketflow .
docker run -v $(pwd)/output:/app/output pocketflow --repo https://github.com/user/repo
```

```english
# Basic CLI usage from main.py
parser = argparse.ArgumentParser()
parser.add_argument('--repo', help='GitHub repository URL')
parser.add_argument('--dir', help='Local directory path')
args = parser.parse_args()
```

### Step-by-Step Setup

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Generate your first tutorial:
   ```bash
   python main.py --repo https://github.com/example/sample-repo --output my_first_tutorial
   ```

3. Check the generated tutorial in the output directory.

### Why It Works
The Dockerfile creates an isolated environment with all dependencies, while main.py orchestrates the analysis process through its argument parsing and initial configuration.

### Common Pitfalls
- Forgetting to mount volumes when using Docker (-v flag)
- Hitting GitHub API rate limits without authentication
- Overlooking file size limits (default 100KB)

### Hands-on Exercise
Generate a tutorial for a small local Python project:
```bash
python main.py --dir ./my_project --include "*.py" --exclude "tests/*"
```

## Chapter 2: Understanding the Core Pipeline

### Introduction
The pipeline architecture is what makes PocketFlow modular and extensible. This chapter explores how the workflow connects processing nodes to transform code into tutorials.

### Key Components

1. **flow.py**: Defines the sequence:
   ```english
   def create_tutorial_flow():
       return Flow() \
           >> FetchRepo() \
           >> IdentifyAbstractions(max_retries=5, wait=20) \
           >> AnalyzeRelationships(max_retries=5, wait=20) \
           >> OrderChapters(max_retries=5, wait=20) \
           >> WriteChapters() \
           >> CombineTutorial()
   ```

2. **nodes.py**: Contains the node implementations that:
   - Process repository files
   - Identify key abstractions
   - Determine learning order

### Pipeline Diagram

```
[FetchRepo] → [IdentifyAbstractions] → [AnalyzeRelationships]
      ↓
[OrderChapters] → [WriteChapters] → [CombineTutorial]
```

### Step-by-Step Execution

1. Pipeline starts with FetchRepo collecting files
2. IdentifyAbstractions detects key components
3. AnalyzeRelationships maps connections
4. OrderChapters structures the learning path
5. WriteChapters generates content
6. CombineTutorial produces final output

### Why It Works
The flow design allows each node to focus on one responsibility while passing state through the shared dictionary. Retry logic (max_retries/wait) makes the pipeline resilient to temporary failures.

### Common Pitfalls
- Modifying node order without understanding dependencies
- Not accounting for shared state requirements between nodes
- Overlooking the batch nature of WriteChapters

### Hands-on Exercise
Modify flow.py to add logging between nodes:
```python
class LogNode(Node):
    def run(self, shared):
        print(f"Processing {shared.get('current_step', 'unknown')}")
        
def create_tutorial_flow():
    return Flow() >> LogNode() >> FetchRepo() >> LogNode() >> ...
```

## Chapter 3: Working with File Crawlers

### Introduction
File crawlers are PocketFlow's eyes - they collect source material for analysis. This chapter covers both GitHub and local file crawling with their filtering capabilities.

### Key Components

1. **crawl_github_files.py**:
   - Uses GitHub API or Git cloning
   - Handles authentication tokens
   - Implements size and pattern filters

2. **crawl_local_files.py**:
   - Recursively scans directories
   - Respects .gitignore
   - Provides progress reporting

### Code Examples

```english
# GitHub crawling with filters
result = crawl_github_files(
    repo_url,
    token=os.getenv('GITHUB_TOKEN'),
    include_patterns={"*.py", "*.md"},
    max_file_size=50000
)
```

```english
# Local crawling with progress
files = crawl_local_files(
    "./project",
    include_patterns={"*.py"},
    exclude_patterns={"*test*"},
    show_progress=True
)
```

### Pattern Matching Table

| Pattern | Matches |
|---------|---------|
| `*.py`  | All Python files |
| `tests/*` | Files in tests directory |
| `*_test.py` | Test files |

### Step-by-Step Configuration

1. For GitHub repos:
   ```bash
   export GITHUB_TOKEN=your_token
   python main.py --repo https://github.com/org/repo --include "*.py" "*.md"
   ```

2. For local projects:
   ```bash
   python main.py --dir ./project --exclude "node_modules/*" "*.min.js"
   ```

### Why It Works
The crawlers use fnmatch for pattern matching and pathspec for .gitignore support. GitHub API handling includes rate limit awareness and retry logic.

### Common Pitfalls
- Forgetting to set GITHUB_TOKEN for private repos
- Overly broad patterns slowing down crawling
- Not accounting for base64 encoding overhead in size calculations

### Hands-on Exercise
Create a custom pattern filter that only includes files with "service" in their name:
```python
def custom_filter(path):
    return "service" in path.lower()
    
files = crawl_local_files("./project", file_filter=custom_filter)
```

## Chapter 4: The LLM Service Integration

### Introduction
The LLM service powers PocketFlow's analysis capabilities. This chapter explores how to configure and optimize LLM interactions.

### Key Components

**call_llm.py** features:
- Multi-provider support (Gemini default)
- Response caching
- Detailed logging
- Error handling

### Code Examples

```english
# Basic LLM call with caching
response = call_llm(
    "Explain this code: print('hello')",
    use_cache=True
)

# Configuration via environment
os.environ["GEMINI_MODEL"] = "gemini-1.5-pro"
os.environ["LLM_PROVIDER"] = "azure"
```

### LLM Providers Table

| Provider | Config Var | Example Model |
|----------|------------|---------------|
| Gemini | GEMINI_API_KEY | gemini-pro |
| Azure | AZURE_API_KEY | gpt-4 |
| OpenAI | OPENAI_KEY | gpt-3.5-turbo |

### Step-by-Step Configuration

1. Set your API key:
   ```bash
   export GEMINI_API_KEY=your_key_here
   ```

2. Test the LLM connection:
   ```python
   from utils.call_llm import call_llm
   print(call_llm("Say hello!"))
   ```

3. Monitor logs:
   ```bash
   tail -f logs/llm_*.log
   ```

### Why It Works
The service abstracts provider differences behind a consistent interface. Caching prevents duplicate calls for identical prompts, saving costs and time.

### Common Pitfalls
- Not setting API keys before running
- Forgetting cache can lead to unnecessary charges
- Overlooking provider-specific rate limits

### Hands-on Exercise
Implement a custom prompt template for code analysis:
```python
def analyze_code_prompt(code):
    return f"""Analyze this code for key concepts:
    
{code}

Identify:
1. Main components
2. Key functions
3. Architectural patterns"""
```

## Chapter 5: Deep Dive into Processing Nodes

### Introduction
Processing nodes are PocketFlow's brain - each specializes in one part of tutorial generation. This chapter examines each node's role and customization points.

### Node Responsibilities

1. **FetchRepo**: Gathers source files
2. **IdentifyAbstractions**: Finds key components
3. **AnalyzeRelationships**: Maps connections
4. **OrderChapters**: Structures learning path
5. **WriteChapters**: Generates content

### Code Examples

```english
# IdentifyAbstractions core logic
abstractions = call_llm(
    f"Identify abstractions in: {file_contents}",
    use_cache=True
)
shared['abstractions'] = validate_abstractions(abstractions)
```

```english
# OrderChapters learning path
chapters = call_llm(
    f"Order these concepts for learning: {abstractions}",
    use_cache=True
)
shared['chapter_order'] = parse_chapter_order(chapters)
```

### Node Characteristics Table

| Node | Retries | Batch | Key Output |
|------|---------|-------|------------|
| FetchRepo | 0 | No | shared['files'] |
| IdentifyAbstractions | 5 | No | shared['abstractions'] |
| WriteChapters | 5 | Yes | shared['chapters'] |

### Step-by-Step Customization

1. Add a new abstraction validator:
   ```python
   class CustomIdentifyAbstractions(IdentifyAbstractions):
       def validate(self, abstractions):
           return [a for a in abstractions if a['complexity'] > 3]
   ```

2. Update the flow:
   ```python
   flow = Flow() >> FetchRepo() >> CustomIdentifyAbstractions() >> ...
   ```

### Why It Works
Each node focuses on a single responsibility while sharing state through the shared dictionary. Retry logic makes the pipeline resilient to transient failures.

### Common Pitfalls
- Modifying shared state keys without updating dependent nodes
- Not maintaining consistent abstraction formats
- Overriding core validation logic incorrectly

### Hands-on Exercise
Create a custom node that counts code lines:
```python
class CountLinesNode(Node):
    def run(self, shared):
        files = shared['files']
        total_lines = sum(f.count('\n') for f in files.values())
        shared['line_count'] = total_lines
```

## Chapter 6: Customizing the Pipeline Flow

### Introduction
PocketFlow's pipeline is designed for customization. This chapter shows how to modify the workflow for specialized use cases.

### Customization Points

1. **Node Order**: Reorganize processing steps
2. **Custom Nodes**: Inject new functionality
3. **Parallel Processing**: Speed up bottlenecks

### Code Examples

```english
# Adding a preprocessing node
class PreprocessNode(Node):
    def run(self, shared):
        shared['files'] = {k: v for k,v in shared['files'].items() 
                          if not k.startswith('migrations')}

flow = Flow() >> FetchRepo() >> PreprocessNode() >> IdentifyAbstractions()...
```

```english
# Parallel chapter writing
from concurrent.futures import ThreadPoolExecutor

class ParallelWriteChapters(WriteChapters):
    def run(self, shared):
        with ThreadPoolExecutor() as executor:
            results = list(executor.map(self.write_chapter, shared['chapter_order']))
        shared['chapters'] = results
```

### Flow Modification Patterns

| Pattern | Use Case | Example |
|---------|----------|---------|
| Pre-processing | Filtering inputs | Remove test files |
| Post-processing | Formatting output | Add table of contents |
| Parallelization | Performance | Concurrent chapter writing |

### Step-by-Step Customization

1. Identify bottleneck (e.g., WriteChapters)
2. Create custom node version
3. Update flow.py to use new node
4. Test with sample repository

### Why It Works
The Flow class manages node execution order while preserving shared state. The >> operator provides clean syntax for composing nodes.

### Common Pitfalls
- Breaking shared state expectations
- Creating circular dependencies between nodes
- Not maintaining node interface contracts

### Hands-on Exercise
Create a flow that generates both beginner and advanced tutorials:
```python
flow = (Flow() >> FetchRepo() >> IdentifyAbstractions()
        >> OrderChapters(difficulty='beginner') >> WriteChapters()
        >> OrderChapters(difficulty='advanced') >> WriteChapters())
```

## Chapter 7: Advanced Configuration and Extension

### Introduction
This chapter covers professional-grade extensions - custom crawlers, specialized prompts, and advanced pipeline configurations.

### Extension Patterns

1. **Custom Crawlers**: Support new sources
2. **Prompt Engineering**: Optimize LLM outputs
3. **Analysis Plugins**: Add new metrics

### Code Examples

```english
# Custom GitLab crawler
class GitLabCrawler:
    def fetch_files(self, repo_url):
        # Implement GitLab API logic
        return files

# In main.py:
if args.repo.startswith('gitlab'):
    shared['files'] = GitLabCrawler().fetch_files(args.repo)
```

```english
# Architecture analysis plugin
class ArchitectureAnalyzer(Node):
    def run(self, shared):
        analysis = call_llm(f"Analyze architecture: {shared['files']}")
        shared['architecture'] = parse_analysis(analysis)
```

### Extension Point Table

| Extension Type | Interface | Example |
|----------------|-----------|---------|
| File Source | dict[str,str] | Bitbucket crawler |
| Analysis Node | Node subclass | Security scanner |
| Output Format | str -> str | Markdown converter |

### Step-by-Step Extension

1. Identify extension point
2. Create conforming implementation
3. Integrate with existing flow
4. Add configuration options

### Why It Works
PocketFlow's modular design isolates components, making extensions safe. The shared dictionary provides a standard integration point.

### Common Pitfalls
- Not maintaining interface contracts
- Overcomplicating simple extensions
- Not properly handling errors in custom code

### Hands-on Exercise
Create a node that generates code diagrams:
```python
class DiagramGenerator(Node):
    def run(self, shared):
        for chapter in shared['chapters']:
            chapter['diagram'] = generate_diagram(chapter['content'])
```

## Chapter 8: Production Deployment Patterns

### Introduction
Taking PocketFlow to production requires optimization and scaling. This chapter covers Docker optimizations, performance tuning, and scaling strategies.

### Production Considerations

1. **Docker Optimization**:
   - Multi-stage builds
   - Layer caching
   - Minimal base images

2. **Performance**:
   - LLM call batching
   - Process isolation
   - Caching strategies

### Code Examples

```english
# Optimized Dockerfile
FROM python:3.10-slim as builder
COPY requirements.txt .
RUN pip install --user -r requirements.txt

FROM python:3.10-slim
COPY --from=builder /root/.local /root/.local
COPY . /app
ENV PATH=/root/.local/bin:$PATH
```

```english
# Batch LLM processing
class BatchLLMNode(Node):
    def run(self, shared):
        batch = [f"Analyze: {f[:1000]}" for f in shared['files'].values()]
        shared['analysis'] = call_llm_batch(batch)
```

### Optimization Techniques

| Technique | Benefit | Implementation |
|-----------|---------|-----------------|
| Multi-stage builds | Smaller images | Dockerfile stages |
| LLM batching | Reduced API calls | BatchNode subclass |
| Memory caching | Faster repeats | @lru_cache |

### Step-by-Step Optimization

1. Analyze performance bottlenecks
2. Implement Docker optimizations
3. Add batch processing where possible
4. Set up monitoring

### Why It Works
Container optimizations reduce deployment footprint while batching and caching improve throughput. Process isolation prevents failures from cascading.

### Common Pitfalls
- Not setting memory limits in Docker
- Over-batching leading to timeouts
- Not monitoring production usage

### Hands-on Exercise
Create a production-ready Dockerfile with:
1. Multi-stage build
2. Non-root user
3. Health check
4. Resource limits

## Where to Go Next

### Project Extensions
1. Add support for additional version control systems (GitLab, Bitbucket)
2. Implement a web interface using FastAPI or Flask
3. Create IDE plugins (VS Code, PyCharm)
4. Add automated testing for tutorial quality
5. Develop specialized analysis modes (security, performance)

### Learning Resources
1. Explore the Gemini API documentation for advanced prompt engineering
2. Study software architecture patterns to improve analysis quality
3. Learn about educational psychology for better tutorial structuring
4. Research CI/CD pipelines for automated tutorial updates

### Community
1. Contribute to the open-source project
2. Share your custom nodes and flows
3. Create issue reports for bugs or feature requests
4. Write about your PocketFlow use cases

## Summary Recap

Throughout this tutorial, you've learned:

1. **Setup**: How to install and run PocketFlow locally and via Docker
2. **Architecture**: The pipeline-based design with specialized nodes
3. **File Handling**: Working with both GitHub and local file crawlers
4. **LLM Integration**: Configuring and optimizing language model interactions
5. **Node Deep Dive**: The responsibilities of each processing node
6. **Customization**: How to modify and extend the pipeline flow
7. **Advanced Patterns**: Professional-grade extensions and configurations
8. **Production**: Deployment optimizations and scaling strategies

You're now equipped to use, customize, and extend PocketFlow for your codebase documentation needs!