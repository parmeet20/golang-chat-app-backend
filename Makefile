# Makefile for Chatly backend

build:
	@echo "Building chatly backend"
	go build -o ./bin/chatly cmd/server/main.go
	@echo "Building chatly backend done"

run-test:
	@echo "Running chatly backend tests"
	go test ./...
	@echo "Running chatly backend tests done"

clean:
	@echo "Cleaning chatly backend"
	rm -rf ./bin
	@echo "Cleaning chatly backend done"

run:
	@echo "Running chatly backend"
	./bin/chatly
	@echo "Running chatly backend done"

build-docker-image:
	@echo "Building chatly backend docker image"
	docker build -t ${DOCKER_USERNAME}/chatly:${GITHUB_SHA} .
	@echo "Building chatly backend docker image done"

docker-run:
	@echo "Running chatly backend docker image"
	docker run -e PORT=8080 \
	-e MONGO_URL="${MONGO_URL}" \
	-e JWT_SECRET="${JWT_SECRET}" \
	-e JWT_TOKEN_EXPIRATION_TIME="15s" \
	-p 8080:8080 \
	-e ALLOWED_ORIGINS="*" \
	${DOCKER_USERNAME}/chatly:${GITHUB_SHA}
	@echo "Running chatly backend docker image done"