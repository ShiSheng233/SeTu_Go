package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RandomPost struct {
	FileUrl string `json:"file_url"`
}

func main() {
	http.HandleFunc("/api/setu", func(w http.ResponseWriter, r *http.Request) {
		tags := r.URL.Query().Get("tags")

		url := "https://danbooru.donmai.us/posts/random.json"
		if tags != "" {
			url += "?tags=" + tags
		}

		client := http.Client{}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error fetching random post", http.StatusInternalServerError)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				http.Error(w, "Error closing response body", http.StatusInternalServerError)
				return
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		var post RandomPost
		err = json.Unmarshal(body, &post)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
			return
		}

		if tags != "" {
			log.Println("[INFO]\nTags:", tags, "\nUrl:", post.FileUrl, "\nUserAgent:", r.UserAgent(), "\nIP:", r.RemoteAddr)
		} else {
			log.Println("[INFO]\nTags: no", "\nUrl:", post.FileUrl, "\nUserAgent:", r.UserAgent(), "\nIP:", r.RemoteAddr)
		}

		http.Redirect(w, r, post.FileUrl, http.StatusFound)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("[INFO]\nStarted to listen on http://localhost:" + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		return
	}
}
