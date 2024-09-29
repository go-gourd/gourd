package test

import (
	"context"
	"fmt"
	"github.com/go-gourd/gourd/event"
	"testing"
)

func TestEvent(t *testing.T) {

	// register event listener
	event.Listen("test.add", func(ctx context.Context) {
		fmt.Println("test.add", ctx)
	})

	event.Listen("test.edit", func(ctx context.Context) {
		fmt.Println("test.edit", ctx)
	})

	event.Listen("user.add", func(ctx context.Context) {
		fmt.Println("user.add", ctx)
	})

	event.Listen("user.edit", func(ctx context.Context) {
		fmt.Println("user.edit", ctx)
	})

	event.Listen("user.address.add", func(ctx context.Context) {
		value := ctx.Value("test_key")
		fmt.Println("user.address.add", value)
	})

	// trigger event
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	fmt.Println("---------test.*-----------")
	event.Trigger("test.*", ctx)

	fmt.Println("---------*.add-----------")
	event.Trigger("*.add", ctx)

	fmt.Println("---------user.address.add-----------")
	event.Trigger("user.address.add", ctx)

	fmt.Println("---------user.*.add-----------")
	event.Trigger("user.*.add", ctx)
}
