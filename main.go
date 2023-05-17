package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
var uploadPath = "./uploads"
//var tmpl = template.Must(template.ParseGlob("templates/*.html"))
var tmpl = template.Must(template.ParseGlob("_includes/*.html"))


func main() {
	http.HandleFunc("/", uploadFileHandler())
	http.HandleFunc("/headless", headlessUpload())
	http.HandleFunc("/delete", deleteUploads)

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	thumfs := http.FileServer(http.Dir("./thumbs"))
	http.Handle("/img/", http.StripPrefix("/img", thumfs))

	log.Print("Server started on :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl.ExecuteTemplate(w, "head.html",nil)
			tmpl.ExecuteTemplate(w,"index.html",nil)
			return
		}
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}
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
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "video/mp4", "video/webm":
		case "image/gif", "image/png", "image/jpg", "image/jpeg":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := randToken(12)
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
	if fileEndings[0] == ".f4v" {fileEndings[0] = ".mp4"}
		newFileName := fileName + fileEndings[0]

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
		

		pics = append(pics,"/files/"+newFileName)
		}
		
	fmt.Println(pics)
	tmpl.ExecuteTemplate(w,"returnUploads.html", pics)

//		w.Write([]byte(fmt.Sprintf("SUCCESS - use /files/%v to access the file", newFileName)))
//			t, _ := template.ParseFiles("templates/returnUploads.html")
//			t, _ := template.ParseFiles("res.html")
//			t.Execute(w, "files/"+newFileName)
//			return
	})
}

func headlessUpload() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.Method != "POST" {
//			w.Write([]byte(fmt.Sprintf("only post is allowed")))
//			return
//		}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}
		var fileEndings string
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
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
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
		case "image/png":
			fileEndings = ".png"
			break
		case "image/jpg", "image/jpeg":
			fileEndings = ".jpg"
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := randToken(12)
//		fileEndings, err := mime.ExtensionsByType(detectedFileType)

		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		fmt.Println(fileEndings)
//		newFileName := fileName + fileEndings[0]
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
		

		pics = append(pics,"/files/"+newFileName)
		}
		
	fmt.Println(pics)
//	tmpl.ExecuteTemplate(w,"returnUploads.html", )

		w.Write([]byte(fmt.Sprintf("SUCCESS : ", pics)))
//			t, _ := template.ParseFiles("templates/returnUploads.html")
//			t, _ := template.ParseFiles("res.html")
//			t.Execute(w, "files/"+newFileName)
//			return
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

func deleteUploads(w http.ResponseWriter, r *http.Request) {
   err := os.Remove(uploadPath)
   if err != nil {
      fmt.Println(err)
   } else {
      fmt.Println("Directory", uploadPath, "removed successfully")
   }
	w.Write([]byte(fmt.Sprintf("deleted all")))
}
