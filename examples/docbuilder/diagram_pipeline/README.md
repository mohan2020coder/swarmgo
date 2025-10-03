# Diagram Rendering Pipeline (Mermaid & PlantUML)

This directory contains notes and helper scripts to render Mermaid or PlantUML diagrams and produce images for embedding in DOCX.

## Options to render diagrams
- **Mermaid CLI (mmdc)**: install via npm (npm i -g @mermaid-js/mermaid-cli)
  - Render: `mmdc -i diagram.mmd -o diagram.png`
- **PlantUML**: use the plantuml.jar or a PlantUML server. Requires Java.
  - Render local: `java -jar plantuml.jar -tpng diagram.puml`

## Example helper script (render_mermaid.sh)
```bash
#!/bin/bash
set -e
INPUT=$1
OUT=$2
mmdc -i "$INPUT" -o "$OUT"
```

## Embedding in DOCX
- Render PNG/SVG to a temporary path and include via `doc.add_picture('/tmp/diagram.png', width=Inches(6))` in the python-docx service.
- Sanitize and validate inputs before writing to disk and running CLI to avoid command injection.
