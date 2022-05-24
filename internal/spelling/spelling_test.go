package spelling

import (
	"reflect"
	"sync"
	"testing"
)

func Test_splitWord(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want []split
	}{
		{
			name: "ok",
			args: args{
				word: "cabbage",
			},
			want: []split{
				{left: "", right: "cabbage"},
				{left: "c", right: "abbage"},
				{left: "ca", right: "bbage"},
				{left: "cab", right: "bage"},
				{left: "cabb", right: "age"},
				{left: "cabba", right: "ge"},
				{left: "cabbag", right: "e"},
				{left: "cabbage", right: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitWord(tt.args.word); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deletes(t *testing.T) {
	type args struct {
		splits []split
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "ok",
			args: args{
				splits: []split{
					{left: "", right: "word"},
					{left: "w", right: "ord"},
					{left: "wo", right: "rd"},
					{left: "wor", right: "d"},
					{left: "word", right: ""},
				},
			},
			wantResult: []string{
				"ord",
				"wrd",
				"wod",
				"wor",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			edits := make(chan string)
			go deletes(tt.args.splits, wg, edits)
			go func() {
				wg.Wait()
				close(edits)
			}()
			if gotResult := chanToSlice(edits); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("deletes() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_transposes(t *testing.T) {
	type args struct {
		splits []split
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "ok",
			args: args{
				splits: []split{
					{left: "", right: "word"},
					{left: "w", right: "ord"},
					{left: "wo", right: "rd"},
					{left: "wor", right: "d"},
					{left: "word", right: ""},
				},
			},
			wantResult: []string{
				"owrd",
				"wrod",
				"wodr",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			edits := make(chan string)
			go transposes(tt.args.splits, wg, edits)
			go func() {
				wg.Wait()
				close(edits)
			}()
			if gotResult := chanToSlice(edits); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("transposes() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_replaces(t *testing.T) {
	type args struct {
		splits []split
	}
	tests := []struct {
		name          string
		args          args
		wantResultLen int
	}{
		{
			name: "ok",
			args: args{
				splits: []split{
					{left: "", right: "word"},
					{left: "w", right: "ord"},
					{left: "wo", right: "rd"},
					{left: "wor", right: "d"},
					{left: "word", right: ""},
				},
			},
			wantResultLen: len("word") * len([]rune(letters)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			edits := make(chan string)
			go replaces(tt.args.splits, wg, edits)
			go func() {
				wg.Wait()
				close(edits)
			}()
			if gotResultLen := len(chanToSlice(edits)); gotResultLen != tt.wantResultLen {
				t.Errorf("len of replaces() = %v, want %v", gotResultLen, tt.wantResultLen)
			}
		})
	}
}

func Test_inserts(t *testing.T) {
	type args struct {
		splits []split
	}
	tests := []struct {
		name          string
		args          args
		wantResultLen int
	}{
		{
			name: "ok",
			args: args{
				splits: []split{
					{left: "", right: "word"},
					{left: "w", right: "ord"},
					{left: "wo", right: "rd"},
					{left: "wor", right: "d"},
					{left: "word", right: ""},
				},
			},
			wantResultLen: (len("word") + 1) * len([]rune(letters)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			edits := make(chan string)
			go inserts(tt.args.splits, wg, edits)
			go func() {
				wg.Wait()
				close(edits)
			}()
			if gotResultLen := len(chanToSlice(edits)); gotResultLen != tt.wantResultLen {
				t.Errorf("len of inserts() = %v, want %v", gotResultLen, tt.wantResultLen)
			}
		})
	}
}

func Test_swapFirstSymbols(t *testing.T) {
	type args struct {
		runes []rune
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "english",
			args: args{
				runes: []rune("water"),
			},
			want: "awter",
		},
		{
			name: "russian",
			args: args{
				runes: []rune("вода"),
			},
			want: "овда",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := swapFirstRunes(tt.args.runes); got != tt.want {
				t.Errorf("swapFirstSymbols() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleEdits(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name          string
		args          args
		wantResultLen int
	}{
		{
			name: "ok",
			args: args{
				word: "cabbage",
			},
			wantResultLen: len("cabbage") + // deletes
				len("cabbage") - 1 + // transposes
				len("cabbage")*len([]rune(letters)) + // replaces
				(len("cabbage")+1)*len([]rune(letters)), // inserts
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResultLen := len(SimpleEdits(tt.args.word)); gotResultLen != tt.wantResultLen {
				t.Errorf("len of SimpleEdits() = %v, want %v", gotResultLen, tt.wantResultLen)
			}
		})
	}
}
