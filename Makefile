cli:
	go build -mod vendor -o bin/report cmd/report/main.go
	go build -mod vendor -o bin/index-properties cmd/index-properties/main.go

docker:
	docker build -t whosonfirst-properties-indexing .	

