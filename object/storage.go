package object

type Storage interface {
	GetUploadParams(key string) (string, string, map[string]string, error)
	GetDownloadUrl(key string) (*string, error)
}
