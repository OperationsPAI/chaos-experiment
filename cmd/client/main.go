package main

import (
	"context"

	"github.com/OperationsPAI/chaos-experiment/client"
	"github.com/k0kubun/pp/v3"
)

func main() {
	list, _ := client.GetContainersWithAppLabel(context.Background(), "ts")
	pp.Print(list)
}
