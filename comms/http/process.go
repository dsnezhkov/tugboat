package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dsnezhkov/tugboat/payloads"
)

type HttpComm struct {
	HttpClient http.Client
	Payman     *payloads.PayloadManager
}

type Method int

const (
	Get  Method = 0
	Post        = 1
)

var httpComm *HttpComm

type HttpCommOptions struct {
	Timeout       time.Duration
	Payman        *payloads.PayloadManager
	SSLCertVerify bool
}

func NewHttpComm(options HttpCommOptions) *HttpComm {
	httpTransport := http.DefaultTransport
	if httpComm == nil {
		httpComm = new(HttpComm)
		// TODO: expose client parameters
		if !options.SSLCertVerify {
			httpTransport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
		httpComm.HttpClient = http.Client{
			Timeout: options.Timeout * time.Second,
		}
		httpComm.HttpClient.Transport = httpTransport
		httpComm.Payman = options.Payman
	}
	return httpComm
}

func (httpc *HttpComm) Get() Method {
	return Get
}
func (httpc *HttpComm) Post() int {
	return Post
}
func (httpc *HttpComm) Fetch(url string, method Method) ([]byte, error) {
	var (
		content []byte
		err     error
	)

	switch method {
	case Get:
		content, err = httpc.get(url)
	case Post:
		return nil, fmt.Errorf("not implemented HTTP method: %d", method)
	default:
		return nil, fmt.Errorf("invalid HTTP method: %d", method)
	}
	return content, err
}
func (httpc *HttpComm) Fetch2FS(url string, method Method, storePayToLoc string) (int, error) {
	var (
		content []byte
		err     error
		written int
	)

	content, err = httpc.Fetch(url, method)
	if err != nil {
		return 0, err
	}

	if storePayToLoc == "" {
		ix := strings.LastIndex(url, "/")
		storePayToLoc = url[ix:]
	}
	fmt.Printf("storePayToLoc: %s", storePayToLoc)
	written, err = httpc.Payman.SetDynamicPayload(storePayToLoc, content)
	if err != nil {
		return written, err
	}
	return written, err
}

func (httpc *HttpComm) get(url string) ([]byte, error) {
	response, err := httpc.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
