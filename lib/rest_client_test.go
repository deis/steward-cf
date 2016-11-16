package lib

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/arschles/testsrv"
	"github.com/deis/steward-framework"
)

const (
	scheme     = "http"
	host       = "randhost"
	port       = 8080
	timeoutSec = 5
	pathElt1   = "path1"
	pathElt2   = "path2"
	pathElt3   = "path3"
)

var (
	testCfg = config{
		BrokerRequestTimeoutSec: timeoutSec,
	}
	brokerSpec = framework.BrokerSpec{
		URL:      fmt.Sprintf("%s://%s:%d", scheme, host, port),
		Username: "testuser",
		Password: "testpass",
	}
)

func strSliceEq(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, elt := range s1 {
		if elt != s2[i] {
			return false
		}
	}
	return true
}

func urlValuesEq(v1, v2 url.Values) bool {
	lengths := len(v1) == len(v2)
	if !lengths {
		return false
	}
	for key, vals := range v1 {
		if len(vals) != len(v2[key]) || !strSliceEq(vals, v2[key]) {
			return false
		}
	}
	return true
}

var (
	queryStr = url.Values(map[string][]string{
		"key1": []string{"val1"},
		"key2": []string{"val2"},
	})
)

func TestURLStr(t *testing.T) {
	rc := newRESTClient(testCfg)
	urlStr := rc.urlStr(brokerSpec, pathElt1, pathElt2, pathElt3)
	assert.Equal(
		t,
		urlStr,
		fmt.Sprintf("%s/%s/%s/%s", rc.fullBaseURL(brokerSpec), pathElt1, pathElt2, pathElt3),
		"url string",
	)
}

func TestFullBaseURL(t *testing.T) {
	rc := newRESTClient(testCfg)
	fullURL := rc.fullBaseURL(brokerSpec)
	assert.Equal(
		t,
		fullURL,
		fmt.Sprintf("%s://%s:%s@%s:%d", scheme, brokerSpec.Username, brokerSpec.Password, host, port),
		"full url string",
	)
}

func testReq(
	r *http.Request,
	method,
	scheme,
	host string,
	port int,
	user,
	pass string,
	queryStr url.Values,
	pathElts ...string,
) error {
	if r.Method != method {
		return fmt.Errorf("method %s wasn't expected %s", r.Method, method)
	}
	if r.URL.Scheme != scheme {
		return fmt.Errorf("scheme %s wasn't expected %s", r.URL.Scheme, scheme)
	}
	hostSpl := strings.Split(r.URL.Host, ":")
	if len(hostSpl) != 2 {
		return fmt.Errorf("invalid host string %s", r.URL.Host)
	}
	actualHost := hostSpl[0]
	actualPort, err := strconv.Atoi(hostSpl[1])
	if err != nil {
		return fmt.Errorf("invalid port %s (%s)", hostSpl[1], err)
	}
	if actualHost != host {
		return fmt.Errorf("host %s wasn't expected %s", actualHost, host)
	}
	if actualPort != port {
		return fmt.Errorf("port %d wasn't expected %d", actualPort, port)
	}
	if user != r.URL.User.Username() {
		return fmt.Errorf("username %s wasn't expected %s", r.URL.User.Username(), user)
	}
	realPass, _ := r.URL.User.Password()
	if pass != realPass {
		return fmt.Errorf("password %s wasn't expected %s", realPass, pass)
	}
	if !urlValuesEq(r.URL.Query(), queryStr) {
		return fmt.Errorf("query string %s wasn't expected %s", r.URL.Query(), queryStr)
	}
	expectedPath := "/" + strings.Join(pathElts, "/")
	if r.URL.Path != expectedPath {
		return fmt.Errorf("path %s wasn't expected %s", r.URL.Path, expectedPath)
	}
	return nil
}

func TestGet(t *testing.T) {
	rc := newRESTClient(testCfg)
	getReq, err := rc.Get(brokerSpec, queryStr, pathElt1, pathElt2, pathElt3)
	assert.NoErr(t, err)
	assert.NoErr(
		t,
		testReq(
			getReq,
			"GET",
			scheme,
			host,
			port,
			brokerSpec.Username,
			brokerSpec.Password,
			queryStr,
			pathElt1,
			pathElt2,
			pathElt3,
		),
	)
}

func TestPut(t *testing.T) {
	rc := newRESTClient(testCfg)
	putReq, err := rc.Put(brokerSpec, queryStr, nil, pathElt1, pathElt2, pathElt3)
	assert.NoErr(t, err)
	assert.NoErr(
		t,
		testReq(
			putReq,
			"PUT",
			scheme,
			host,
			port,
			brokerSpec.Username,
			brokerSpec.Password,
			queryStr,
			pathElt1,
			pathElt2,
			pathElt3,
		),
	)
}

func TestDelete(t *testing.T) {
	rc := newRESTClient(testCfg)
	delReq, err := rc.Delete(brokerSpec, queryStr, pathElt1, pathElt2, pathElt3)
	assert.NoErr(t, err)
	assert.NoErr(
		t,
		testReq(
			delReq,
			"DELETE",
			scheme,
			host,
			port,
			brokerSpec.Username,
			brokerSpec.Password,
			queryStr,
			pathElt1,
			pathElt2,
			pathElt3,
		),
	)
}

func TestDo(t *testing.T) {
	srv := testsrv.StartServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	u, err := url.Parse(srv.URLStr())
	assert.NoErr(t, err)
	hostSpl := strings.Split(u.Host, ":")
	assert.True(t, len(hostSpl) == 2, "host string was invalid")
	assert.NoErr(t, err)
	rc := newRESTClient(config{
		BrokerRequestTimeoutSec: timeoutSec,
	})
	req, err := http.NewRequest("GET", srv.URLStr(), nil)
	assert.NoErr(t, err)
	ctx := context.Background()
	cancCtx, canc := context.WithCancel(ctx)
	defer canc()
	resp, err := rc.Do(cancCtx, req)
	assert.NoErr(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "status code")
	reqs := srv.AcceptN(1, 1*time.Second)
	assert.Equal(t, len(reqs), 1, "number of requests")

	cancCtx, canc = context.WithCancel(ctx)
	canc()
	resp, err = rc.Do(cancCtx, req)
	assert.True(t, err != nil, "no error returned when expected")
	assert.Nil(t, resp, "response")
}
