package files

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HTTPServer(serverAddr, srvDir string, allowPut bool, uploadToken string) {
	fs := http.FileServer(http.Dir(srvDir))

	if allowPut {
		log.Println("enabling file uploads")
		http.HandleFunc("PUT /", authPutHandler(uploadToken))
	}

	http.Handle("/", addHeaders(fs))

	log.Printf("Serving %s on %s\n", srvDir, serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		for k, v := range r.Header {
			log.Printf(">\tHeader[%q] = %q\n", k, v)
		}

		// There seems to be a race condition in fd_unix.go FD.Write with ignoringEINTRIO()
		// If you put a break point right after `n, err := ignoringEINTRIO(syscall.Write, fd.Sysfd, p[nn:max])`
		// everything seems to work.
		//fs.ServeHTTP(w, r)

		// Just creating a wrapper seems to fix this. I imagine it is because the additional time to resolve the parent
		// methods is all we need. so if it starts failing again, it means something became slower (maybe the S1 agent?)
		// if that happens, try uncommenting the additional methods.
		fs.ServeHTTP(&MyResponseWriter{ResponseWriter: w}, r)
	}
}

type MyResponseWriter struct {
	http.ResponseWriter
}

//func (w *MyResponseWriter) Write(what []byte) (int, error) {
//	log.Println("Method Called: write", string(what[0:20]))
//	return w.ResponseWriter.Write(what)
//}

//func (w *MyResponseWriter) WriteHeader(code int) {
//	log.Println("Method Called: writeHeader", code)
//
//	w.ResponseWriter.WriteHeader(code)
//
//}

func authPutHandler(authToken string) func(http.ResponseWriter, *http.Request) {
	validAuthToken := fmt.Sprintf("Bearer %s", authToken)
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is PUT
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Authorization") != validAuthToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		uploadFile := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
		uploadDir := filepath.Dir(uploadFile)

		log.Printf("upload path: %s\n", uploadFile)

		err := os.MkdirAll(uploadDir, 0777)
		if err != nil {
			http.Error(w, errors.New("unable to upload file").Error(), http.StatusInternalServerError)
			return
		}

		if Exists(uploadFile) {
			http.Error(w, errors.New("unable to upload file").Error(), http.StatusConflict)
			return
		}

		saveFile, err := os.Create(uploadFile)
		if err != nil {
			log.Println("failed to create file", err)
			http.Error(w, errors.New("failed to create file").Error(), http.StatusInternalServerError)
			return
		}
		defer saveFile.Close()
		defer r.Body.Close()
		io.Copy(saveFile, r.Body)

		// Send a response
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "PUT request successfully processed!")
	}
}
