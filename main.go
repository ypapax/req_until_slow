package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ypapax/logrus_conf"
	"net/http"
	"strings"
	"time"
)

func main() {
	var (
		u       string
		timeout time.Duration
		randomQuery bool
	)
	flag.StringVar(&u, "url", "", "url to request")
	flag.DurationVar(&timeout, "timeout", time.Second, "request timeout")
	flag.BoolVar(&randomQuery, "random-query", true, "ads random query to avoid caching")
	flag.Parse()
	if err := func() error {
		if err := logrus_conf.PrepareFromEnv("req_until_slow"); err != nil {
			return errors.WithStack(err)
		}
		if len(u) == 0 {
			return errors.Errorf("missing url")
		}
		if timeout == 0 {
			return errors.Errorf("missing timeout")
		}
		i := 0
		for {
			i++
			l := logrus.WithField("i", i).WithField("timeout", timeout)
			t1 := time.Now()
			uFull := u
			if randomQuery {
				const querySep = "?"
				var prefix = querySep
				if strings.Contains(uFull, querySep) {
					prefix = "&"
				}
				uFull+=prefix+"random="+fmt.Sprintf("%+v", time.Now().UnixNano())
			}
			l = l.WithField("url", uFull)
			if err := req(uFull, timeout); err != nil {
				return errors.Wrapf(err, "for i %+v, time spent: %+v", i, time.Since(t1))
			}
			l.Infof("requested for %+v", time.Since(t1))
		}
	}(); err != nil {
		logrus.Errorf("err: %+v", err)
	}
}

func req(u string, timeout time.Duration) error {
	cl := http.Client{Timeout: timeout}
	res, err := cl.Get(u)
	if err != nil {
		return errors.WithStack(err)
	}
	if res.StatusCode >= 400 {
		return errors.Errorf("bad status code: %+v", res.StatusCode)
	}
	return nil
}