package doctest

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
	"time"

	maybetls "github.com/aaw/maybe_tls"
)

//nolint:funlen,gosec,gocognit // keep the test boilerplate compact and ignore security issues since this is for testing only.
// RunTest mocks an test http server.
func RunTest(expectRequest bool, test func()) {
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

	config := tls.Config{Certificates: []tls.Certificate{keypair}, InsecureSkipVerify: true}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("Unbable to listen %s", err))
	}

	mln := maybetls.Listener{
		Listener: listener,
		Config:   &config,
	}

	var gotRequest bool

	// echo server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		gotRequest = true
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

	mux.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		gotRequest = true
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
	mux.HandleFunc("/post", func(writer http.ResponseWriter, request *http.Request) {
		gotRequest = true
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

	srv := http.Server{
		Addr:    listener.Addr().String(),
		Handler: mux,
	}

	go func() {
		if err := srv.Serve(&mln); err != nil {
			if err != http.ErrServerClosed {
				panic(fmt.Sprintf("unable to serve: %s", err))
			}
			return
		}
	}()

	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("tcp", listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // this is a test server so allow insecure TLS connections.
		},
		MaxIdleConns:          30, //nolint:gomnd // this is the default
		IdleConnTimeout:       0,
		TLSHandshakeTimeout:   10 * time.Second, //nolint:gomnd // this is the default
		ExpectContinueTimeout: time.Second,
	}

	test()

	_ = srv.Close()
	_ = listener.Close()

	if expectRequest && !gotRequest {
		panic("expected at least one request")
	}
}
