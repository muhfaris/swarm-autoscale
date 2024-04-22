package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/muhfaris/swarm-autoscale/cmd"
)

const (
	dockerVersion     = "1.41"
	dockerHost        = "unix:///var/run/docker.sock"
	labelAutoScale    = "swarm.autoscale=true"
	labelAutoScaleMin = "swarm.autoscale.min"
	labelAutoScaleMax = "swarm.autoscale.max"
	labelAutoScaleCPU = "swarm.autoscale.cpu"
	labelAutoScaleMem = "swarm.autoscale.mem"
)

type ServiceConfig struct {
	MinReplicas uint64
	MaxReplicas uint64
	CPU         float64
	Memory      float64
}

func main() {
	cmd.HTTPServe()
}

func main2() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	cli, err := client.NewClientWithOpts(client.WithHost(dockerHost), client.WithVersion(dockerVersion))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				err = watchServices(context.Background(), cli)
				if err != nil {
					return
				}
			}
		}
	}()

	<-done
}

func watchServices(ctx context.Context, cli *client.Client) error {
	services, err := cli.ServiceList(ctx, types.ServiceListOptions{
		Status:  true,
		Filters: filters.NewArgs(filters.Arg("label", labelAutoScale))})
	if err != nil {
		return err
	}

	for _, service := range services {
		if service.ServiceStatus.RunningTasks == 0 {
			continue
		}

		svcConfig := getLabelsService(service)
		go scaleService(ctx, svcConfig, cli, service)
	}

	return nil
}

func getLabelsService(service swarm.Service) ServiceConfig {
	var svcConfig ServiceConfig
	for key, value := range service.Spec.Labels {
		switch key {
		case labelAutoScaleMin:
			// validate input arg scale number
			scale, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				fmt.Println("error parsing min scale")
				continue
			}
			svcConfig.MinReplicas = scale

		case labelAutoScaleMax:
			// validate input arg scale number
			scale, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				fmt.Println("error parsing max scale")
				continue
			}
			svcConfig.MaxReplicas = scale

		case labelAutoScaleCPU:
			cpu, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("error parsing cpu scale")
				continue
			}
			svcConfig.CPU = cpu

		case labelAutoScaleMem:
			mem, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("error parsing mem scale")
				continue
			}
			svcConfig.Memory = mem
		}
	}

	return svcConfig
}

func scaleService(ctx context.Context, svcConfig ServiceConfig, cli *client.Client, service swarm.Service) error {
	var (
		currentReplicas = uint64(*service.Spec.Mode.Replicated.Replicas)
	)

	if currentReplicas < svcConfig.MinReplicas {

	}
	return nil
}

func runScaleService(ctx context.Context, cli *client.Client, service swarm.Service) error {
	scale := uint64(1)
	specMode := &service.Spec.Mode
	switch {
	case specMode.Replicated != nil:
		specMode.Replicated.Replicas = &scale

	case specMode.ReplicatedJob != nil:
		specMode.ReplicatedJob.TotalCompletions = &scale
	}

	response, err := cli.ServiceUpdate(ctx, service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("response: %v\n", response)
	return nil
}
