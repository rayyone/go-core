package imagehelper

import (
	"github.com/rayyone/go-core/errors"

	"gopkg.in/h2non/bimg.v1"
)

type ImageHelper struct {
	File     *[]byte
	BImgFile *bimg.Image
	Error    error
}

type ImageSize struct {
	Width  int
	Height int
}

func NewImageHelper(file *[]byte) *ImageHelper {
	bImgFile := bimg.NewImage(*file)
	return &ImageHelper{File: file, BImgFile: bImgFile}
}

func (i *ImageHelper) GetImage() ([]byte, error) {
	imageData := i.BImgFile.Image()
	return imageData, i.Error
}

func (ih *ImageHelper) Resize(width int, height int, options ...interface{}) ([]byte, error) {
	newImage, err := ih.BImgFile.Resize(width, height)
	if err != nil {
		return nil, errors.BadRequest.Newf("Cannot resize image. Error: %v", err)
	}
	return newImage, err
}

func (ih *ImageHelper) ResizeByHeight(height int) *ImageHelper {
	options := bimg.Options{
		Height:  height,
		Enlarge: true,
	}
	bs, err := ih.BImgFile.Process(options)
	if err != nil {
		ih.Error = errors.BadRequest.Newf("Cannot resize image. Error: %v", err)
		return nil
	}
	ih.BImgFile = bimg.NewImage(bs)
	return ih
}

func (ih *ImageHelper) ResizeByWidth(width int) *ImageHelper {
	options := bimg.Options{
		Width:   width,
		Enlarge: true,
	}
	bs, err := ih.BImgFile.Process(options)
	if err != nil {
		ih.Error = errors.BadRequest.Newf("Cannot resize image. Error: %v", err)
		return nil
	}
	ih.BImgFile = bimg.NewImage(bs)
	return ih
}

func (ih *ImageHelper) Compress(quality int) *ImageHelper {
	options := bimg.Options{Quality: quality}
	bs, err := ih.BImgFile.Process(options)
	if err != nil {
		ih.Error = errors.BadRequest.Newf("Cannot resize image. Error: %v", err)
		return nil
	}
	ih.BImgFile = bimg.NewImage(bs)
	return ih
}

func (ih *ImageHelper) GetSize(bImage *[]byte) (*ImageSize, error) {
	bSize, err := bimg.Size(*bImage)
	if err != nil {
		return nil, errors.BadRequest.Newf("Cannot get image size. Error: %v", err)
	}
	return &ImageSize{Width: bSize.Width, Height: bSize.Height}, nil
}

func (ih *ImageHelper) MakeThumbnail(width int, height int, options ...interface{}) ([]byte, error) {
	bs, err := ih.Resize(width, height, options)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (ih *ImageHelper) MakeSquareThumbnail(size int) ([]byte, error) {
	options := bimg.Options{
		Width:   size,
		Height:  size,
		Crop:    true,
		Quality: 80,
		Gravity: bimg.GravitySmart,
	}
	bs, err := ih.BImgFile.Process(options)
	if err != nil {
		return nil, errors.BadRequest.Newf("Cannot make square thumbnail. Error: %v", err)
	}

	return bs, nil
}

func (ih *ImageHelper) FileType() string {
	return ih.BImgFile.Type()
}
