package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var (
	statusLineRe = regexp.MustCompile(`^(.+)\s+(\d+)\s+(.+)$`)
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Returns the response code and both to be returned for the given path.
//
// This is read from responses.txt, which contains lines like
// /path 200 path/to/file.json
func responseFor(path string) (int, []byte) {
	f, err := os.Open("responses.txt")
	panicErr(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		match := statusLineRe.FindStringSubmatch(s.Text())
		if len(match) == 0 {
			continue
		}
		candidatePath, codeStr, rspFile := match[1], match[2], match[3]
		if candidatePath == path {
			code, err := strconv.Atoi(codeStr)
			panicErr(err)
			b, err := ioutil.ReadFile(rspFile)
			panicErr(err)
			return code, b
		}
	}
	return http.StatusNotFound, []byte(`{}`)
}

func getEndpoints(rw http.ResponseWriter, req *http.Request) {
	log.Printf("➡️  %v", req.URL)
	status, body := responseFor(req.URL.Path)
	rw.Header().Set("Content-Type", `application/json`)
	rw.WriteHeader(status)
	log.Printf("    ⬅️  %d", status)
	rw.Write(body)
}

func main() {
	handler := http.HandlerFunc(getEndpoints)
	log.Printf("👂  Listening on :5000")
	panic(http.ListenAndServe(":5000", handler))
}
