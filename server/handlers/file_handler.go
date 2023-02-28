package handlers

import (
	"errors"
	"file-server/mio"
	"file-server/server/presenter"
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
// @Failure 500 {object} presenter.ResponseHTTP{}
// @Router /{project}/{folder}/files [post]
func UploadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the multipart form:
		form, err := c.MultipartForm()
		// => *multipart.Form
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.FileErrorResponse(err))
		}

		// Get all files from "documents" key:
		files := form.File["files"]
		// => []*multipart.FileHeader

		// Get project & folder
		params := c.AllParams()

		// Loop through files:
		data := make([]presenter.FileObject, 0)
		for _, file := range files {
			contentType, _ := GetFileContentType(file)
			uf, err := mio.UploadFile(params["project"], params["folder"], file, contentType)
			if err == nil && uf != nil {
				data = append(data, *uf)
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.FileErrorResponse(err))
		}

		return c.JSON(presenter.FileUploadSuccessResponse(&data))
	}
}

// DownloadFile godoc
//
// @Summary Download file
// @Description Allow download
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param project path string true "Project Name"
// @Param folder path string true "Folder Name"
// @param filename path string true "File Name"
// @Success 301 {string} string
// @Failure 400 {object} presenter.ResponseHTTP{}
// @Failure 401 {object} presenter.ResponseHTTP{}
// @Failure 503 {object} presenter.ResponseHTTP{}
// @Router /{project}/{folder}/{filename} [get]
func DownloadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get project & folder
		params := c.AllParams()
		if params["project"] == "" || params["folder"] == "" || params["filename"] == "" {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.FileErrorResponse(errors.New("Tham số không hợp lệ")))
		}
		fmt.Println("DEBUG", params["filename"])
		link, err := mio.GenerateDownloadLink(params["project"], params["folder"], params["filename"])

		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenter.FileErrorResponse(err))
		}

		return c.Redirect(link.String(), http.StatusMovedPermanently)
	}
}
