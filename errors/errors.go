package errors

import (
	"encoding/json"
	"fmt"
	"log"
)

type Op string
type Kind uint8

const (
	Other         Kind = iota // Unclassified error.
	Invalid                   // Invalid operation for this type of item.
	Permission                // Permission denied.
	IO                        // External I/O error such as network failure.
	Exist                     // Item already exists.
	NotExist                  // Item does not exist.
	IsDir                     // Item is a directory.
	NotDir                    // Item is not a directory.
	NotEmpty                  // Directory not empty.
	Private                   // Information withheld.
	Internal                  // Internal error or inconsistency.
	CannotDecrypt             // No wrapped key for user with read access.
	Transient                 // A transient error.
	BrokenLink                // Link target does not exist.
	Security                  // Security check failed
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case Permission:
		return "permission denied"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case BrokenLink:
		return "link target does not exist"
	case IsDir:
		return "item is a directory"
	case NotDir:
		return "item is not a directory"
	case NotEmpty:
		return "directory not empty"
	case Private:
		return "information withheld"
	case Internal:
		return "internal error"
	case CannotDecrypt:
		return `no wrapped key for user; owner must "upspin share -fix"`
	case Transient:
		return "transient error"
	case Security:
		return "Security check failed"
	}
	return "unknown error kind"
}

type Pair struct {
	Key, Value interface{}
}

const ContextPair string = "context"

func P(key, value interface{}) Pair {
	return Pair{
		key,
		value,
	}
}

// C constructs a context Pair
func C(context string) Pair {
	return Pair{
		Key:   ContextPair,
		Value: context,
	}
}

type Error struct {
	Op       Op     `json:"op"`
	Kind     Kind   `json:"kind"`
	SKind    string `json:"kindString"`
	err      error
	SErr     string `json:"error"`
	Args     []Pair `json:"args"`
	CausedBy *Error `json:"cause"`
}

func (e *Error) Error() string {
	var cause string = ""
	if e.CausedBy != nil {
		cause = e.CausedBy.Error()
	} else if e.SErr != "" {
		cause = e.SErr
	} else if e.err != nil {
		cause = e.err.Error()
	}
	return fmt.Sprintf("%s - %s: %v", e.Kind, e.Op, cause)
}

func MarshalJSON(e *Error) ([]byte, error) {
	if e.err != nil {
		e.SErr = e.err.Error()
	}
	return json.Marshal(e)
}

func Ops(e *Error) []Op {
	res := []Op{e.Op}

	if e.CausedBy == nil {
		return res
	}

	res = append(res, Ops(e.CausedBy)...)

	return res
}

// E constructs a new Error
// It will take an operation Op, which should be unique
// to the method that returns this Error.
// kind can be one of [DEBUG, INFO, WARN, ERROR].
// args can be arbitrary key value pairs
func E(op Op, kind Kind, err error, args ...Pair) error {
	e := &Error{
		op,
		kind,
		kind.String(),
		err,
		"",
		args,
		nil,
	}

	if e.err != nil {
		e.SErr = e.err.Error()
	}

	prev, ok := err.(*Error)
	if ok {
		e.err = nil
		e.CausedBy = prev
	}

	return e
}

func Print(err error) {
	e, ok := err.(*Error)
	if ok {
		j, err := MarshalJSON(e)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(j))
	}
	log.Fatal("unable to mount ramfs: ", err)
}
