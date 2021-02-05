// Package server provides a test server that can be used as a test environment for the hit package.
package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	maybetls "github.com/aaw/maybe_tls"
	"github.com/gorilla/mux"
)

// Server is a test server that can be used for the test environment.
type Server struct {
	listener     net.Listener
	server       *http.Server
	requestCount uint
	mu           sync.Mutex
}

// Transport returns a http.Transport that uses the Server.
func (s *Server) Transport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("tcp", s.listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // this is a test server so allow insecure TLS connections.
		},
		MaxIdleConns:          30, //nolint:gomnd // this is the default
		IdleConnTimeout:       0,
		TLSHandshakeTimeout:   10 * time.Second, //nolint:gomnd // this is the default
		ExpectContinueTimeout: time.Second,
	}
}

// RequestCount returns the count of requests that have been made to the server.
func (s *Server) RequestCount() uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.requestCount
}

// Close closes the server.
func (s *Server) Close() error {
	if err := s.server.Close(); err != nil {
		return err
	}
	return s.listener.Close()
}

//nolint:funlen,gocognit //ignore the size, this is not production ready code.
// NewServer starts a new test server.
func NewServer() *Server {
	s := &Server{}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Error generating RSA key: %s", err))
	}
	key := x509.MarshalPKCS1PrivateKey(priv)
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"example.com"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	var cert []byte
	cert, err = x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic(fmt.Sprintf("Failed to create certificate: %s", err))
	}

	pemKey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: key})
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})

	keypair, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		panic(fmt.Sprintf("Error generating keypair: %s", err))
	}

	config := tls.Config{
		Certificates:       []tls.Certificate{keypair},
		InsecureSkipVerify: true, //nolint:gosec // allow all certificates
	}
	s.listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("Unbable to listen %s", err))
	}

	mln := maybetls.Listener{
		Listener: s.listener,
		Config:   &config,
	}

	// echo server
	r := mux.NewRouter()
	r.HandleFunc("/index.html", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		writer.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(writer, "Hello World")
	})

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		writer.Header()["Date"] = nil
		for k, v := range request.Header {
			writer.Header()[k] = v
		}

		for k := range request.Trailer {
			writer.Header().Add("Trailers", k)
		}

		writer.WriteHeader(http.StatusOK)

		n, _ := io.Copy(writer, request.Body)
		if n == 0 {
			_, _ = io.WriteString(writer, "Hello World")
		}

		for k, v := range request.Trailer {
			writer.Header()[k] = v
		}
	})

	r.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		writer.Header()["Date"] = nil
		for k, v := range request.Header {
			writer.Header()[k] = v
		}

		writer.Header().Set("Content-Type", "application/json")

		for k := range request.Trailer {
			writer.Header().Add("Trailers", k)
		}

		writer.WriteHeader(http.StatusOK)

		n, _ := io.Copy(writer, request.Body)
		if n == 0 {
			_, _ = io.WriteString(writer, `{"ID": 10,"Name":"Joe","Roles":["Admin", "User"]}`)
		}

		for k, v := range request.Trailer {
			writer.Header()[k] = v
		}
	})

	// this endpoint should mimic httpbin.org/post
	r.HandleFunc("/post", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = io.WriteString(writer, `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<title>405 Method Not Allowed</title>
<h1>Method Not Allowed</h1>
<p>The method is not allowed for the requested URL.</p
`)
			return
		}

		var jsonData interface{}
		var data string

		if request.Body != nil {
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			data = string(body)

			_ = json.Unmarshal(body, &jsonData)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(writer).Encode(map[string]interface{}{
			"args":    []interface{}{},
			"data":    data,
			"files":   []interface{}{},
			"form":    []interface{}{},
			"headers": request.Header,
			"json":    jsonData,
			"origin":  request.RemoteAddr,
			"url":     "/post",
		})
	})

	r.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = io.WriteString(writer, `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<title>405 Method Not Allowed</title>
<h1>Method Not Allowed</h1>
<p>The method is not allowed for the requested URL.</p
`)
			return
		}

		_ = json.NewEncoder(writer).Encode(map[string]interface{}{
			"args":    []interface{}{},
			"headers": request.Header,
			"origin":  request.RemoteAddr,
			"url":     "/post",
		})
	})

	// this endpoint should mimic httpbin.org/status/

	statusHandler := func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		i, err := strconv.ParseInt(mux.Vars(request)["code"], 10, 32)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(int(i))
	}

	r.HandleFunc("/status/{code:[0-9]+}", statusHandler)
	r.HandleFunc("/cookies/set/{cookie}/{value}", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		vars := mux.Vars(request)
		http.SetCookie(writer, &http.Cookie{
			Name:    vars["cookie"],
			Value:   vars["value"],
			Path:    "/",
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})
		writer.Header().Set("Location", "/cookies")
		writer.WriteHeader(http.StatusFound)
		_, _ = io.WriteString(writer, `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<title>Redirecting...</title>
<h1>Redirecting...</h1>
<p>You should be redirected automatically to target URL: <a href="/cookies">/cookies</a>.  If not click the link.`)
	})
	r.HandleFunc("/cookies", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()
		writer.WriteHeader(http.StatusOK)

		cookies := make(map[string]string)

		for _, cookie := range request.Cookies() {
			cookies[cookie.Name] = cookie.Value
		}

		_ = json.NewEncoder(writer).Encode(struct {
			Cookies map[string]string `json:"cookies"`
		}{
			Cookies: cookies,
		})
	})

	r.HandleFunc("/basic-auth/{username}/{password}", func(writer http.ResponseWriter, request *http.Request) {
		s.mu.Lock()
		s.requestCount++
		s.mu.Unlock()

		vars := mux.Vars(request)

		username, password, ok := request.BasicAuth()
		if !ok || vars["username"] != username || vars["password"] != password {
			writer.Header().Set("WWW-Authenticate", `Basic realm="test"`)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	s.server = &http.Server{
		Addr:    s.listener.Addr().String(),
		Handler: r,
	}

	go func() {
		if err := s.server.Serve(&mln); err != nil {
			if err != http.ErrServerClosed {
				panic(fmt.Sprintf("unable to serve: %s", err))
			}
			return
		}
	}()

	return s
}
