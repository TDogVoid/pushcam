package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/radovskyb/watcher"
)

func main() {
	//set token
	Token := "" // Access Token under Account Settings
	title := "Motion Detected"
	link := ""                     //link sent out to watch camera so you can easily click on it
	push(Token, title, link, true) // bool value is if you want it too send a image in second message
}

//watchFolder directory until it sees a new image then pushes out image and stops
// This may not work if video is not being recorded?
func watchFolder(dir string, pb *pushbullet.Pushbullet, title string) {

	cameraFolder := ""
	if len(os.Args) > 1 {
		cameraFolder = os.Args[1]
		dir += "/" + cameraFolder
	}

	w := watcher.New()
	w.IgnoreHiddenFiles(true)
	w.Ignore()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if filepath.Ext(event.Path) == ".jpg" {
					sendImage(pb, title, event.Path)
					w.Close()
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(dir); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

//push out notification
func push(token, title, link string, sendImage bool) {

	pb := pushbullet.New(token)

	sendLink(pb, title, link) // send link out
	if sendImage {
		watchFolder("/data/output", pb, title) // watch folder where photos are output and push file out once detected
	}

}

//send link out
func sendLink(pb *pushbullet.Pushbullet, title string, link string) {
	l := requests.NewLink()
	l.Title = title
	l.Url = link
	if _, err := pb.PostPushesLink(l); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
}

//sendImage sends the image to pushbullet and notifys
func sendImage(pb *pushbullet.Pushbullet, title string, filepath string) {
	file, filename := getNewestImage(filepath)
	defer file.Close()
	fileType := "image/jpeg"

	// request permission to upload file
	res, err := pb.PostUploadRequest(filename, fileType)
	if err != nil {
		fmt.Println(err)
	}

	//upload file
	pushbullet.Upload(pb.Client(), res, file)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.FileUrl)

	//Push Notification
	f := requests.NewFile()
	f.Title = filename
	f.FileUrl = res.FileUrl
	f.FileName = filename

	f.FileType = "image/jpeg"

	if _, err := pb.PostPushesFile(f); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
}

//getNewestImage gets image file and filename
func getNewestImage(path string) (*os.File, string) {

	openFile, err := os.Open(path)
	basename := filepath.Base(path)
	filename := strings.TrimSpace(strings.TrimSuffix(basename, filepath.Ext(basename)))

	if err != nil {
		panic(err)
	}
	return openFile, filename
}
