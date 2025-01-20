package connhandler

import (
	"fmt"
	"testing"
)

func Test_Dummy(t *testing.T) {
	b := make([]byte, 0, 512)
	fmt.Println(len(b[0:100]), cap(b[0:100]))
	fmt.Println(len(b[0:512]), cap(b[0:512]))
	fmt.Println(len(b[100:200]), cap(b[100:200]))
	fmt.Println(len(b[:200]), cap(b[:200]))
}
