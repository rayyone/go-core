package mails

import (
	"fmt"
	"os"
	"path/filepath"
)

type Encoding string

const (
	EncodingBase64 Encoding = "base64"
)

type Attachment struct {
	Name        string
	Content     []byte
	ContentType string
	Encoding    Encoding
	Disposition string
}

func AttachmentFromPath(path string) (*Attachment, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	name := filepath.Base(path)
	return &Attachment{
		Name:        name,
		Content:     content,
		ContentType: detectContentTypeByExt(path),
		Encoding:    EncodingBase64,
		Disposition: "attachment",
	}, nil
}

func AttachmentFromData(name string, content []byte, contentType string) *Attachment {
	return &Attachment{
		Name:        name,
		Content:     content,
		ContentType: contentType,
		Encoding:    EncodingBase64,
		Disposition: "attachment",
	}
}

func (a *Attachment) AsInline() *Attachment {
	a.Disposition = "inline"
	return a
}

func (a *Attachment) AsAttachment() *Attachment {
	a.Disposition = "attachment"
	return a
}

func (a *Attachment) WithContentType(ct string) *Attachment {
	a.ContentType = ct
	return a
}

func (a *Attachment) WithEncoding(enc Encoding) *Attachment {
	a.Encoding = enc
	return a
}

func (a *Attachment) WithName(name string) *Attachment {
	a.Name = name
	return a
}

func detectContentTypeByExt(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
