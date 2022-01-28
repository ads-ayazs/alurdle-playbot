package playerbot

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"sync"
	"time"

	"aluance.io/wordleplayer/internal/config"
	"github.com/matryer/resync"
)

type PlayerDictionary interface {
	Remember(string, bool) error
	Generate() (string, error)
	IsValid(string) bool
	DescribeSize(bool) int
}

type playerDictionary struct {
	lexicon map[string]bool
	mu      sync.Mutex
}

var muDictionaries sync.Mutex
var dictionaries = map[string]*playerDictionary{}

func (d *playerDictionary) initialize() {
	d.lexicon = make(map[string]bool)
}

func CreateDictionary(botId string) PlayerDictionary {
	rand.Seed(time.Now().UnixNano())

	muDictionaries.Lock()
	d, ok := dictionaries[botId]
	if !ok {
		d = new(playerDictionary)
		d.initialize()

		dictionaries[botId] = d
	}
	muDictionaries.Unlock()

	return d
}

func (d *playerDictionary) Remember(word string, valid bool) error {
	(*d).mu.Lock()
	d.lexicon[word] = valid
	(*d).mu.Unlock()

	return nil
}

func (d *playerDictionary) Generate() (string, error) {
	w := ""

	// Generate a word and discard if it is known to be invalid
	valid := false
	for ok := true; ok; valid, ok = d.lexicon[w] {
		if ok && valid {
			break
		}
		if rand.Float32() > 0.0 {
			// create a word from random letters
			w = createRandomWord()
		} else {
			// draw a random word from the real words dictionary
			w = injectValidWord()
		}
	}

	return w, nil
}

func createRandomWord() string {
	buf := bytes.NewBufferString("")

	for i := 0; i < config.CONFIG_GAME_WORDLENGTH; i++ {
		letterNum := rand.Intn(26) + 65
		c := string(rune(letterNum))
		buf.WriteString(c)
	}

	return buf.String()
}

func (d *playerDictionary) IsValid(word string) bool {
	if validity, ok := d.lexicon[word]; ok {
		return validity
	}

	return false
}

func (d *playerDictionary) DescribeSize(validOnly bool) int {
	if validOnly {
		validCount := 0
		for _, v := range d.lexicon {
			if v {
				validCount++
			}
		}
		return validCount
	}

	return len(d.lexicon)
}

func injectValidWord() string {
	if err := loadDictionaryFile(""); err != nil {
		return ""
	}

	idx := rand.Intn(len(dictionaryWords))
	return strings.ToUpper(dictionaryWords[idx])
}

var dictionaryLoaded = false
var dictionaryLoad_once resync.Once
var dictionaryWordMap = map[string]bool{}
var dictionaryWords = []string{}

func loadDictionaryFile(filename string) error {

	// Only initialize dictionary once
	if dictionaryLoaded {
		return nil
	}

	if len(filename) < 1 {
		filename = config.CONFIG_DICTIONARY_FILEPATH
	}

	f, err := config.LoadEmbedFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Do this only once (unless reset)
	dictionaryLoad_once.Do(func() {

		// Load only words of configured length from the file
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			word := scanner.Text()
			if len(word) == config.CONFIG_GAME_WORDLENGTH {
				dictionaryWords = append(dictionaryWords, word)
				dictionaryWordMap[word] = true
			}
		}

		if err := scanner.Err(); err != nil {
			return
		}

		dictionaryLoaded = true
	})

	return nil
}

type dict struct {
	init_once  resync.Once
	initalized bool
	words      []string
	wordMap    map[string]bool
}
