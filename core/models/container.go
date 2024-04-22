package models

import (
	"math"

	"github.com/gofiber/fiber/v2/log"
)

type ReqContainer struct {
	ServiceID   string          `json:"service_id"`
	ServiceName string          `json:"service_name"`
	Current     ReqCurrentStats `json:"current"`
	Based       ReqBasedStats   `json:"based"`
}

func (r *ReqContainer) Scale() (desiredReplicas float64) {
	var (
		currentCPU = r.Current.CPUPercentage
		currentMem = r.Current.MemoryPercentage
		basedCPU   = r.Based.CPUPercentage
		basedMem   = r.Based.MemoryPercentage
		replicas   = float64(r.Current.Replicas)
		minReplica = float64(r.Based.Min)
		maxReplica = float64(r.Based.Max)
	)

	defer func() {
		if replicas == desiredReplicas {
			// dont need update service
			desiredReplicas = 0
			return
		}

		if desiredReplicas < minReplica {
			desiredReplicas = 0
			return
		}

		if desiredReplicas < replicas {
			return
		}

		if replicas >= maxReplica {
			desiredReplicas = 0
			log.Warnf("replicas over max (%d)", int(maxReplica))
			return
		}
	}()

	// desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
	if basedCPU > 0 {
		desiredReplicas = math.Ceil(replicas * (currentCPU / basedCPU))
		return
	}

	if basedMem > 0 {
		desiredReplicas = math.Ceil(replicas * (currentMem / basedMem))
		return
	}

	return 0
}

type ReqBasedStats struct {
	CPUPercentage    float64 `json:"cpu_percentage"`
	MemoryPercentage float64 `json:"memory_percentage"`
	Min              int64   `json:"min"`
	Max              int64   `json:"max"`
}

type ReqCurrentStats struct {
	CPUPercentage    float64 `json:"cpu_percentage"`
	MemoryPercentage float64 `json:"memory_percentage"`
	Replicas         int64   `json:"replicas"`
}
