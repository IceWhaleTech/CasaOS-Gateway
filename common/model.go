package common

type Route struct {
	Path   string `json:"path" binding:"required"`
	Target string `json:"target" binding:"required"`
}

type UpdatePortRequest struct {
	Port string `json:"port" binding:"required"`
}
