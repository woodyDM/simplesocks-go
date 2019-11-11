package pro

import (
	"fmt"
	"testing"
	"time"
)

func TestPanic(t *testing.T) {
	i := 0
	for {
		i++
		time.Sleep(1 * time.Second)

		if i > 5 {

			fmt.Println("After panic")
		}

	}

	fmt.Println("impossible, After for")

}

func TestSlice(t *testing.T) {
	f := []int{1, 2, 3, 4, 5}
	c := f[1:4]
	t.Log(len(c))
}

func Test_when_byte_overflow(t *testing.T) {
	var a byte = 0
	var b byte = 1
	if a-b != 255 {
		t.Fatal("255")
	}
	a = 255
	b = 1
	t.Log(a + b)
	if a+b != 0 {
		t.Fatal(0)
	}
}

//
