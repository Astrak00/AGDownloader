package main

import (
	courseParser "AGDownloader/functions"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

func processCourse(wg *sync.WaitGroup, token string, course courseParser.Course) {
	defer wg.Done()
	fmt.Printf("Course: %s\n", course.Name)
	modules, err := courseParser.GetCourseContent(token, course.ID)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	if err := courseParser.SaveFiles(token, course.Name, modules, "cursosDescargados"); err != nil {
		color.Red("Error: %s", err)
	}
}

func main() {
	color.Blue("Download UC3M Aula Global files from Command Line using 'aulaglobalmovil' Security key")

	token_ptr := flag.String("token", "00", "aulaglobalmobile token used to authenticate")
	flag.Parse()
	token := *token_ptr

	userID, _, err := courseParser.GetUserInfo(token)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	fmt.Println("Obtaining courses...")

	courses, err := courseParser.ParseXmlCourses(token, userID)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	startTime := time.Now()
	wg := sync.WaitGroup{}
	for _, course := range courses {
		wg.Add(1)
		go processCourse(&wg, token, course)
	}
	wg.Wait()
	elapsedTime := time.Since(startTime)
	color.Green("All files downloaded in %s", elapsedTime)
}
