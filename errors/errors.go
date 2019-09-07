package errors

import (
	"encoding/json"
	"fmt"
)

type Op string
type Kind uint8

const (
	Other       Kind = iota // Unclassified error.
	Invalid                 // Invalid operation for this type of item.
	IO                      // External I/O error such as network failure.
	IsDir                   // Item is a directory.
	NotDir                  // Item is not a directory.
	Exists                  // Item exists
	NotExists               // Item does not exist
	NotEmpty                // Directory not empty.
	Private                 // Information withheld.
	Internal                // Internal error or inconsistency.
	Security                // Security check failed
	Permissions             // Not enough permissions
	Parse                   // Failed to parse
)

func (k Kind) HumanReadable() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case IO:
		return "I/O error"
	case IsDir:
		return "item is a directory"
	case NotDir:
		return "item is not a directory"
	case Exists:
		return "item exists"
	case NotExists:
		return "item does not exist"
	case NotEmpty:
		return "directory not empty"
	case Private:
		return "information withheld"
	case Internal:
		return "internal error"
	case Security:
		return "Security check failed"
	case Permissions:
		return "not enough permissions"
	case Parse:
		return "failed to parse entity"
	}
	return "unknown error kind"
}

type Pair struct {
	Key, Value interface{}
}

func (p Pair) Equal(other Pair) bool {
	if p.Key != other.Key {
		return false
	}

	return p.Value == other.Value
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

func (e Error) Unwrap() error {
	if e.CausedBy != nil {
		return e.CausedBy
	}

	if e.err != nil {
		return e.err
	}

	return fmt.Errorf(e.SErr)
}

func (e *Error) Error() string {
	j, err := MarshalJSON(e)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (e *Error) Equal(other *Error) bool {
	if e.Op != other.Op {
		return false
	}

	if uint8(e.Kind) != uint8(other.Kind) {
		return false
	}

	if len(e.Args) != len(other.Args) {
		return false
	}
	for i, v := range e.Args {
		ov := other.Args[i]
		if !v.Equal(ov) {
			return false
		}
	}

	if (e.CausedBy == nil && other.CausedBy != nil) || (e.CausedBy != nil && other.CausedBy == nil) {
		return false
	}
	if e.CausedBy != nil {
		return e.CausedBy.Equal(other.CausedBy)
	}

	return true
}

func MarshalJSON(e *Error) ([]byte, error) {
	if e.CausedBy != nil {
		e.err = nil
		e.SErr = ""
	}
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
		kind.HumanReadable(),
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
