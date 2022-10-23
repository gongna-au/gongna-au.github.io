package posts

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	re := Random("192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4")
	fmt.Println(re)
}
