.PHONY: dep build docker release install test backup

build:
	go build

run: build
	./example

docker:
	docker build -t example .

run-docker: docker
	docker run --rm -e G_KEY -e G_PROJECT_ID -e PORT=8000 -p 8000:8000 example
