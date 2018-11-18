package controllers

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/ms304w-client/basis/errors"
)

type FileController struct {
	BaseController
}

// 添加
func (c *FileController) AddMaterial() {
	f, h, err := c.GetFile("image")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}
	defer f.Close()

	fileName := h.Filename

	arr := strings.Split(fileName, ":")
	if len(arr) > 1 {
		fileName = arr[len(arr)-1]
	}

	// dir
	dir := beego.AppConfig.String("img_material_dir")

	fileName = fmt.Sprintf("%d%s", time.Now().Unix(), path.Ext(fileName))

	// newFileName := fmt.Sprintf("%s%s", dir, fileName)

	// 保存
	if err := c.SaveToFile("image", path.Join(dir, fileName)); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, struct {
		Img string `json:"img"`
	}{
		Img: fileName,
	}, nil)

	return
}
