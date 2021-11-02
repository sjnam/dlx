package dlx

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestDLX(t *testing.T) {
	dlxInput := `
| A simple example
A B C D E | F G
C E F
A D G
B C F
A D
B G
D E G
`
	xccInput := `
|A simple example of color controls
A B C | X Y
A B X:0 Y:0
A C X:1 Y:1
X:0 Y:1
B X:1
C Y:1
`
	cases := []struct {
		in, want string
	}{
		{dlxInput, "[[A D] [B G] [C E F]]"},
		{xccInput, "[[A C X:1 Y:1] [B X:1]]"},
	}
	for _, c := range cases {
		d, err := NewDLX(strings.NewReader(c.in))
		if err != nil {
			t.Fatal(err)
		}
		res := <-d.Dance()
		if res.Err != nil {
			t.Fatal(res.Err)
		}
		sol := res.Solution
		sort.Slice(sol, func(i, j int) bool {
			for _, opt := range sol {
				sort.Strings(opt)
			}
			return sol[i][0] < sol[j][0]
		})
		got := fmt.Sprint(sol)
		if got != c.want {
			t.Errorf("DLX.Dance() == %q, want %q", got, c.want)
		}
	}
}
