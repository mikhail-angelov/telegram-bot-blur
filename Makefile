include .env

test:
	go test ./...
.PHONY: test

api:
	go run main.go
.PHONY: api

run:
	vercel dev
.PHONY: run

register:
	curl -F "url=${LAMBDA_URL}" https://api.telegram.org/bot${BOT_TOKEN}/setWebhook
.PHONY: register

unregister:
	curl https://api.telegram.org/bot${BOT_TOKEN}/deleteWebhook
.PHONY: unregister
