package model

type Account struct {
	Name string
	Meta map[string]string
}

type Cluster struct {
	Name     string
	Location string
	Meta     map[string]string
}
