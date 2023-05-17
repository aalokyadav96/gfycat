package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
	"path/filepath"
	"log"
)

func postFile(filename string, targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    // this step is very important
    fileWriter, err := bodyWriter.CreateFormFile("imgfile", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    // open file handle
    fh, err := os.Open(filename)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }
    defer fh.Close()

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return nil
}

// sample usage
func main() {
    target_url := "http://localhost:4000/headless"
	
		files := TraverseDir("f:/lets_get_it/Lets Go/Go2C/tt/*") 
		for _,file := range files {
			if isDirectory(file) == false && isFileValid(file){
			filename := file
			postFile(filename, target_url)
			}
		}
}


func TraverseDir(path string) []string {
    files, err := filepath.Glob(path)
    if err != nil {
        log.Fatal(err)
    }
	return files
}
	

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func isFileValid(path string) bool {
	fileExtension := filepath.Ext(path)
	switch fileExtension {
		case ".mp4" :
			fmt.Println("case4")
			return true
		case ".webm" :
			fmt.Println("casem ")
			return true
	}
	return false		
}