package spelling

import (
	"sync"
)

const letters = "abcdefghijklmnopqrstuvwxyzабвгдеёжзиклмнопрстуфхцчшщъыьэюя"

type split struct {
	left, right []rune
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
		if len(s.right) == 0 {
			continue
		}
		edits <- string(append(s.left, s.right[1:]...))
	}

	wg.Done()
}

func transposes(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		if len(s.right) < 2 {
			continue
		}
		edits <- string(append(s.left, swapFirstChars(s.right)...))
	}

	wg.Done()
}

func replaces(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		if len(s.right) == 0 {
			continue
		}
		for _, l := range letters {
			edit := append(s.left, l)
			edit = append(edit, s.right[1:]...)
			edits <- string(edit)
		}
	}

	wg.Done()
}

func inserts(splits []split, wg *sync.WaitGroup, edits chan<- string) {
	for _, s := range splits {
		for _, l := range letters {
			edit := append(s.left, l)
			edit = append(edit, s.right...)
			edits <- string(edit)
		}
	}

	wg.Done()
}

func splitWord(word string) []split {
	runes := []rune(word)
	splits := make([]split, len(runes)+1)

	for i := 0; i <= len(runes); i++ {
		splits[i] = split{
			left:  runes[:i],
			right: runes[i:],
		}
	}

	return splits
}

func swapFirstChars(runes []rune) []rune {
	runes[1], runes[0] = runes[0], runes[1]
	return runes
}

func chanToSlice[T any](c chan T) (result []T) {
	for entry := range c {
		result = append(result, entry)
	}
	return
}
