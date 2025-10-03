# LLM Integration Examples (Ollama & LongCat)

This file contains example HTTP payloads and notes for calling Ollama (local) and LongCat (example cloud).

## Ollama (local - example)
Ollama exposes local models via HTTP (example). Adjust the endpoint & model name per your setup.

Endpoint (example):
```
POST http://localhost:11434/completions
Content-Type: application/json
```
Payload example:
```json
{
  "model": "llama2-13b",
  "prompt": "Create a detailed outline for a document about: Document Builder Project Plan. Audience: developer. Level: detailed",
  "max_tokens": 800,
  "temperature": 0.2
}
```

Response (example):
```json
{
  "id": "...",
  "object": "text_completion",
  "choices": [
    {"text": "1. Introduction\n2. System Architecture\n..."} 
  ]
}
```

## LongCat (cloud - example)
LongCat (or other cloud-hosted LLM) will have a different API. Example payload might be:

Endpoint (example):
```
POST https://api.longcat.ai/v1/generate
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

Payload example:
```json
{
  "model": "longcat-1.0",
  "input": "Create a detailed outline for a document about: Document Builder Project Plan. Audience: developer. Level: detailed",
  "max_output_tokens": 800,
  "temperature": 0.2
}
```

Response example (simplified):
```json
{
  "id": "...",
  "output": "1. Introduction\n2. System Architecture\n..."
}
```

## Notes
- Ollama local payloads and endpoints vary by version; inspect `ollama --help` or the local docs.
- LongCat is an example name here — follow your cloud provider docs for exact request/response fields and authentication.
- Always cache prompt -> response to avoid repeated calls and save costs.
