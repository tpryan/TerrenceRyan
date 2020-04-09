package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var postJSON string
var presoJSON string
var repoJSON string

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/presos", handlePresos)
	http.HandleFunc("/repos", handleRepos)
	http.HandleFunc("/posts", handlePosts)
	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.HandleFunc("/resume", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handlePresos(w http.ResponseWriter, r *http.Request) {
	if presoJSON == "" {
		presoXML, err := get("https://speakerdeck.com/tpryan.atom")
		if err != nil {
			writeError(w, err)
		}

		var f Feed
		if err := xml.Unmarshal(presoXML, &f); err != nil {
			writeError(w, err)
		}

		presoJSON, err = f.JSON()
		if err != nil {
			writeError(w, err)
		}

	}

	writeResponse(w, http.StatusOK, presoJSON)
	return
}

func handleRepos(w http.ResponseWriter, r *http.Request) {
	if repoJSON == "" {
		temp, err := get("https://api.github.com/users/tpryan/repos?sort=pushed")
		if err != nil {
			writeError(w, err)
		}
		repoJSON = string(temp)
	}

	writeResponse(w, http.StatusOK, repoJSON)
	return
}

func handlePosts(w http.ResponseWriter, r *http.Request) {

	if postJSON == "" {
		postXML, err := get("https://tpryan.blog/feed/")
		if err != nil {
			writeError(w, err)
		}

		var f RssFeed
		if err := xml.Unmarshal(postXML, &f); err != nil {
			writeError(w, err)
		}
		postJSON, err = f.JSON()
		if err != nil {
			writeError(w, err)
		}

	}

	writeResponse(w, http.StatusOK, postJSON)
	return
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "ok")
	return
}

func get(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func writeResponse(w http.ResponseWriter, code int, msg string) {

	if code != http.StatusOK {
		log.Printf(msg)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write([]byte(msg))

	return
}

func writeError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("{\"err\":\"%s\"}", err)
	writeResponse(w, http.StatusInternalServerError, msg)
}

type Feed struct {
	Entries []Entry `json:"entries" xml:"entry"`
}

// JSON Returns the given DFResponse struct as a JSON string
func (f Feed) JSON() (string, error) {

	items := []Item{}

	for _, v := range f.Entries {
		item := Item{}
		item.Title = v.Title
		item.Published = v.Published
		item.Link = v.Link.HREF
		item.Content = v.Content
		items = append(items, item)
	}

	b, err := json.Marshal(items)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(b), nil
}

type Entry struct {
	Title     string `json:"title" xml:"title"`
	Content   string `json:"content" xml:"content"`
	Published string `json:"published" xml:"published"`
	Link      struct {
		HREF string `json:"href" xml:"href,attr"`
	} `json:"link" xml:"link"`
	Author struct {
		Name string `json:"name" xml:"name"`
	} `json:"author" xml:"author"`
}

type RssFeed struct {
	Channel struct {
		Items []Item `json:"items" xml:"item"`
	} `json:"channel" xml:"channel"`
}

// JSON Returns the given DFResponse struct as a JSON string
func (f RssFeed) JSON() (string, error) {
	items := f.Channel.Items

	b, err := json.Marshal(items)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(b), nil
}

type Item struct {
	Title     string `json:"title" xml:"title"`
	Content   string `json:"content" xml:"description"`
	Published string `json:"published" xml:"pubDate"`
	Link      string `json:"link" xml:"link"`
}
