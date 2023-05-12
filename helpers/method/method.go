package method

import (
	"github.com/rayyone/go-core/ryerr"
	"mime/multipart"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//func EmptyJsonb() postgres.Jsonb {
//	return postgres.Jsonb{RawMessage: []byte("{}")}
//}

func IsSliceEmpty(data interface{}) bool {
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return len(data.([]interface{})) == 0
	}
	reflectElem := reflect.ValueOf(data).Elem()
	if reflectElem.Kind() == reflect.Slice {
		return reflectElem.Len() == 0
	}
	return false
}

func NewBool(val bool) *bool {
	b := val
	return &b
}

func NewInt(val int) *int {
	i := val
	return &i
}

func NewString(val string) *string {
	i := val
	return &i
}

func Ptr[T any](val T) *T {
	return &val
}

func ExtractMimeFromType(typeStr string) []string {
	mimeType := getMimeTypeMap()
	if mimes, ok := mimeType[typeStr]; ok {
		return mimes
	}
	return []string{}
}

func ExtractTypeFromMime(mime string) string {
	mimesTypes := getMimeTypeMap()
	for fileType, mimesType := range mimesTypes {
		for _, mimeType := range mimesType {
			if mimeType == mime {
				return fileType
			}
		}
	}
	return ""
}

func getMimeTypeMap() map[string][]string {
	return map[string][]string{
		"pdf":   {"application/pdf"},
		"image": {"image/jpg", "image/jpeg", "image/png", "image/gif", "image/svg+xml"},
		"video": {"video/x-flv", "video/mp4", "application/x-mpegURL", "video/MP2T", "video/3gpp", "video/quicktime", "video/x-msvideo", "video/x-ms-wmv"},
		"csv": {"text/plain", "application/vnd.ms-excel", "text/x-csv", "text/csv", "application/csv", "application/x-csv", "text/comma-separated-values", "text/x-comma-separated-values", "text/tab-separated-values"},
	}
}

func TraceCaller(skip int) (file string, line int, fnName string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(skip, pc)

	if pc[0] == uintptr(0) {
		return
	}

	f := runtime.FuncForPC(pc[0])
	file, line = f.FileLine(pc[0])
	fnName = f.Name()

	return
}

func GuessTime(t interface{}) (*time.Time, error) {
	var res time.Time
	var err error
	if t != nil {
		switch t := t.(type) {
		case float64:
			res = time.Unix(int64(t), 0)
		case int64:
			res = time.Unix(t, 0)
		case string:
			res, err = dateparse.ParseAny(t)
			if err != nil {
				return nil, err
			}
		case time.Time:
			res = t
		default:
			return nil, ryerr.New("Date time format is not supported.")
		}
	}

	return &res, nil
}

func ValidateStruct(param interface{}) error {
	validate := validator.New()
	validate.SetTagName("binding")
	if err := validate.Struct(param); err != nil {
		return err
	}
	return nil
}

// IsZero check if any type is its zero value
func IsZero(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Optional return default value if any type is its zero value
func Optional(x interface{}, def interface{}) interface{} {
	if IsZero(x) {
		return def
	}

	return x
}

// Get body params from request (gin context) & allow multiple body reading
func GetBodyParams(c *gin.Context) (interface{}, error) {
	if c.Request.Method == http.MethodGet || c.Request.Body == http.NoBody {
		return nil, nil
	}

	objectBodyParams := make(map[string]interface{})
	b := binding.Default(c.Request.Method, c.ContentType())
	var i interface{} = b
	var err error
	bBody, ok := i.(binding.BindingBody)
	if ok && bBody.Name() == "json" {
		var bodyParams interface{}                     // because body params might not be an object
		err = c.ShouldBindBodyWith(&bodyParams, bBody) // application/json
		if err == nil {
			return bodyParams, nil
		}
	} else if b == binding.FormMultipart {
		err = c.Request.ParseMultipartForm(0) // multipart/form-data
		assignMultipartForm(objectBodyParams, c.Request.MultipartForm)
	} else if b == binding.Form {
		err = c.Request.ParseForm() // application/x-www-form-urlencoded
		for key, value := range c.Request.PostForm {
			objectBodyParams[key] = value
		}
	} else {
		err = c.ShouldBind(&objectBodyParams)
	}

	return objectBodyParams, err
}

func assignMultipartForm(bodyParams map[string]interface{}, multipartForm *multipart.Form) {
	for key, value := range multipartForm.Value {
		bodyParams[key] = value
	}
	for key, files := range multipartForm.File {
		var fileDetails []map[string]interface{}
		for _, file := range files {
			fileDetail := make(map[string]interface{})
			fileDetail["filename"] = file.Filename
			fileDetail["header"] = file.Header
			fileDetail["filesize"] = file.Size
			fileDetails = append(fileDetails, fileDetail)
		}
		bodyParams[key] = fileDetails
	}
}

func GetType(variable interface{}) string {
	return reflect.TypeOf(variable).String()
}
