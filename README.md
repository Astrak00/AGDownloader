# AGDownloader

AGDownloader is a simple and easy-to-use downloader for AulaGlobal (the use of moodle by the University Carlos III of Madrid). It is a command-line tool that allows you to download all the files from all the courses a user has access to.

You can also indicate the courseID or names of the courses you want to download, and the program will only download the files from those courses.

## Usage
To download the programm, go to the [releases page](https://github.com/Astrak00/AGDownloader/releases/latest) and download the latest version for your operating system. You can also [build](#build-from-source) the program from source by following the instructions below. 
> [!NOTE] Key
> Darwin = MacOS
> ARM = M1/M2/M3 and Snapdragon

To use AGDownloader, execute the following command:

```bash
./AGDownloader 
```
You can specify some parameters to customize the download process but if you don't the program will ask you for them (except for the parameters that have a default value).

```bash
./AGDownloader -h

Usage of ./AGDownloader:
      --courses strings   Ids or names of the courses to be downloaded, enclosed in ", separated by spaces. 
                          "all" downloads all courses
      --dir string        Directory where you want to save the files
      --l int             Choose your language: 1: Español, 2:English
      --p int             Cores to be used while downloading (default 4)
      --token string      Aula Global user security token 'aulaglobalmovil'
```

### Obtaining the token

To obtain the token, you must log in to AulaGlobal and go to the preferences panel. There, you will find the token under the "Security keys" section. Copy the token and paste it into the program when prompted.

![Retrieving token](assets/instructions-token.gif)

### Example

This is an example of a full command:

```bash
./AGDownloader --l 1 --token aaaa1111bbbb2222cccc3333dddd4444 --dir AulaGlobal-Copy --p 5 --courses "Inteligencia Distribuidos 123445"
```                           

This program will run in Spanish, with the secret token aaaa1111bbbb2222cccc3333dddd4444, in the folder AulaGlobal-Copy inside the folder where the program is being run, using 5 cores, and downloading the courses that contain the words "Inteligencia" or "Distributed" in their name or the course with the ID 123445.

#### Language

You can choose the language in which the program will run by using the `--l` parameter. The possible values are:
- 1: Spanish
- 2: English

#### Courses

You can specify the courses you want to download by using the `--courses` parameter. You can specify the courses by their ID or by their name. If you want to download all the courses, you can use the keyword "all".

```
./AGDownload --courses "Ingeniería Inteligencia"
```
This parameter will download all the courses that contain the words "Ingeniería" or "Inteligencia" in their name.

#### Directory

You can specify the directory where you want to save the files by using the `--dir` parameter. You must specify the path to the directory where you want to save the files. If you want to download the files in the same directory where the program is being run, you can put a dot.

```
./AGDownload --dir .
```

#### Cores

You can specify the number of cores you want to use while downloading the files by using the `--p` parameter. The default value is 4.

```
./AGDownload --p 8
```





## Build from source

To build the program from source, you will need to have Go installed on your computer. You can download it from the [official website](https://golang.org/). Once you have installed Go, you can clone the repository and build the program by running the following commands:

```bash
git clone git@github.com:Astrak00/AGDownloader.git
go build
```
This will create an executable file called AGDownloader that you can run.


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## Inspiration
This project was inspired by the need to download all the files from all the courses in AulaGlobal. This is a tedious task that can take a lot of time, so I decided to create a tool that would allow me to do this in a simple and easy way.

Previously, I had created a similar tool in Python, but I decided to create a new one in Go because I wanted to learn more about this language and the faster execution time it offers was very appealing to me.

The idea of using the mobile token to authenticate the user was inspired by the project created by [Josersanvil](https://github.com/Josersanvil/AulaGlobal-CoursesFiles). I decided to use this method because it is more secure than using the user's password and because this year, the university removed the login with user and password.

## License

This project is licensed under the MIT License

