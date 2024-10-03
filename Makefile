run:
	go run main.go --p 7 --l 1 --dir courses --courses all

build:
	goreleaser release --snapshot --clean 

