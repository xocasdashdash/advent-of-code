package day15

import "testing"

func BenchmarkPlayRounds(b *testing.B) {

	g := Game{
		StartNumbers: []int{
			0,
			3,
			6,
		},
	}
	for i := 0; i < b.N; i++ {
		g.PlayNRounds(30000000)
	}
}

func TestPlayRounds(t *testing.T) {
	type tc struct {
		g      Game
		rounds int
		result int
	}
	g := Game{
		StartNumbers: []int{
			0,
			3,
			6,
		},
	}
	testCases := []tc{
		{
			g:      g,
			result: 0,
			rounds: 10,
		},
		{
			g:      g,
			result: 436,
			rounds: 2020,
		},
		{
			g:      g,
			result: 175594,
			rounds: 30000000,
		},
	}
	for _, testCase := range testCases {
		r := testCase.g.PlayNRounds(testCase.rounds)
		if testCase.result != r {
			t.Errorf("Expected %d, got %d", testCase.result, r)
		}
	}
}
