package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	client  *http.Client
	options options
}

func (h *HttpClient) Post(ctx context.Context, url string, body interface{}, v interface{}, opts ...RequestOption) (err error) {
	opt := RequestOptions{
		ContentType: "application/json",
		Header:      http.Header{},
	}
	for _, o := range opts {
		o(&opt)
	}

	opt.Header.Set("Content-Type", opt.ContentType)
	reader, err := bodyReader(body)
	if err != nil {
		return
	}

	resp, err := h.do("POST", url, reader, opt.Header)
	if err != nil {
		return
	}

	if opt.RespHandler != nil {
		return opt.RespHandler(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code:%v url:%v", resp.StatusCode, url)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, v)
	return
}

func bodyReader(body interface{}) (io.Reader, error) {
	if _, ok := body.(io.Reader); ok {
		return body.(io.Reader), nil
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func (h *HttpClient) Get(ctx context.Context, url string, v interface{}) (err error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusNoContent {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code:%v url:%v", resp.StatusCode, url)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, v)
	return
}

func (h *HttpClient) do(method string, url string, reader io.Reader, header http.Header) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header = header
	return h.client.Do(req)
}

func NewHttpClient(opts ...Option) *HttpClient {
	opt := options{
		maxConnectionNum:    100,
		timeout:             20 * time.Second,
		dialTimeout:         30 * time.Second,
		idleConnTimeout:     90 * time.Second,
		keepAlive:           30 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}

	return &HttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   opt.dialTimeout,
					KeepAlive: opt.keepAlive,
				}).DialContext,
				MaxIdleConns:          opt.maxConnectionNum,
				MaxIdleConnsPerHost:   opt.maxConnectionNum,
				IdleConnTimeout:       opt.idleConnTimeout,
				TLSHandshakeTimeout:   opt.tlsHandshakeTimeout,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: opt.timeout,
		},
		options: opt,
	}
}
