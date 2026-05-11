package vault

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (v *Vault) TokenInfo() (ttl int64, policies []string, err error) {
	secret, err := v.vault.Auth().Token().LookupSelf()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to lookup token: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return 0, nil, fmt.Errorf("empty token lookup response")
	}

	if raw, ok := secret.Data["ttl"]; ok {
		switch n := raw.(type) {
		case float64:
			ttl = int64(n)
		case int64:
			ttl = n
		case json.Number:
			ttl, _ = n.Int64()
		}
	}

	if raw, ok := secret.Data["policies"]; ok {
		if list, ok := raw.([]interface{}); ok {
			for _, p := range list {
				if s, ok := p.(string); ok {
					policies = append(policies, s)
				}
			}
		}
	}

	return ttl, policies, nil
}

func (v *Vault) SealStatus() (status string, clusterName string, err error) {
	health, err := v.Sys.Health()
	if err != nil {
		return "Unknown", "", fmt.Errorf("failed to get health: %w", err)
	}

	switch {
	case health.Sealed:
		status = "Sealed"
	case strings.Contains(strings.ToLower(health.ReplicationPerformanceMode), "standby"),
		strings.Contains(strings.ToLower(health.ReplicationDRMode), "standby"):
		status = "Standby"
	default:
		status = "Active"
	}

	return status, health.ClusterName, nil
}
