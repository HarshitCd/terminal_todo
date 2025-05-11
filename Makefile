build:
	@go mod tidy
	@go build -o target/todo main.go

run: build
	@./target/todo

clean:
	rm -rf ./target

