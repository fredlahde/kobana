package errors

import (
	json2 "encoding/json"
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	prev := E("prev", IO, fmt.Errorf("foobar"), P("userid", 2423))
	e := E("test", Other, prev)

	_, err := MarshalJSON(e.(*Error))
	if err != nil {
		t.Fatal(err)
	}
}

func TestMuchError(t *testing.T) {
	e, _ := E(Op("start"), Other, fmt.Errorf("could not load user")).(*Error)
	var prev = e
	i := 15
	for i > 0 {
		var op Op = Op(fmt.Sprintf("test %d", i))

		p, ok := E(op, Other, prev).(*Error)
		if !ok {
			t.Fatal("Could not cast to Error")
		}
		prev = p
		i = i - 1
	}

	json, err := MarshalJSON(prev)
	if err != nil {
		t.Fatal(err)
	}

	var gotErr *Error
	err = json2.Unmarshal(json, &gotErr)
	if err != nil {
		t.Fatal(err)
	}
	if !gotErr.Equal(prev) {
		t.Fatal("Errors do not match")
	}
}
