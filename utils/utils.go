package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"

	"github.com/tg123/go-htpasswd"
	"golang.org/x/sync/errgroup"
)

// buffer pools
var (
	SPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 576)
		},
	} // small buff pool
	LPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 64*1024+262)
		},
	} // large buff pool for udp
)

// Transport rw1 and rw2
func Transport(rw1, rw2 io.ReadWriter) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		b := LPool.Get().([]byte)
		defer LPool.Put(b)

		_, err := io.CopyBuffer(rw1, rw2, b)
		return err
	})

	g.Go(func() error {
		b := LPool.Get().([]byte)
		defer LPool.Put(b)

		_, err := io.CopyBuffer(rw2, rw1, b)
		return err
	})

	err := g.Wait()
	return err
}

// StrEQ returns whether s1 and s2 are equal
func StrEQ(s1, s2 string) bool {
	return subtle.ConstantTimeCompare([]byte(s1), []byte(s2)) == 1
}

// StrInSlice return whether str in slice
func StrInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str || strings.HasSuffix(str, fmt.Sprintf(".%s", s)) {
			return true
		}
	}
	return false
}

// VerifyByMap returns an verifier that verify by an username-password map
func VerifyByMap(users map[string]string) func(string, string) bool {
	return func(username, password string) bool {
		pw, ok := users[username]
		if !ok {
			return false
		}
		return StrEQ(pw, password)
	}
}

// VerifyByHtpasswd returns an verifier that verify by a htpasswd file
func VerifyByHtpasswd(users string) func(string, string) bool {
	f, err := htpasswd.New(users, htpasswd.DefaultSystems, nil)
	if err != nil {
		log.Fatalf("Load htpasswd file failed: %s", err)
	}
	return func(username, password string) bool {
		return f.Match(username, password)
	}
}

func HttpBasicAuth(auth string, verify func(string, string) bool) bool {
	prefix := "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return false
	}
	auth = strings.Trim(auth[len(prefix):], " ")
	dc, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return false
	}
	groups := strings.Split(string(dc), ":")
	if len(groups) != 2 {
		return false
	}
	return verify(groups[0], groups[1])
}

func PProf() {
	go func() {
		log.Println(http.ListenAndServe(":8090", nil))
	}()
}
