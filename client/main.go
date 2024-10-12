package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

const (
	url      = "http://localhost:8000/api/items"
	filename = "file.txt"
)

type SyncWriter struct {
	m      sync.Mutex
	Writer io.Writer
}

func (w *SyncWriter) Write(b []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.Writer.Write(b)
}

func main() {

	engine := html.New("./public", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Index!",
		})
	})

	app.Get("/message", func(c *fiber.Ctx) error {
		go Worker()
		return nil
	})

	log.Fatal(app.Listen(":3000"))

}

func Worker() {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}
	defer file.Close()
	wr := &SyncWriter{sync.Mutex{}, file}
	read, write := io.Pipe()

	go func() {
		defer write.Close()
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to get request: %v", err)
		}
		defer resp.Body.Close()
		io.Copy(write, resp.Body)

	}()

	io.Copy(wr, read)

}
