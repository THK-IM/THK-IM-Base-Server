package object

type Storage interface {
	UploadObject(key string, path string) (*string, error)
	GetUploadParams(key string) (string, string, map[string]string, error)
	GetDownloadUrl(key string) (*string, error)
	DeleteObjectsByKeys(keys []string) error
	KeyExists(key string) (bool, error)
}
