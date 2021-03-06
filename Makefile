NAME=ntc-gnats
VERSION=0.0.1


.PHONY: req
req:
	@go run req.go

.PHONY: res
res:
	@go run res.go

.PHONY: sub
sub:
	@go run sub.go

.PHONY: pub
pub:
	@go run pub.go

.PHONY: worker
worker:
	@go run worker.go

.PHONY: main
main:
	@go run main.go

.PHONY: build
build:
	@echo "Build project..."
	@go build -o $(NAME)

.PHONY: run
run:
	@echo "Run project mode development..."
	@./$(NAME) -e development

.PHONY: run-test
run-test:
	@echo "Run project mode test..."
	@nohup ./$(NAME) -e test >/dev/null 2>&1 &

.PHONY: run-stag
run-stag:
	@echo "Run project mode staging..."
	@nohup ./$(NAME) -e staging >/dev/null 2>&1 &

.PHONY: run-prod
run-prod:
	@echo "Run project mode production..."
	@nohup ./$(NAME) -e production >/dev/null 2>&1 &

.PHONY: clean
clean:
	@echo "Clean project..."
	@rm -f $(NAME)

.PHONY: deps
deps:
	@./deps.sh

.PHONY: test
test:
	@echo "Run test..."
	@go test -v ./test/*
