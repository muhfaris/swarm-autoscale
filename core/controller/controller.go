package controller

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2/log"
	"github.com/muhfaris/swarm-autoscale/core/models"
	"github.com/muhfaris/swarm-autoscale/pkg/cache"
)

func UpdateService(ctx context.Context, params models.ReqContainer) error {
	var dockerPath = "/usr/bin/docker"

	containerRaw, exist := cache.Get(params.ServiceID)
	if exist {
		container, ok := containerRaw.(models.ReqContainer)
		if !ok {
			return fmt.Errorf("error parsing container from cache")
		}

		container.Based = params.Based
		container.Current.CPUPercentage = container.Current.CPUPercentage + params.Current.CPUPercentage
		container.Current.MemoryPercentage = container.Current.MemoryPercentage + params.Current.MemoryPercentage
		desiredReplica := container.Scale()
		if desiredReplica > 0 {
			scaleUpCommand := fmt.Sprintf("%s=%d", container.ServiceName, int(desiredReplica))
			commands := []string{"service", "scale", scaleUpCommand}
			err := exec.Command(dockerPath, commands...).Run()
			if err != nil {
				log.Error("error existing scale up: ", err)
				return err
			}
		}

		cache.Set(params.ServiceName, params)
		return nil
	}

	desiredReplica := params.Scale()
	if desiredReplica > 0 {
		scaleUpCommand := fmt.Sprintf("%s=%d", params.ServiceName, int(desiredReplica))
		commands := []string{"service", "scale", scaleUpCommand}
		err := exec.Command(dockerPath, commands...).Run()
		if err != nil {
			log.Error("error scale up", err)
			return err
		}
	}

	return nil
}
