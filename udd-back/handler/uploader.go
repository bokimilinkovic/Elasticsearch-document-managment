package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Uploader struct {
	Dir string
}

func (u *Uploader) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// Destination
	dst, err := os.Create(u.Dir + file.Filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
		return err
	}

	return c.String(http.StatusOK, "Uploaded successfully")
}
