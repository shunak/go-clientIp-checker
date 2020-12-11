package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
    "regexp"
    "errors"
)

func getConnectionData() string {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// header 見たいので request の dump を response で返す
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, string(dump))
	}))

	rpURL, err := url.Parse(backendServer.URL)
	if err != nil {
		log.Fatal(err)
	}
	frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(rpURL))
	defer frontendProxy.Close()

	// リクエスト定義
	req, err := http.NewRequest(http.MethodGet, frontendProxy.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Connection ヘッダー追加
	req.Header.Set("Connection", "keep-alive")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%s", b)
	return string(b)
	// return string(frontendProxy.URL)
}
func extractCientIp() string {
	res := getConnectionData()
	reg := regexp.MustCompile(`X-.*:.*[\d]`)
	xForwardedFor := reg.FindString(res)
    clientIp := xForwardedFor[17:]
    return clientIp
}

func genError() (bool, error) {
    return false, errors.New("you don't have access permission.")
}

func checkIp(clientIp string) (bool, error) {
    const allowIp = "127.0.0.1"
    if clientIp != allowIp {
        return genError()
    } 
    return true, nil
}

func main() {
    ip := extractCientIp()
    _, err := checkIp(ip)
    if err != nil {
        panic(err)
    }
	fmt.Printf("your ip is %v. you are permitted access!\n", ip)
}
