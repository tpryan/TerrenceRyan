package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/go-github/v33/github"
	"github.com/mmcdole/gofeed"
)

var cache *Cache
var cacheEnabled = true
var verbose = true
var logger *log.Logger
var parser = gofeed.NewParser()

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

		f, err := parser.ParseURL("https://speakerdeck.com/tpryan.atom")

		content := Content{}

		if err := content.LoadFeed(f); err != nil {
			writeError(w, fmt.Errorf("error converting to content: %s", err))
			return
		}

		presoJSON, err = content.JSONandCache(cache, "presos")
		if err != nil {
			writeError(w, err)
		}

	}

	writeResponse(w, http.StatusOK, presoJSON)
	return
}

func handleRepos(w http.ResponseWriter, r *http.Request) {

	repoJSON, err := cache.Get("github")

	if err != nil {

		client := github.NewClient(nil)
		opt := &github.RepositoryListOptions{Type: "public", Sort: "pushed"}
		f, _, err := client.Repositories.List(context.Background(), "tpryan", opt)

		if err != nil {
			writeError(w, fmt.Errorf("error retrieving repos: %s", err))
			return
		}

		content := Content{}

		if err := content.LoadGithub(f); err != nil {
			writeError(w, fmt.Errorf("error converting to content: %s", err))
			return
		}

		repoJSON, err = content.JSONandCache(cache, "github")
		if err != nil {
			writeError(w, err)
		}
	}

	writeResponse(w, http.StatusOK, repoJSON)
	return
}

func handlePosts(w http.ResponseWriter, r *http.Request) {

	postJSON, err := cache.Get("blog")

	if err != nil {
		f, err := parser.ParseURL("https://tpryan.blog/feed/")

		content := Content{}

		if err := content.LoadFeed(f); err != nil {
			writeError(w, fmt.Errorf("error converting to content: %s", err))
			return
		}

		postJSON, err = content.JSONandCache(cache, "blog")
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

	c.log("Successfully retrieved content from cache as key: %s", key)

	return s, nil
}

// Content is a collection of nuggets with the
type Content struct {
	Nuggets []Nugget `json:"nuggets"`
	Cached  bool     `json:"cached"`
}

// Nugget is a what goes on the front page of the website.
type Nugget struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Timestamp   time.Time `json:"timestamp"`
}

// JSON Returns the given content struct as a JSON string
func (c Content) JSON() (string, error) {

	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(b), nil
}

// JSONandCache Caches and Returns the given content struct as a JSON string
func (c *Content) JSONandCache(cache *Cache, key string) (string, error) {
	outJSON, err := c.JSON()
	if err != nil {
		return "", fmt.Errorf("error converting content to json: %s", err)
	}

	c.Cached = true
	saveJSON, err := c.JSON()
	if err != nil {
		return "", fmt.Errorf("error converting content to savable json: %s", err)
	}

	if err := cache.Save(key, saveJSON); err != nil {
		cache.logger.Printf("error saving to cache: %s", err)
	}

	return outJSON, nil
}

// LoadGithub takes a github response and loads it into content.
func (c *Content) LoadGithub(g []*github.Repository) error {
	nuggets := []Nugget{}

	for _, v := range g {

		n := Nugget{}
		n.Title = v.GetName()
		n.Timestamp = v.GetUpdatedAt().Time
		n.URL = v.GetHTMLURL()
		n.Description = v.GetDescription()
		nuggets = append(nuggets, n)

	}
	c.Nuggets = nuggets

	return nil
}

// LoadFeed takes a feed response and loads it into content.
func (c *Content) LoadFeed(f *gofeed.Feed) error {

	nuggets := []Nugget{}

	for _, v := range f.Items {

		var err error

		n := Nugget{}
		n.Title = v.Title
		n.URL = v.Link

		if f.FeedType == "atom" {
			n.Description = v.Content
			n.Timestamp, err = time.Parse("2006-1-2T15:04:05-07:00", v.Published)
			if err != nil {
				return err
			}
		}

		if f.FeedType == "rss" {
			n.Description = v.Description
			n.Timestamp, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", v.Published)
			if err != nil {
				return err
			}
		}

		nuggets = append(nuggets, n)
	}

	c.Nuggets = nuggets
	return nil
}
