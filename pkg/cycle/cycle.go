package cycle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type Tester struct {
	Cycles  []Cycle
	Current int
}

type Cycle struct {
	Headers  map[string]string
	Method   string
	Path     string
	Request  []byte
	Code     int
	Response []byte
}

func (c *Cycle) Match(r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	matched := true

	if r.Method != c.Method || r.URL.Path != c.Path || bytes.Compare(c.Request, body) != 0 {
		matched = false
	}

	mh := map[string]string{}

	if c.Headers != nil {
		for k, v := range c.Headers {
			if rv := r.Header.Get(k); rv != v {
				matched = false
				mh[k] = rv
			}
		}
	}

	if !matched {
		return errors.WithStack(fmt.Errorf("cycle\n  headers: %s\n  method: %s\n  path: %s\n  body: %s\nrequest\n  headers: %s\n  method: %s\n  path: %s\n  body: %s", c.Headers, c.Method, c.Path, c.Request, mh, r.Method, r.URL.Path, body))
	}

	return nil
}

func (ct *Tester) Next() *Cycle {
	if ct.Current >= len(ct.Cycles) {
		return nil
	}

	c := ct.Cycles[ct.Current]
	ct.Current++

	return &c
}

func (ct *Tester) Register(c Cycle) {
	if ct.Cycles == nil {
		ct.Cycles = []Cycle{}
	}

	ct.Cycles = append(ct.Cycles, c)
}

func (ct *Tester) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := ct.Next()
	if c == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no cycles remaining\n")
		w.WriteHeader(404)
		return
	}

	if err := c.Match(r); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid cycle\n%s\n", err)
		w.WriteHeader(404)
		return
	}

	if c.Code > 0 {
		w.WriteHeader(c.Code)
	}

	w.Write(c.Response)
}
