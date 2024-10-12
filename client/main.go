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

type syncWriter struct {
	m      sync.Mutex
	writer io.Writer
}

func (w *syncWriter) Write(b []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.writer.Write(b)
}

func main() {

	engine := html.New("./views", ".html")

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
		go worker()
		return nil
	})

	log.Fatal(app.Listen(":3000"))

}

func worker() {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}
	defer file.Close()
	wr := &syncWriter{sync.Mutex{}, file}
	var wg sync.WaitGroup
	read, write := io.Pipe()

	wg.Add(1)
	// чтение из response
	go func() {
		defer write.Close()
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to get request: %v", err)
		}
		defer resp.Body.Close()
		io.Copy(write, resp.Body)

	}()
	// запись в файл
	go func() {
		defer wg.Done()
		io.Copy(wr, read)
	}()
	//wait group
	wg.Wait()
}
