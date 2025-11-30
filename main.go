package main

import (
	"log"
	"net/http"
	"os"
)

func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "app/index.html")
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Generate directory listing for /app/assets
	entries, err := os.ReadDir("app/assets")
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("<pre>\n"))
	for _, entry := range entries {
		w.Write([]byte(`<a href="` + entry.Name() + `">` + entry.Name() + `</a>\n`))
	}
	w.Write([]byte("</pre>"))
}

func main() {
	mux := http.NewServeMux() // router
	mux.HandleFunc("/healthz", readinessCheckHandler)
	mux.HandleFunc("/app", appHandler)
	mux.HandleFunc("/app/assets", assetsHandler)
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := &http.Server{
		Addr:    ":8080", // Bind to localhost:8080
		Handler: mux,     // Use our empty mux
	}

	log.Println("server running on :8080")
	log.Fatal(server.ListenAndServe())

}
