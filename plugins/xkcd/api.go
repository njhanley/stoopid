package xkcd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Info struct {
	Year, Month, Day string
	Title, Alt       string
	Num              int
	Img              string
}

func parse(b []byte) (*Info, error) {
	var x Info
	err := json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func get(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

// empty num for current comic
func Get(num string) (*Info, error) {
	url := "https://xkcd.com/"
	if num != "" {
		url += num + "/"
	}
	url += "info.0.json"

	b, err := get(url)
	if err != nil {
		return nil, err
	}

	return parse(b)
}

var noRedirectClient = &http.Client{
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func GetRandom() (*Info, error) {
	r, err := noRedirectClient.Get("https://c.xkcd.com/random/comic")
	if err != nil {
		return nil, err
	}
	r.Body.Close()

	u, err := r.Location()
	if err != nil {
		return nil, err
	}

	b, err := get(u.String() + "info.0.json")
	if err != nil {
		return nil, err
	}

	return parse(b)
}
