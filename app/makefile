run: ## runs app localy
	@go run cmd/app/main.go

commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/dz5-app:2.0.0 . && docker push dmitryavdonin/dz5-app:2.0.0
