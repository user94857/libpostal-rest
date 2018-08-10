package main

import (
        "context"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "os/signal"
        "time"
        "github.com/gorilla/mux"
        expand "github.com/openvenues/gopostal/expand"
        parser "github.com/openvenues/gopostal/parser"
)

func main() {
        host := os.Getenv("LISTEN_HOST")
        if host == "" {
                host = "0.0.0.0"
        }
        port := os.Getenv("LISTEN_PORT")
        if port == "" {
                port = "8080"
        }
        listenSpec := fmt.Sprintf("%s:%s", host, port)

        certFile := os.Getenv("SSL_CERT_FILE")
        keyFile := os.Getenv("SSL_KEY_FILE")

        router := mux.NewRouter()
        router.HandleFunc("/health", HealthHandler).Methods("GET")
        router.HandleFunc("/parser", ParserUrlHandler).Queries("address", "{address}", "language", "{language}", "country", "{country}").Methods("GET")
        router.HandleFunc("/parser", ParserUrlHandler).Queries("address", "{address}", "country", "{country}").Methods("GET")
        router.HandleFunc("/parser", ParserUrlHandler).Queries("address", "{address}", "language", "{language}").Methods("GET")
        router.HandleFunc("/parser", ParserUrlHandler).Queries("address", "{address}").Methods("GET")
        router.HandleFunc("/parser", ParserParamHandler).Queries("language", "{language}", "country", "{country}").Methods("POST")
        router.HandleFunc("/parser", ParserParamHandler).Queries("language", "{language}").Methods("POST")
        router.HandleFunc("/parser", ParserParamHandler).Queries("country", "{country}").Methods("POST")
        router.HandleFunc("/parser", ParserHandler).Methods("POST")
        router.HandleFunc("/expand", ExpandHandler).Methods("POST")

        s := &http.Server{Addr: listenSpec, Handler: router}
        go func() {
                if certFile != "" && keyFile != "" {
                        fmt.Printf("listening on https://%s\n", listenSpec)
                        s.ListenAndServeTLS(certFile, keyFile)
                } else {
                        fmt.Printf("listening on http://%s\n", listenSpec)
                        s.ListenAndServe()
                }
        }()

        stop := make(chan os.Signal)
        signal.Notify(stop, os.Interrupt)
        <-stop
        fmt.Println("\nShutting down the server...")
        ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
        s.Shutdown(ctx)
        fmt.Println("Server stopped")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("READY"))
}

func ParserUrlHandler(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        address := vars[("address")]
        language := vars[("language")]
        country := vars[("country")]
        
        if (len(address) == 0) {
                http.Error(w, "Address cannot be empty", 500)
                return
        }

        parser_options := parser.ParserOptions {
                Language: language,
                Country: country,
        }

        parsed := parser.ParseAddressOptions(address, parser_options)
        parseThing, _ := json.Marshal(parsed)

        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(parseThing))
}

func ExpandHandler(w http.ResponseWriter, r *http.Request) {
        q, _ := ioutil.ReadAll(r.Body)
        address := string(q[:])
        expansions := expand.ExpandAddress(address)
        expansionThing, _ := json.Marshal(expansions)

        w.Header().Set("Content-Type", "application/json")
        w.Write(expansionThing)
}

func ParserHandler(w http.ResponseWriter, r *http.Request) {
        q, _ := ioutil.ReadAll(r.Body)
        address := string(q[:])
        parsed := parser.ParseAddress(address)
        parseThing, _ := json.Marshal(parsed)

        w.Header().Set("Content-Type", "application/json")
        w.Write(parseThing)
}

func ParserParamHandler(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        language := vars[("language")]
        country := vars[("country")]

        q, _ := ioutil.ReadAll(r.Body)
        address := string(q[:])

        parser_options := parser.ParserOptions {
                Language: language,
                Country: country,
        }

        parsed := parser.ParseAddressOptions(address, parser_options)
        parseThing, _ := json.Marshal(parsed)

        w.Header().Set("Content-Type", "application/json")
        w.Write(parseThing)
}
