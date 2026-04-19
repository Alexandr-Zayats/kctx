package ui

import (
	"kctx/internal/model"
	"kctx/internal/provider"
)

func SelectProvider(ps []provider.Provider) (provider.Provider, error) {
	var items []string
	for _, p := range ps {
		items = append(items, p.Name())
	}

	selected, err := Select("Provider", items)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.Name() == selected {
			return p, nil
		}
	}

	return nil, err
}

func SelectAccount(a []model.Account) (model.Account, error) {
	var items []string
	for _, v := range a {
		items = append(items, v.Name)
	}

	selected, err := Select("Account", items)
	if err != nil {
		return model.Account{}, err
	}

	for _, v := range a {
		if v.Name == selected {
			return v, nil
		}
	}

	return model.Account{}, err
}

func SelectCluster(c []model.Cluster) (model.Cluster, error) {
	var items []string
	for _, v := range c {
		items = append(items, v.Name)
	}

	selected, err := Select("Cluster", items)
	if err != nil {
		return model.Cluster{}, err
	}

	for _, v := range c {
		if v.Name == selected {
			return v, nil
		}
	}

	return model.Cluster{}, err
}

func SelectString(label string, items []string) (string, error) {
	return Select(label, items)
}
