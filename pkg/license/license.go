package license

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
)

const PublicKey = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAu2/MMK91F2AJTYIDeMPX
XziHsQaKz4ODgRy9tiPubqVBinT2ygbuhW9JQ87jSSJfgOcHZNzKTEMT7WpbKLoc
Z4fM6b9D+wfInvU/p7t82mlDAiSyYS0ICj0lbFqMMeZ0i602xlbblnVnJNvcJN0B
QeUfy7rDoHbCUO3FawCduWV6bb7vFiAVJMWr8E5w6NGYrKhQQWuqA9pNzop3nKBS
ntRPTr5DaFQWVRotnytDKfEcncHF2irKgQzt1AnCTlEj3/h4sdsEa+bAa7EqYYy1
R4oBH3qQE/+CwQEftoDapJvbMjAQ7TFBWhxIKLWWMeAYxEiEA5co0I+6OsCzWSgt
x5gHgYEg/1lB7o1aopT2RrGW4ZTgrl5UfV9WzjBJt6m5/ppYvYLbjtYGlA3DUYNQ
oQfVDA/vhrO0Q9KU7r3V+dcQ6263iM5/h+oDr4LTkFIsgx0L0uiCSUsgpf4Jf74h
DK9yPBg88RZ3584+Hc19tg4QRqedDYl9OVfsFfIhvTvfDVbD4RGYBzJONaHWc1B4
y3SCMstQ61IPbqZO2vGU4hSHPXEU3chUcf+a5fF4Y6/jRxhRN7ojgOmwQorNc5+O
gIaWodjEFFAJzfodCeKabngG/wNwLER3DvN+Mf728jeZ19SWKE4De+Gv6U9/hcER
t4UKEuJlUUlrH/50Lhm3m+cCAwEAAQ==
-----END PUBLIC KEY-----
`

var (
	Current *License
)

type envelope struct {
	License   []byte `json:"l"`
	Signature []byte `json:"s"`
}

type License struct {
	Backend string    `json:"backend"`
	Expires time.Time `json:"expires"`
	Public  bool      `json:"public"`
	Seats   int       `json:"seats"`
}

func init() {
	if err := load(); err != nil {
		fmt.Fprintf(os.Stderr, "license error: %s\n", err)
		os.Exit(1)
	}
}

func load() error {
	if os.Getenv("TEST") == "true" {
		return nil
	}

	res, err := http.Get(fmt.Sprintf("https://console:%s@enterprise.convox.com/license?host=%s", settings.LicenseKey, settings.Host))
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	eed, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return errors.WithStack(err)
	}

	gz, err := gzip.NewReader(bytes.NewReader(eed))
	if err != nil {
		return errors.WithStack(err)
	}
	defer gz.Close()

	ed, err := ioutil.ReadAll(gz)
	if err != nil {
		return errors.WithStack(err)
	}

	var e envelope

	if err := json.Unmarshal(ed, &e); err != nil {
		return errors.WithStack(err)
	}

	b, _ := pem.Decode([]byte(PublicKey))

	der, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return errors.WithStack(err)
	}

	pub, ok := der.(*rsa.PublicKey)
	if !ok {
		return errors.WithStack(fmt.Errorf("invalid public key"))
	}

	hash := sha256.Sum256(e.License)

	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], e.Signature); err != nil {
		return errors.WithStack(err)
	}

	var l License

	if err := json.Unmarshal(e.License, &l); err != nil {
		return errors.WithStack(err)
	}

	Current = &l

	// if settings.Development {
	//   ub, err := url.Parse(Current.Backend)
	//   if err != nil {
	//     return errors.WithStack(err)
	//   }

	//   ub.Host = "convox-enterprise-dev.ngrok.io"

	//   Current.Backend = ub.String()
	// }

	fmt.Printf("ns=license host=%s expires=%s seats=%d\n", settings.Host, Current.Expires.Format("2006-01-02"), Current.Seats)

	return nil
}

func (l *License) Expired() bool {
	return l.Expires.Before(time.Now().UTC())
}
