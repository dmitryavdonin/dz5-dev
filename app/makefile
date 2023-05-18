run: ## runs app localy
	@go run cmd/app/main.go

commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/app:latest . && docker push dmitryavdonin/auth:latest
