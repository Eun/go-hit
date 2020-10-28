maybe_tls
=========

An implementation of golang's
[`net.Listener` interface](https://golang.org/pkg/net/#Listener) that accepts both
TCP and TLS-over-TCP connections on the same port.

Example
-------

To use, create a `net.Listener` and a `tls.Config` object first, then wrap them
in a `maybe_tls.Listener`:

```
config := tls.Config{Certificates: []tls.Certificate{keypair}}
ln, err := net.Listen("tcp", ":1234")
if err != nil {
  // Handle error
}
mln := maybe_tls.Listener{ln, &config}
```
Then use the `maybe_tls.Listener` just as you would a `net.Listener`. You can
detect whether the accepted connection is encrypted or not by testing whether
it's a `tls.Conn`:

```
conn, err := mln.Accept()
switch conn.(type) {
default:
    // It's a plain TCP connection
case *tls.Conn:
    // It's a TLS connection
}
```

Restrictions
------------

`maybe_tls.Listener` works by trying to figure out whether the first few bytes
that a client sends look like a TLS ClientHello message. If so, the data from
the connection is "rewound" and wrapped in a
[`tls.Conn`](https://golang.org/pkg/crypto/tls/#Conn). If not, the data from the
connection is rewound and replayed as is. Because of this detection mechanism,
the `maybe_tls.Listener` can only be used with application protocols where the
client is expected to send the first message after a connection is made.

Installation
------------

    go get github.com/aaw/maybe_tls