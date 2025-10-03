.PHONY: build up orchestrator-run py-run clean

build:
	docker-compose build

up:
	docker-compose up

orchestrator-run:
	cd orchestrator && go run ./cmd/server

py-run:
	cd python-docx-service && python app.py

clean:
	rm -rf outputs cache/*.db
