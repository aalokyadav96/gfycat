package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const PORT = "4000"
const maxUploadSize = 10 * 1024 * 1024 // 8 mb
var uploadPath = "./uploads"
var posterpath =  "./poster"
var tmpl = template.Must(template.ParseGlob("_includes/*.html"))

type Res struct {
	Video string
	Poster string
}

func main() {
	http.HandleFunc("/", uploadFileHandler())
	http.HandleFunc("/v/", viewPost())
	http.HandleFunc("/getall", GetAll())
	http.HandleFunc("/del/", DeleteVid())

	vidfs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", vidfs))

	posterfs := http.FileServer(http.Dir(posterpath))
	http.Handle("/poster/", http.StripPrefix("/poster", posterfs))

	log.Print("Server started on localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

type PostHead struct {
	PostURL string
	Title string
	WebsiteName string
}

//"https://gifs-ba0f.onrender.com"
//http://localhost:4000

func viewPost()  http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			title := r.URL.Path[len("/v/"):]
			fmt.Println("path", r.URL.Path)
			res := Res {Video: "/files/"+ title, Poster: "/poster/"+ title}
			jg := PostHead{PostURL: title, Title: "Test GIF", WebsiteName: "https://gifs-ba0f.onrender.com"}
			tmpl.ExecuteTemplate(w, "head.html", jg)
			tmpl.ExecuteTemplate(w, "viewpost.html", res)
			return
		}
	})
}

func GetAll()  http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("path", r.URL.Path)
			entries, err := os.ReadDir("./uploads")
			if err != nil {
				log.Fatal(err)
			}
		var vids []string
			for _, e := range entries {
				vids = append(vids, e.Name())
			}
			tmpl.ExecuteTemplate(w, "head.html", nil)
			tmpl.ExecuteTemplate(w, "allfiles.html", vids)
			return
		}
	})
}

func DeleteVid()  http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			title := r.URL.Path[len("/del/"):]
			err := os.Remove("./uploads/"+title+".mp4")  // remove a single file
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("path", r.URL.Path)
			w.Write([]byte("done"))
			return
		}
	})
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("path", r.URL.Path)
			tmpl.ExecuteTemplate(w, "head.html", nil)
			tmpl.ExecuteTemplate(w, "index.html", nil)
			return
		}
			fmt.Println("path", r.URL.Path)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}
		var fileEndings string
		var fileName string
		var pics []string
		files := r.MultipartForm.File["imgfile"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer file.Close()

			// Get and print out file size
			fileSize := fileHeader.Size
			fmt.Printf("File size (bytes): %v\n", fileSize)
			// validate file size
			if fileSize > maxUploadSize {
				renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
				return
			}
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				renderError(w, "INVALID_FILE"+http.DetectContentType(fileBytes), http.StatusBadRequest)
				return
			}

			// check file type, detectcontenttype only needs the first 512 bytes
			detectedFileType := http.DetectContentType(fileBytes)
			switch detectedFileType {
			case "video/mp4":
				fileEndings = ".mp4"
				break
			case "video/webm":
				fileEndings = ".webm"
				break
			case "image/gif":
				fileEndings = ".gif"
				break
			default:
				renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
				return
			}
			fileName = randToken(12)
			//		fileEndings, err := mime.ExtensionsByType(detectedFileType)

			if err != nil {
				renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
				return
			}
			newFileName := fileName + fileEndings

			newPath := filepath.Join(uploadPath, newFileName)
			fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}

			pics = append(pics, "/files/"+newFileName)
		}
		fmt.Println(pics)
	w.Write([]byte("/v/"+fileName))
//		tmpl.ExecuteTemplate(w, "player.html", pics)

	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
