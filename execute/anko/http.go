package anko

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var reqErr = errors.New("request error")
var reqNotOk = errors.New("request error")

func httpGet(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", reqErr
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%d", resp.StatusCode), reqNotOk
	}
	defer resp.Body.Close()
	db, err := io.ReadAll(resp.Body)
	return string(db), nil
}
