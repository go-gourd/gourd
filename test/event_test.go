package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-gourd/gourd/event"
)

func TestEvent(t *testing.T) {

	// register event listener with pattern
	event.Listen("test.*", func(ctx context.Context) {
		fmt.Println("test.*", ctx)
	})

	event.Listen("*.add", func(ctx context.Context) {
		fmt.Println("*.add", ctx)
	})

	event.Listen("user.*", func(ctx context.Context) {
		fmt.Println("user.*", ctx)
	})

	event.Listen("user.*.edit", func(ctx context.Context) {
		fmt.Println("user.*.edit", ctx)
	})

	event.Listen("user.address.*", func(ctx context.Context) {
		value := ctx.Value("test_key")
		fmt.Println("user.address.*", value)
	})

	// trigger event with exact name
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	fmt.Println("---------test.add-----------")
	event.Trigger("test.add", ctx)

	fmt.Println("---------user.edit-----------")
	event.Trigger("user.edit", ctx)

	fmt.Println("---------user.address.add-----------")
	event.Trigger("user.address.add", ctx)
	
	fmt.Println("---------user.profile.edit-----------")
	event.Trigger("user.profile.edit", ctx)
}