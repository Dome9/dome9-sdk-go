package client

import (
	"dome9"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux *http.ServeMux
	client *Client
	server *httptest.Server
	config *dome9.Config
)

func setupMuxConfig() *dome9.Config {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	config = dome9.DefaultConfig()
	err := config.SetBaseURL(server.URL)
	if err != nil {
		panic(err)
	}
	return config
}

func teardown()  {
	server.Close()
}
type dummyStruct struct {
	ID int `json:"id"`
}

const getResponse = `{"id": 1234}`

func TestClient_NewRequestDo(t *testing.T) {

	type args struct {
		method string
		url    string
		body   interface{}
		v      interface{}
	}
	tests := []struct {
		name       string
		args       args
		muxHandler func(w http.ResponseWriter, r *http.Request)
		wantResp   *http.Response
		wantErr    bool
		wantVal    *dummyStruct
	}{
		// NewRequestDo test cases
		{
			name: "GET happy path",
			args: struct {
				method string
				url    string
				body   interface{}
				v      interface{}
			}{
				method: "GET",
				url: "/test",
				body: nil,
				v: new(dummyStruct),
			},
			muxHandler: func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(getResponse))
				w.Header().Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
			},
			wantResp: &http.Response{
				StatusCode: 200,
			},
			wantVal: &dummyStruct{
				ID: 1234,
			},
		},
	}

	for _, tt := range tests {
		client := NewClient(setupMuxConfig())
		client.WriteLog("Server URL: %v", client.Config.BaseURL)
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.args.url, tt.muxHandler)
			got, err := client.NewRequestDo(tt.args.method, tt.args.url, tt.args.body, tt.args.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewRequestDo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp.StatusCode != got.StatusCode{
				t.Errorf("Client.NewRequestDo() = %v, want %v", got, tt.wantResp)
			}

			if !reflect.DeepEqual(tt.args.v, tt.wantVal) {
				t.Errorf("returned %#v; want %#v", tt.args.v, tt.wantVal)
			}
		})
	}
	teardown()
}

func TestNewClient(t *testing.T) {
	type args struct {
		config *dome9.Config
	}
	tests := []struct {
		name  string
		args  args
		wantC *Client
	}{
		// NewClient test cases
		{
			name: "Successful Client creation with default config values",
			args: struct{ config *dome9.Config }{config: nil},
			wantC: &Client{dome9.DefaultConfig()},
		},
		{
			name: "Successful Client creation with custom config values",
			args: struct{ config *dome9.Config }{config: &dome9.Config{
				BaseURL: &url.URL{Host: "https://otherhost.com"},

			}},
			wantC: &Client{&dome9.Config{
				BaseURL:    &url.URL{Host: "https://otherhost.com"},
			}},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotC := NewClient(tt.args.config); !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewClient() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
