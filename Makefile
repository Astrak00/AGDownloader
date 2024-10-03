run:
	go run main.go --p 7 --l 1 --dir courses
build:
	goreleaser release --snapshot --clean 

