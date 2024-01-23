package vault

import "strings"

// func (v *Vault) AllNamespaces() ([]string, error) {

// }

func (v *Vault) ChangeNamespace(ns string) []string {
	v.Logger.Debug().Msgf("Changing namespace to: %v", ns)
	v.vault.SetNamespace(ns)
	list, err := v.ListNamespaces()
	if err != nil {
		v.Logger.Err(err).Msgf("Failed to list namespaces: %v", err)
	}
	v.Logger.Debug().Msgf("New available namespaces are: %v", list)
	return list
}

func (v *Vault) SetNamespace(ns string) {
	v.Logger.Debug().Msgf("Changing namespace to: %v", ns)
	v.vault.SetNamespace(ns)
}

func (v *Vault) ListNamespaces() ([]string, error) {
	namespaces := []string{}
	secret, err := v.vault.Logical().List("sys/namespaces")
	if err != nil {
		return nil, err
	}

	if secret != nil {
		keys, ok := secret.Data["keys"].([]interface{})
		if !ok || len(keys) == 0 {
			return namespaces, nil
		}
		for _, namespace := range keys {
			trimmedNamespace := strings.TrimSuffix(namespace.(string), "/")
			namespaces = append(namespaces, trimmedNamespace)
		}
	} else {
		return namespaces, nil
	}
	return namespaces, nil
}
