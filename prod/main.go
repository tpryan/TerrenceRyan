package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
)

var cache *Cache
var cacheEnabled = true
var verbose = true
var logger *log.Logger

func main() {

	logger = log.New(os.Stderr, "MAIN : ", log.Ldate|log.Ltime|log.Lmsgprefix)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Printf("Defaulting to port %s", port)
	}

	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")

	var err error
	cache, err = NewCache(redisHost, redisPort, cacheEnabled, verbose)
	if err != nil {
		logger.Fatal(err)
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

	logger.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handlePresos(w http.ResponseWriter, r *http.Request) {

	presoJSON, err := cache.Get("presos")

	if err != nil {
		presoXML, err := get("https://speakerdeck.com/tpryan.atom")
		if err != nil {
			writeError(w, err)
		}

		var f Feed
		if err := xml.Unmarshal(presoXML, &f); err != nil {
			writeError(w, err)
		}

		presoJSON, err = f.JSON()
		err = cache.Save("presos", presoJSON)
		if err != nil {
			logger.Printf("error saving to cache: %s", err)
		}

	}

	writeResponse(w, http.StatusOK, presoJSON)
	return
}

func handleRepos(w http.ResponseWriter, r *http.Request) {

	repoJSON, err := cache.Get("github")

	if err != nil {
		temp, err := get("https://api.github.com/users/tpryan/repos?sort=pushed")
		if err != nil {
			writeError(w, err)
		}
		repoJSON = string(temp)

		err = cache.Save("github", repoJSON)
		if err != nil {
			logger.Printf("error saving to cache: %s", err)
		}
	}

	writeResponse(w, http.StatusOK, repoJSON)
	return
}

func handlePosts(w http.ResponseWriter, r *http.Request) {

	postJSON, err := cache.Get("blog")

	if err != nil {
		postXML, err := get("https://tpryan.blog/feed/")
		if err != nil {
			writeError(w, err)
		}

		var f RssFeed
		if err := xml.Unmarshal(postXML, &f); err != nil {
			writeError(w, err)
		}
		postJSON, err = f.JSON()

		err = cache.Save("blog", postJSON)
		if err != nil {
			logger.Printf("error saving to cache: %s", err)
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write([]byte(msg))
	return
}

func writeError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("{\"err\":\"%s\"}", err)
	logger.Printf("error serving content %s: ", err)
	writeResponse(w, http.StatusInternalServerError, msg)
}

// Feed is a representation of a Atom feed
type Feed struct {
	Entries []Entry `json:"entries" xml:"entry"`
}

// JSON Returns the given Feed struct as a JSON string
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

// Entry are individual items in a Feed
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

// RssFeed is a representation of a RSS feed
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

// Item are individual entries in an RSSFeed
type Item struct {
	Title     string `json:"title" xml:"title"`
	Content   string `json:"content" xml:"description"`
	Published string `json:"published" xml:"pubDate"`
	Link      string `json:"link" xml:"link"`
}

// RedisPool is an interface that allows us to swap in an mock for testing cache
// code.
type RedisPool interface {
	Get() redis.Conn
}

// ErrCacheMiss error indicates that an item is not in the cache
var ErrCacheMiss = fmt.Errorf("item is not in cache")

// Cache abstracts all of the operations of caching for the application
type Cache struct {
	redisPool RedisPool
	enabled   bool
	verbose   bool
	logger    *log.Logger
}

// NewCache returns an initialized cache ready to go.
func NewCache(redisHost, redisPort string, enabled, verbose bool) (*Cache, error) {
	c := &Cache{}
	c.logger = log.New(os.Stderr, "CACHE : ", log.Ldate|log.Ltime|log.Lmsgprefix)
	c.redisPool = c.InitPool(redisHost, redisPort)
	c.enabled = enabled
	c.verbose = verbose
	return c, nil
}

// InitPool starts the cache off
func (c Cache) InitPool(redisHost, redisPort string) RedisPool {
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	c.log("Initialized Redis at %s", redisAddr)
	const maxConnections = 10

	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}, maxConnections)

	return pool
}

func (c *Cache) log(msg string, args ...interface{}) {
	if len(args) > 0 {
		c.logger.Printf(msg, args...)
		return
	}
	c.logger.Printf(msg)
	return
}

// Save saves a list of all of the games a player is in.
func (c *Cache) Save(key, content string) error {
	if !c.enabled {
		return nil
	}

	conn := c.redisPool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", key, content); err != nil {
		return err
	}
	c.log("Successfully saved content to cache as key: %s", key)
	return nil
}

// Get retrieves content from the cache
func (c *Cache) Get(key string) (string, error) {
	if !c.enabled {
		return "", ErrCacheMiss
	}

	conn := c.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", ErrCacheMiss
	} else if err != nil {
		return "", err
	}

	c.log("Successfully retrieved content from cache as key :%s", key)

	return s, nil
}
