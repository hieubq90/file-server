package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"net/http"
)

func GetFileContentType(out *multipart.FileHeader) (string, error) {
	file, err := out.Open()
	if err != nil {
		return "", err
	}
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// UploadFile godoc
//
// @Summary Upload file
// @Description Allow upload single or multiple file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param project path string true "Project Name"
// @Param folder path string true "Folder Name"
// @param files formData file true "Files"
// @Success 200 {object} presenter.ResponseHTTP{data=[]presenter.FileObject}
// @Failure 400 {object} presenter.ResponseHTTP{}
// @Failure 401 {object} presenter.ResponseHTTP{}
// @Failure 503 {object} presenter.ResponseHTTP{}
// @Router /{project}/{folder}/files [post]
func UploadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the multipart form:
		form, err := c.MultipartForm() // => *multipart.Form
		if err != nil {
			return err
		}

		// Get all files from "documents" key:
		files := form.File["files"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			contentType, _ := GetFileContentType(file)
			fmt.Println(file.Filename, file.Size, contentType)
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to minio:
			//err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))

			// Check for errors
			if err != nil {
				return err
			}
		}
		return nil
	}
}
