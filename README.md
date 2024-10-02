# AGDownloader

AGDownloader is a simple and easy-to-use downloader for AulaGlobal (the use of moodle by the University Carlos III of Madrid). It is a command-line tool that allows you to download all the files from all the courses a user has access to.

You can also indicate the courseID or names of the courses you want to download, and the program will only download the files from those courses.

## Usage

To use AGDownloader, execute the following command:

```bash
./AGDownloader 
```
You can specify some parameters to customize the download process but if you don't the program will ask you for them (except for the parameters that have a default value).

```bash
./AGDownloader -h

Usage of ./AGDownloader:
      --courses strings   Ids or names of the courses you want to download enclosed in ", separated by spaces. 
                          "all" downloads all courses
      --dir string        Directory where you want to save the files (default "courses")
      --l int             Choose your language: 1: Español, 2:English
      --p int             Cores to be used while downloading (default 7)
      --token string      Aula Global user security token 'aulaglobalmovil'
```

Here, `<token>` is the token available in the preferences panel, and `<directory>` specifies the location where you want to save the downloaded files. The `<num_cores>` parameter indicates how many cores you wish to allocate for the download. If you do not provide this, the program will utilize all available cores minus one.


### Obtaining the token

To obtain the token, you must log in to AulaGlobal and go to the preferences panel. There, you will find the token under the "Security keys" section. Copy the token and paste it into the program when prompted.

![Retrieving token](assets/instructions-token.gif)


## Build from source

To build the program from source, you will need to have Go installed on your computer. You can download it from the [official website](https://golang.org/). Once you have installed Go, you can clone the repository and build the program by running the following commands:

```bash
git clone git@github.com:Astrak00/AGDownloader.git
go build
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## Inspiration
This project was inspired by the need to download all the files from all the courses in AulaGlobal. This is a tedious task that can take a lot of time, so I decided to create a tool that would allow me to do this in a simple and easy way.

Previously, I had created a similar tool in Python, but I decided to create a new one in Go because I wanted to learn more about this language and the faster execution time it offers was very appealing to me.

The idea of using the mobile token to authenticate the user was inspired by the project created by [Josersanvil](https://github.com/Josersanvil/AulaGlobal-CoursesFiles). I decided to use this method because it is more secure than using the user's password and because this year, the university removed the login with user and password.

## License

This project is licensed under the MIT License

