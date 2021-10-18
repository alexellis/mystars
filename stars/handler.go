package function

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	input := readBody(r)
	repo := strings.TrimSpace(string(input))

	if len(repo) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Please provide a repository name")
	}

	starCount := getStars(repo)
	fmt.Fprintf(w, "Stars for: %s, %d", repo, starCount)
}

func getStars(repoOwner string) int64 {
	return 28000
}

func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return body
}
