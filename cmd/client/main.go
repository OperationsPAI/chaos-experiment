package main

import (
	"context"

	"github.com/LGU-SE-Internal/chaos-experiment/client"
	"github.com/k0kubun/pp/v3"
)

func main() {
	ctx := context.Background()
	list, _ := client.GetContainersWithAppLabel(ctx, "ts")
	pp.Print(list)
}
