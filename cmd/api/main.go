package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sohan-14/auth-system-go/internal/app"
)

func main() {
	if err := app.Run(context.Background(), os.Getenv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
