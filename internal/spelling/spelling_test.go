package spelling

import (
	"reflect"
	"sync"
	"testing"
)

var englishSplits = []split{
	{left: []rune(""), right: []rune("word")},
	{left: []rune("w"), right: []rune("ord")},
	{left: []rune("wo"), right: []rune("rd")},
	{left: []rune("wor"), right: []rune("d")},
	{left: []rune("word"), right: []rune("")},
}

var russianSplits = []split{
	{left: []rune(""), right: []rune("огонь")},
	{left: []rune("о"), right: []rune("гонь")},
	{left: []rune("ог"), right: []rune("онь")},
	{left: []rune("ого"), right: []rune("нь")},
	{left: []rune("огон"), right: []rune("ь")},
	{left: []rune("огонь"), right: []rune("")},
}

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
				word: "word",
			},
			want: englishSplits,
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
				splits: englishSplits,
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
				splits: englishSplits,
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
			name: "english",
			args: args{
				splits: englishSplits,
			},
			wantResultLen: len("word") * len([]rune(letters)),
		},
		{
			name: "russian",
			args: args{
				splits: russianSplits,
			},
			wantResultLen: len([]rune("огонь")) * len([]rune(letters)),
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
				splits: englishSplits,
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
		want []rune
	}{
		{
			name: "english",
			args: args{
				runes: []rune("water"),
			},
			want: []rune("awter"),
		},
		{
			name: "russian",
			args: args{
				runes: []rune("огонь"),
			},
			want: []rune("гоонь"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := swapFirstChars(tt.args.runes); !reflect.DeepEqual(got, tt.want) {
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
			if gotResultLen := len(SimpleEdits(tt.args.word, true)); gotResultLen != tt.wantResultLen {
				t.Errorf("len of SimpleEdits() = %v, want %v", gotResultLen, tt.wantResultLen)
			}
		})
	}
}
