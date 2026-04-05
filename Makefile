build:
	@echo "Building chatly backend"
	go build -o ./bin/chatly cmd/server/main.go
	@echo "Building chatly backend done"

clean:
	@echo "Cleaning chatly backend"
	rm -r ./bin
	@echo "Cleaning chatly backend done"

run:
	@echo "Running chatly backend"
	./bin/chatly
	@echo "Running chatly backend done"

build-docker-image:
	@echo "Building chatly backend docker image"
	docker build -t your_username/chatly .
	@echo "Building chatly backend docker image done"

docker-run:
	@echo "Running chatly backend docker image"
	docker run -e PORT=8080 -e MONGO_URL="your-mongo-url" -e JWT_SECRET="your-jwt-secret" -e JWT_TOKEN_EXPIRATION_TIME=15s  -p 8080:8080 -e ALLOWED_ORIGINS=*   your_username/chatly
	@echo "Running chatly backend docker image done"