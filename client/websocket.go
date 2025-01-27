package client

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func (c *Client) wrapWSS(conn net.Conn) net.Conn {
	return c.wrapWS(tls.Client(conn, c.TLSConfig))
}

func (c *Client) wrapWS(conn net.Conn) net.Conn {
	return newWSWrapper(conn, c)
}

type wsWrapper struct {
	net.Conn
	client *Client
	buf    *bytes.Buffer

	wsConn *websocket.Conn
}

func newWSWrapper(conn net.Conn, client *Client) *wsWrapper {
	return &wsWrapper{
		Conn:   conn,
		client: client,
		buf:    bytes.NewBuffer(make([]byte, 0, 1024)),
		wsConn: nil,
	}
}

func (w *wsWrapper) Read(b []byte) (n int, err error) {
	if w.wsConn == nil {
		w.wsConn, err = w.handshake()
		if err != nil {
			return
		}
	}

	if w.buf.Len() > 0 {
		return w.buf.Read(b)
	}

	_, p, err := w.wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}
	n = copy(b, p)
	w.buf.Write(p[n:])

	return
}

func (w *wsWrapper) Write(b []byte) (n int, err error) {
	if w.wsConn == nil {
		w.wsConn, err = w.handshake()
		if err != nil {
			return
		}
	}

	err = w.wsConn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *wsWrapper) handshake() (conn *websocket.Conn, err error) {
	config := w.client.Config
	log.Printf("[websocket] upgrade to websocket at %s", config.WSPath)
	u := url.URL{
		Scheme: "ws",
		Host:   config.ServerAddr,
		Path:   config.WSPath,
	}
	var header http.Header
	if config.Username != "" && config.Password != "" {
		header = make(http.Header)
		s := base64.StdEncoding.EncodeToString([]byte(config.Username + ":" + config.Password))
		header.Add("Authorization", "Basic "+s)
	}

	wsDial := websocket.Dialer{
		NetDial: func(net, addr string) (net.Conn, error) {
			return w.Conn, nil
		},
	}
	conn, res, err := wsDial.Dial(u.String(), header)
	if err == nil {
		log.Printf("[websocket] connection established: %s", res.Status)
	}
	return
}
