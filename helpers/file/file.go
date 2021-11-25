package filehelper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/rayyone/go-core/errors"
	"github.com/rayyone/go-core/helpers/str"
)

// GetMimeType return file mime type
func GetMimeType(filename string) (mimeType string) {
	fileExt := filepath.Ext(filename)

	return mime.TypeByExtension(fileExt)
}

// RandomFilename return randomize filename
func RandomFilename(file *multipart.FileHeader) string {
	return fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), str.Random(32), filepath.Ext(file.Filename))
}

func ReadFormFile(file *multipart.FileHeader) (*[]byte, error) {
	imageFile, err := file.Open()
	if err != nil {
		return nil, errors.BadRequest.Newf("Cannot open file. Error: %v", err)
	}
	fileByte, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return nil, errors.BadRequest.Newf("Cannot read file. Error: %v", err)
	}
	return &fileByte, nil
}

func GetFileMd5(file multipart.File) (md5Str string, err error) {
	h := md5.New()
	if _, err := file.Seek(0, 0); err != nil {
		return "", errors.BadRequest.Newf("Get file md5 error: %v", err)
	}
	if _, err := io.Copy(h, file); err != nil {
		return "", errors.BadRequest.Newf("Get file md5 error: %v", err)
	}
	md5Str = hex.EncodeToString(h.Sum(nil))
	// Call file.Seek to allow re-reading the file stream again by set it to the start point
	if _, err := file.Seek(0, 0); err != nil {
		return "", errors.BadRequest.Newf("Get file md5 error: %v", err)
	}
	return md5Str, nil
}
