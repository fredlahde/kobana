package errors

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	prev := E("prev", IO, fmt.Errorf("foobar"), P("userid", 2423))
	e := E("test", Other, prev, P("userid", 2423))

	json, err := MarshalJSON(e.(*Error))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(json))
}

func TestMuchError(t *testing.T) {
	e, _ := E(Op("start"), Other, fmt.Errorf("could not load user")).(*Error)
	var prev = e
	i := 5
	for i > 0 {
		var op Op = Op(fmt.Sprintf("test %d", i))

		p, ok := E(op, Other, prev).(*Error)
		if !ok {
			t.Fatal("Could not cast to Errro")
		}
		prev = p
		i = i - 1
	}

	json, err := MarshalJSON(prev)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))
}
