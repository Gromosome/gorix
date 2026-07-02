package controller

import (
	"io"
	"os"

	"github.com/Gromosome/gorix/gorix"
)

type MediaController struct{}

func NewMediaController() *MediaController {
	return &MediaController{}
}

func (c *MediaController) UploadRaw() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.POST, "/raw", func(ctx *gorix.Context) (any, error) {
		var savedPath string

		return ctx.
			Status(gorix.StatusCreated).
			LimitBody(20<<20). // 20 MB
			StreamFile("file", func(file gorix.FileStream) error {
				dstPath := "../tmp/" + file.FileName

				dst, err := os.Create(dstPath)
				if err != nil {
					return err
				}
				defer dst.Close()

				_, err = io.Copy(dst, file.Reader)
				if err != nil {
					return err
				}

				savedPath = dstPath
				return nil
			}).
			ResponseEntityJSON(func() (any, error) {
				return map[string]any{
					"message": "file uploaded",
					"path":    savedPath,
				}, nil
			})
	}
}
