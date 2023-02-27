package presenter

import "github.com/gofiber/fiber/v2"

// FileObject is the presenter object which will be passed in the response by Handler
type FileObject struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
}

func FileSuccessResponse(data *[]FileObject) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

func FileErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   nil,
		"error":  err.Error(),
	}
}
