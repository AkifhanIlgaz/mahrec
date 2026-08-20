.PHONY: generate css build dev

generate:
	templ generate

css:
	npx @tailwindcss/cli -i input.css -o internal/web/static/css/app.css

build: generate css
	go build -o bin/server ./cmd/server

dev: generate css
	go run ./cmd/server
