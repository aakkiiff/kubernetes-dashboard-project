package models

type Pod struct {
	Name          string `json:"name"`
	NamespaceName string `json:"namespacename"`
	Image         string `json:"image"`
}
