# AGDownloader

AGDownloader is a simple and easy-to-use downloader for AulaGlobal (the use of moodle by the University Carlos III of Madrid). It is a command-line tool that allows you to download all the files from all the courses a user has access to.

## Installation

To install AGDownloader, you will need to compile it from source using go. You can do this by running the following command:

```bash
go build -o agdownloader main.go
```

## Usage

To use AGDownloader, you will need to run the following command:

```bash
./agdownloader -l <language:1\|2> -t <token> -d <directory>
```

Where `<token>` is the token you can get from the preferences panel and `<directory>` is the directory where you want to download the files.


## License

This project is licensed under the MIT License

