package spelling

import (
	"sync"
)

const letters = "abcdefghijklmnopqrstuvwxyzабвгдеёжзиклмнопрстуфхцчшщъыьэюя"

type split struct {
	left, right string
}

func SimpleEdits(word string) (result []string) {
	splits := splitWord(word)

	editFuncs := []func(splits []split, wg *sync.WaitGroup, edits chan<- string){
		deletes, transposes, replaces, inserts,
	}

	wg := &sync.WaitGroup{}
	edits := make(chan string)
	for _, f := range editFuncs {
		wg.Add(1)
		go f(splits, wg, edits)
	}

	go func() {
		wg.Wait()
		close(edits)
	}()

	return chanToSlice(edits)
}

func deletes(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		if s.right == "" {
			continue
		}
		edits <- s.left + s.right[1:]
	}

	wg.Done()
}

func transposes(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		runes := []rune(s.right)
		if len(runes) < 2 {
			continue
		}
		edits <- s.left + swapFirstRunes(runes)
	}

	wg.Done()
}

func replaces(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		if s.right == "" {
			continue
		}
		for _, l := range letters {
			edits <- s.left + string(l) + s.right[1:]
		}
	}

	wg.Done()
}

func inserts(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		for _, l := range letters {
			edits <- s.left + string(l) + s.right
		}
	}

	wg.Done()
}

func splitWord(word string) []split {
	splits := make([]split, len(word)+1)

	for i := 0; i <= len(word); i++ {
		splits[i] = split{
			left:  word[:i],
			right: word[i:],
		}
	}

	return splits
}

func swapFirstRunes(runes []rune) string {
	runes[1], runes[0] = runes[0], runes[1]
	return string(runes)
}

func chanToSlice[T any](c chan T) (result []T) {
	for entry := range c {
		result = append(result, entry)
	}
	return
}
