package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func generateHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Transfer-Encoding", "chunked")
    
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }

    // Generate small chunks rapidly to test for missing data
    for i := 1; i <= 100; i++ {
        chunk := fmt.Sprintf("Chunk %d [%s]\n", i, time.Now().Format("15:04:05.000"))
        _, err := w.Write([]byte(chunk))
        if err != nil {
            log.Printf("Write error: %v", err)
            return
        }
        flusher.Flush()
        time.Sleep(time.Millisecond * 100) // 10 chunks per second
    }
}
func main() {
    http.HandleFunc("/generate", generateHandler)
    log.Printf("Generator service starting at http://localhost:8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}
