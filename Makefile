run:
	go run main.go --p 7 --l 1 --dir courses_temp

build:
	go build -o AGDownload main.go

release:
	goreleaser release --snapshot --clean 

test:
	go test -v ./...

clean:
	rm -rf dist 
