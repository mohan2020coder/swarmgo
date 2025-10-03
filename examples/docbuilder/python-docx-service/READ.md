sudo apt install python3-pip
pip install flask python-docx
pip install -r requirements.txt
python3 app.py


-----------------


curl -X POST http://localhost:5000/generate \
  -H "Content-Type: application/json" \
  --data @payload.json

curl -X POST http://127.0.0.1:5000/generate \
  -H "Content-Type: application/json" \
  --data @payload.json


response

{
  "path": "/tmp/docx_outputs/7a6d9f1e96c04a9a9f4a1c57e1f1e2b8.docx"
}

download

curl "http://localhost:5000/download?path=/tmp/docx_outputs/7a6d9f1e96c04a9a9f4a1c57e1f1e2b8.docx" --output professional_report.docx

curl "http://localhost:5000/download?path=/tmp/docx_outputs/2b544683425340d8b18a510c01f61546.docx" --output professional_report.docx

/tmp/docx_outputs/2b544683425340d8b18a510c01f61546.docx


curl "http://localhost:5000/download?path=/tmp/docx_outputs/61ae9a2054c1410c938580478083cd8f.docx" --output professional_report1.docx


curl "http://localhost:5000/download?path=/tmp/docx_outputs/3d787827ea8b45fd898b96bc8b019eff.docx" --output professional_report1.docx

