package gourd

import (
	"fmt"
	"github.com/go-gourd/gourd/event"
	"testing"
)

func TestEvent(t *testing.T) {

	event.Listen("test.add", func(params any) {
		fmt.Println("test.add", params)
	})

	event.Listen("test.edit", func(params any) {
		fmt.Println("test.edit", params)
	})

	event.Listen("user.add", func(params any) {
		fmt.Println("user.add", params)
	})

	event.Listen("user.edit", func(params any) {
		fmt.Println("user.edit", params)
	})

	event.Listen("user.address.add", func(params any) {
		fmt.Println("user.address.add", params)
	})

	fmt.Println("---------test.*-----------")
	event.Trigger("test.*", 1)

	fmt.Println("---------*.add-----------")
	event.Trigger("*.add", 2)

	fmt.Println("---------user.address.add-----------")
	event.Trigger("user.address.add", 3)

	fmt.Println("---------user.*.add-----------")
	event.Trigger("user.*.add", 4)
}
