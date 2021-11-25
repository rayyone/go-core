package contenttype

type ContentType string

const (
	ApplicationJSON ContentType = "application/json"
	FormUrlencoded  ContentType = "application/x-www-form-urlencoded"
	FormData        ContentType = "multipart/form-data"
)
