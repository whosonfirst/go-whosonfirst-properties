cli:
	go build -mod vendor -o bin/report-properties cmd/report-properties/main.go
	go build -mod vendor -o bin/index-properties cmd/index-properties/main.go

docker:
	docker build -t whosonfirst-properties-indexing .	

