package views

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var allowedTypes = [...]string{
	"png",
	"gif",
	"jpg",
	"mp3",
	"ogg",
	"opus",
	"mp4",
	"webm",
}

type file struct {
	Url string `json:"url"`
}

type uploadResponse struct {
	Success bool     `json:"success"`
	Files   []file   `json:"files"`
	Errors  []string `json:"errors"`
}

func HandleUpload(r *http.Request, w http.ResponseWriter) (error, string) {
	if r.Method == "POST" {
		return HandleUploadAPI(r, w)
	}
	return nil, "Not yet implemented!"
}

func HandleUploadAPI(r *http.Request, w http.ResponseWriter) (error, string) {
	if r.Method != "POST" {
		response := uploadResponse{
			false,
			[]file{},
			[]string{"GET not supported. Use POST."},
		}
		b, _ := json.Marshal(response)
		return nil, string(b)
	}
	upFile, header, err := r.FormFile("file")
	if err != nil || upFile == nil {
		upFile, header, err = r.FormFile("files[]") // This is legacy!
	}
	if err != nil || upFile == nil {
		response := uploadResponse{
			false,
			[]file{},
			[]string{"No file sent. Please send a file as \"file\"."},
		}
		b, _ := json.Marshal(response)
		return nil, string(b)
	}
	defer upFile.Close()

	localname := header.Filename
	dotSplit := strings.Split(localname, ".")
	extension := strings.ToLower(dotSplit[len(dotSplit)-1])

	// Check if we allow the extension
	extensionAllowed := false
	for _, b := range allowedTypes {
		if extension == b {
			extensionAllowed = true
		}
	}

	if !extensionAllowed {
		response := uploadResponse{
			false,
			[]file{},
			[]string{"Extension \"" + extension + "\" not supported"},
		}
		b, _ := json.Marshal(response)
		return nil, string(b)
	}

	filename := genFilename()
	if filename == "" {
		response := uploadResponse{
			false,
			[]file{},
			[]string{"We failed to create a filename"},
		}
		b, _ := json.Marshal(response)
		return nil, string(b)
	}

	out, err := os.Create("uploads/" + filename + "." + extension)
	if err != nil {
		response := uploadResponse{
			false,
			[]file{},
			[]string{"Failed to create a file"},
		}
		b, _ := json.Marshal(response)
		return nil, string(b)
	}
	defer out.Close()

	_, err = io.Copy(out, upFile)
	if err != nil {
		response := uploadResponse{
			false,
			[]file{},
			[]string{"Failed to save to file"},
		}
		b, _ := json.Marshal(response)
		return nil, string(b)
	}

	response := uploadResponse{
		true,
		[]file{file{filename}},
		[]string{},
	}
	b, _ := json.Marshal(response)
	return nil, string(b)
}

func genFilename() string {
	iter := 0
	exists := true
	filename := ""
	for exists {
		if iter > 25 {
			return ""
		}
		filename = rngStr(3 + int(iter/5)) // Add one letter per 5 retries
		e := false
		for _, extension := range allowedTypes {
			if fileExists("uploads/" + filename + "." + extension) {
				e = true
			}
		}
		if !e {
			exists = false
		}
		iter++
	}
	return filename
}

func rngStr(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
