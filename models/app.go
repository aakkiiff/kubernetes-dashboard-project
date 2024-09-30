package models

type App struct {
	Name          string `json:"name" binding:"required"`
	Image         string `json:"image" binding:"required"`
	ContainerPort int    `json:"containerport" binding:"required"`
}
