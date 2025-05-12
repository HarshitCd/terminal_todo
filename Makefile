ENV_PATH?=.env

build:
	@go mod tidy
	@go build -ldflags="-X 'main.envPath=$(ENV_PATH)'" -o target/todo main.go

run: build
	@./target/todo $(ARGS)

clean:
	rm -rf ./target

