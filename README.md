# AGDownloader

AGDownloader is a simple and easy-to-use downloader for AulaGlobal (the use of moodle by the University Carlos III of Madrid). It is a command-line tool that allows you to download all the files from all the courses a user has access to.

## Installation

To install AGDownloader, you will need to compile it from source using go. You can do this by running the following command:

```bash
go build
```

## Usage

To use AGDownloader, you will need to run the following command:

```bash
./AGDownloader -l <language=1|2> -t <token> -d <directory> -c <num_cores>
```

Where `<token>` is the token you can get from the preferences panel and `<directory>` is the directory where you want to download the files.
`<num_cores>` is the number of cores you want to use to download the files. If you don't specify this parameter, the program will use all the cores available.


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Inspiration
This project was inspired by the need to download all the files from all the courses in AulaGlobal. This is a tedious task that can take a lot of time, so I decided to create a tool that would allow me to do this in a simple and easy way.

Previously, I had created a similar tool in Python, but I decided to create a new one in Go because I wanted to learn more about this language and the faster execution time it offers was very appealing to me.

The idea of using the mobile token to authenticate the user was inspired by the project created by [Josersanvil](github.com/Josersanvil/AulaGlobal-CoursesFiles). I decided to use this method because it is more secure than using the user's password and because this year, the university removed the login with user and password.

## License

This project is licensed under the MIT License

