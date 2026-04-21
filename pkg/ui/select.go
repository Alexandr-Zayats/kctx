package ui

import (
	"fmt"

	"kctx/internal/model"
	"kctx/internal/provider"
)

func SelectString(title string, items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("%s list is empty", title)
	}

	return FzfSelect(title, items)
}

func SelectProvider(items []provider.Provider) (provider.Provider, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("provider list is empty")
	}

	names := make([]string, 0, len(items))
	index := make(map[string]provider.Provider, len(items))

	for _, p := range items {
		name := p.Name()
		names = append(names, name)
		index[name] = p
	}

	selected, err := FzfSelect("Provider", names)
	if err != nil {
		return nil, err
	}

	p, ok := index[selected]
	if !ok {
		return nil, fmt.Errorf("selected provider not found: %s", selected)
	}

	return p, nil
}

func SelectAccount(items []model.Account) (model.Account, error) {
	if len(items) == 0 {
		return model.Account{}, fmt.Errorf("account list is empty")
	}

	names := make([]string, 0, len(items))
	index := make(map[string]model.Account, len(items))

	for _, a := range items {
		name := a.Name
		names = append(names, name)
		index[name] = a
	}

	selected, err := FzfSelect("Account", names)
	if err != nil {
		return model.Account{}, err
	}

	acc, ok := index[selected]
	if !ok {
		return model.Account{}, fmt.Errorf("selected account not found: %s", selected)
	}

	return acc, nil
}

func SelectCluster(items []model.Cluster) (model.Cluster, error) {
	if len(items) == 0 {
		return model.Cluster{}, fmt.Errorf("cluster list is empty")
	}

	names := make([]string, 0, len(items))
	index := make(map[string]model.Cluster, len(items))

	for _, c := range items {
		label := c.Name
		if c.Location != "" {
			label = fmt.Sprintf("%s (%s)", c.Name, c.Location)
		}

		names = append(names, label)
		index[label] = c
	}

	selected, err := FzfSelect("Cluster", names)
	if err != nil {
		return model.Cluster{}, err
	}

	cluster, ok := index[selected]
	if !ok {
		return model.Cluster{}, fmt.Errorf("selected cluster not found: %s", selected)
	}

	return cluster, nil
}
