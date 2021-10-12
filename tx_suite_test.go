package tx_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tx Suite")
}

// JSONFuzz is the same as the Fuzz testing function exposed by surge, but it
// uses JSON.
func JSONFuzz(t reflect.Type) {
	// Fuzz data
	data, ok := quick.Value(reflect.TypeOf([]byte{}), rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		panic(fmt.Errorf("cannot generate value of type %v", t))
	}
	// Unmarshal
	x := reflect.New(t)
	if err := json.Unmarshal(data.Bytes(), x.Interface()); err != nil {
		// Ignore the error, because we are only interested in whether or not
		// the unmarshaling causes a panic.
	}
}

// JSONMarshalUnmarshalCheck is the same as the MarshalUnmarshalCheck testing
// function exposed by surge, but it uses JSON.
func JSONMarshalUnmarshalCheck(t reflect.Type) error {
	// Generate
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal
	data, err := json.Marshal(x.Interface())
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}
	// Unmarshal
	y := reflect.New(t)
	if err := json.Unmarshal(data, y.Interface()); err != nil {
		return fmt.Errorf("cannot unmarshal: %v", err)
	}
	// Equality
	if !reflect.DeepEqual(x.Interface(), y.Elem().Interface()) {
		return fmt.Errorf("unequal")
	}
	return nil
}
