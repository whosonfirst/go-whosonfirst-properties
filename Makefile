cli:
	go build -mod vendor -o bin/report cmd/report/main.go
	go build -mod vendor -o bin/index cmd/index/main.go
