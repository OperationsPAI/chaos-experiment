package main

import (
	"context"

	"github.com/OperationsPAI/chaos-experiment/chaos"
	"github.com/OperationsPAI/chaos-experiment/client"
	"github.com/OperationsPAI/chaos-experiment/controllers"
	"k8s.io/utils/pointer"
)

func main() {
	ctx := context.Background()
	client := client.GetK8sClient()
	controllers.CreateJVMRuntimeMutatorChaos(client, ctx, "ts", "ts-execute-service", "execute.serivce.ExecuteServiceImpl", "getOrderByIdFromOrder", "string", pointer.String("5m"), nil, nil, chaos.WithRuntimeMutatorStrategy("empty"))
}
