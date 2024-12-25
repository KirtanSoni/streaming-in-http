package main

import (
	"embed"
	"io"
	"log"
	"net/http"
)

//go:embed index.html
var content embed.FS
func proxyHandler(w http.ResponseWriter, r *http.Request) {
    // Connect to generator service
    resp, err := http.Get("http://localhost:8081/generate")
    if err != nil {
        http.Error(w, "Generator service error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Set headers for streaming
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Transfer-Encoding", "chunked")

    // Get flusher for real-time forwarding
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }

    // Create buffer for copying
    buf := make([]byte, 1024)

    // Forward chunks in real-time
    for {
        n, err := resp.Body.Read(buf)
        if n > 0 {
            _, writeErr := w.Write(buf[:n])
            if writeErr != nil {
                log.Printf("Write error: %v", writeErr)
                return
            }
            flusher.Flush()
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Read error: %v", err)
            return
        }
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data, err := content.ReadFile("index.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "text/html")
        w.Write(data)
    })

    http.HandleFunc("/stream", proxyHandler)

    log.Printf("Server starting at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}