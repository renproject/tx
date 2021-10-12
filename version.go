package tx

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/renproject/surge"
)

// Version is an alias of the String ABI type. It adds a small amount of
// semantic meaning to the type.
type Version string

// Enumeration of all transaction versions.
const (
	Version0 = Version("0")
	Version1 = Version("1")
)

func (v Version) String() string {
	switch v {
	case Version0:
		return "0"
	case Version1:
		return "1"
	default:
		return ""
	}
}

// SizeHint returns the number of bytes required to represent the version in
// binary.
func (v Version) SizeHint() int {
	return surge.SizeHintString(v.String())
}

// Marshal the version into binary.
func (v Version) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.MarshalString(v.String(), buf, rem)
}

// Unmarshal the version from binary.
func (v *Version) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.UnmarshalString((*string)(v), buf, rem)
}

// Generate allows us to quickly generate random transaction selectors. This is
// mostly used for writing tests.
func (Version) Generate(r *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(Version(fmt.Sprintf("%v", r.Int()%1)))
}
