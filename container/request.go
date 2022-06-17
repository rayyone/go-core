package corecontainer

import (
	"context"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/rayyone/go-core/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rayyone/go-core/helpers/array"
	"github.com/rayyone/go-core/helpers/method"
	"github.com/rayyone/go-core/helpers/pagination"
	"github.com/rayyone/go-core/ryerr"
)

type Auth struct{}
type ExtraData struct{}

type UrlParams struct {
	Str map[string]string
	Arr map[string][]string
}

type RequestInf interface {
	GetDBM() *Database
	SetPostParams(params interface{}) error
	ValidateFileType(file *multipart.FileHeader, allowTypes []string) error
}

type Request struct {
	Auth
	ExtraData
	Ctx         context.Context
	DBM         *Database
	GinCtx      *gin.Context
	Pagination  pagination.Config
	UrlParams   UrlParams
	PostParams  interface{}
	QueryParams interface{}
}

func (r *Request) GetDBM() *Database {
	return r.DBM
}

func (r *Request) SetQueryParams(params interface{}) error {
	if params == nil {
		return nil
	}
	err := r.GinCtx.ShouldBindQuery(params)
	if err != nil {
		_ = r.GinCtx.Error(err)
		return err
	}

	r.QueryParams = params
	return nil
}

func (r *Request) SetPostParams(params interface{}) error {
	if params == nil {
		return nil
	}
	b := binding.Default(r.GinCtx.Request.Method, r.GinCtx.ContentType())
	var i interface{} = b
	var err error
	bBody, ok := i.(binding.BindingBody)
	if ok {
		// Use ShouldBindBodyWith so we can reuse request body after we read it (so we can have multiple binding)
		err = r.GinCtx.ShouldBindBodyWith(params, bBody)
	} else {
		err = r.GinCtx.ShouldBind(params)
	}

	if err != nil {
		_ = r.GinCtx.Error(err)
		return err
	}

	r.PostParams = params
	return nil
}

func (r *Request) ValidateFileType(file *multipart.FileHeader, allowTypes []string) error {
	if file == nil {
		return nil
	}

	var mimes []string
	for _, allowType := range allowTypes {
		mimes = append(mimes, method.ExtractMimeFromType(allowType)...)
	}
	fileType := file.Header.Get("Content-Type")
	if exists := array.InArray(fileType, mimes); exists {
		return nil
	}

	return ryerr.Validation.Newf("File is only allow %s type", strings.Join(allowTypes, " / "))
}

func InitCoreRequest(c *gin.Context) *Request {
	var r Request
	r.GinCtx = c
	r.Ctx = context.Background()
	r.DBM = NewCoreDBManager(database.GetDB())
	initUrlParams(c, &r)
	initPagination(c, &r)
	return &r
}

func initUrlParams(c *gin.Context, r *Request) {
	r.UrlParams.Str = make(map[string]string)
	r.UrlParams.Arr = make(map[string][]string)

	if c != nil {
		params := c.Request.URL.Query()
		for param, value := range params {
			if len(value) <= 1 {
				r.UrlParams.Str[param] = value[0]
			}
			r.UrlParams.Arr[param] = value
		}
	}
}

func initPagination(c *gin.Context, r *Request) {
	var page, limit int
	var err error
	if c != nil {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil || page == 0 {
			page = 1
		}

		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 25
		}
	}

	r.Pagination = pagination.GetPaginationConfig(page, limit)
}
