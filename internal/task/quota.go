package task

import (
	"context"
	"time"

	"github.com/bestruirui/octopus/internal/op"
	"github.com/bestruirui/octopus/internal/utils/log"
)

func CheckAndResetQuotas() {
	ctx := context.Background()
	keys, err := op.APIKeyList(ctx)
	if err != nil {
		log.Errorf("failed to list api keys for quota reset: %v", err)
		return
	}

	now := time.Now().Unix()
	for _, key := range keys {
		if key.AutoResetQuota && key.ResetDuration > 0 {
			if key.NextResetTime == 0 {
				key.NextResetTime = now + key.ResetDuration
				op.APIKeyUpdate(&key, ctx)
			} else if now >= key.NextResetTime {
				if err := op.StatsAPIKeyReset(key.ID); err == nil {
					key.NextResetTime = now + key.ResetDuration
					if err := op.APIKeyUpdate(&key, ctx); err != nil {
						log.Errorf("failed to update api key next reset time for key %s: %v", key.Name, err)
					} else {
						log.Infof("reset quota for api key %s (id: %d)", key.Name, key.ID)
					}
				} else {
					log.Errorf("failed to reset stats for api key %s: %v", key.Name, err)
				}
			}
		}
	}
}
