package vault

func (v *Vault) AllPolicies() ([]string, error) {
	pl, err := v.Sys.ListPolicies()
	if err != nil {
		return nil, err
	}

	policies := []string{}
	for _, p := range pl {
		policies = append(policies, p)
	}

	return policies, nil
}

func (v *Vault) GetPolicyInfo(name string) (string, error) {
	//TODO: Might need to make it a custom function
	policy, err := v.Sys.GetPolicy(name)
	if err != nil {
		return "", err
	}
	return policy, nil
}
