# Document Builder - Repo Scaffold (Demo)

This scaffold provides a minimal, runnable demo of the **Document Builder**:
- Go-based orchestrator and agents (skeleton)
- Python DOCX service (Flask + python-docx)
- SQLite caching
- Diagram rendering pipeline stubs (Mermaid/PlantUML)
- Docker Compose to run everything locally

## Quickstart (development)

1. Build and run with Docker Compose:
```bash
docker-compose build
docker-compose up
```

2. Submit a job:
```bash
curl -X POST http://localhost:8080/generate -H "Content-Type: application/json" -d '{
  "topic": "Document Builder Project Plan",
  "format": "docx",
  "include_diagrams": true,
  "audience": "developer",
  "detail_level": "detailed"
}'
```

3. Download the returned path via Python service `/download` endpoint (the orchestrator returns the path of the generated file).

## What is included
- `orchestrator/` - Go server and agent stubs
- `python-docx-service/` - Flask docx generator
- `docker-compose.yml` and Dockerfiles
- `integration_docs/` - example HTTP payloads for Ollama and LongCat
- `diagram_pipeline/` - helper scripts and notes for rendering mermaid & PlantUML

## Notes
This scaffold is purposely minimal to be a starter. Read `integration_docs/` for exact LLM payloads used with Ollama and LongCat (examples included).

