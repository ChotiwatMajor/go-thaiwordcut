package gothaiwordcut

import (
	"github.com/tchap/go-patricia/patricia"
	"os"
	"bufio"
	"regexp"
)


type Segmenter struct {
	Trie *patricia.Trie

	minLength int
}

type Option func(*Segmenter)


func (w *Segmenter) loadFileIntoTrie(filePath string) {
	f, err := os.Open(filePath)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w.Trie.Insert(patricia.Prefix(scanner.Text()), 1)
	}

	check(scanner.Err())
}

func (w *Segmenter) findSegment(c string) []string {
	i := 0
	N := len(c)
	arr := make([]string, 0)
	for i < N {
		// search tree
		j := w.searchTrie(c[i:N])
		if j == "" {
			i = i + 1
		} else {
			arr = append(arr, j)
			i = i + len(j)
		}
	}

	return arr
}

func (w *Segmenter) searchTrie(s string) string {
	// check if the word is latin
	latinResult := simpleRegex("[A-Za-z\\d]*", s)
	if latinResult != "" {
		return latinResult
	}

	// check if its number
	numberResult := simpleRegex("[\\d]*", s)
	if numberResult != "" {
		return numberResult
	}

	var longestWord string

	// loop word character, trying to find longest word
	for i := len(s); i > 0; i-- {
		if w.Trie.Match(patricia.Prefix(s[0:i])) {
			longestWord = s[0:i]
			break
		}
	}

	return longestWord
}

func simpleRegex(expr string, s string) string {
	r, err := regexp.Compile(expr)
	check(err)
	return r.FindString(s)
}

func (w *Segmenter) Segment(txt string) []string {
	return w.findSegment(txt)
}

func Wordcut(options ...Option) *Segmenter {
	segmenter := &Segmenter{}
	segmenter.Trie = patricia.NewTrie()
	return segmenter
}

func (w *Segmenter) LoadDefaultDict() {
	w.loadFileIntoTrie("./dict/lexitron.txt")
}

/*
 * If error, then we should PANIC!
 */
func check(e error) {
	if e != nil {
		panic(e)
	}
}