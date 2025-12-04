package files

import (
	"io"
	"net/http"
	"net/url"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetURLReader(parseURL *url.URL) (io.Reader, error) {
	req, err := http.NewRequest("GET", parseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone)")
}

func DownloadURL() {}
