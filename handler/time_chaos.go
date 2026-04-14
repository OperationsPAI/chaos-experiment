package handler

import (
	"context"
	"fmt"
	"strconv"

	controllers "github.com/OperationsPAI/chaos-experiment/controllers"
	"k8s.io/utils/pointer"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
)

type TimeSkewSpec struct {
	Duration     int `range:"1-60" description:"Time Unit Minute"`
	System       int `range:"0-0" dynamic:"true"`
	ContainerIdx int `range:"0-0" dynamic:"true" description:"Container Index"`
	TimeOffset   int `range:"-600-600" description:"Time offset in seconds"`
}

func (s *TimeSkewSpec) Create(cli cli.Client, opts ...Option) (string, error) {
	conf := Conf{}
	for _, opt := range opts {
		opt(&conf)
	}

	annotations := make(map[string]string)
	if conf.Annotations != nil {
		annotations = conf.Annotations
	}

	ctx := context.Background()
	if conf.Context != nil {
		ctx = conf.Context
	}

	labels := make(map[string]string)
	if conf.Labels != nil {
		labels = conf.Labels
	}

	ns := conf.Namespace
	system := conf.System

	containerInfo, err := getContainerInfoByIndex(ctx, system, ns, s.ContainerIdx)
	if err != nil {
		return "", err
	}
	appName := containerInfo.AppLabel
	containerName := containerInfo.ContainerName

	duration := pointer.String(strconv.Itoa(s.Duration) + "m")
	// Format the TimeOffset with "s" unit
	timeOffset := fmt.Sprintf("%ds", s.TimeOffset)

	return controllers.CreateTimeChaosWithContainer(cli, ctx, ns, appName, timeOffset, duration, annotations, labels, []string{containerName})
}
