package customercontroller

type createCustomCausesRequest struct {
	Name string `json:"name" binding:"required"`
}
