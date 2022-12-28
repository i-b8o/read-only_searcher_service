APP_BIN = app/build/app

lint:
	golangci-lint run

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/main.go

clean:
	rm -rf ./app/build || true


git:
	git add .
	git commit -a -m "$m"
	git push -u origin main

mod:
	cd app;go mod tidy

migrate-up:
	migrate -path ./migrations -database 'postgres://reader:$p@0.0.0.0:5436/reader?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://reader:$p@0.0.0.0:5436/reader?sslmode=disable' down
update_contracts:
	go get -u github.com/i-b8o/read-only_contracts@$m
test:
	go test -p 1 ./...