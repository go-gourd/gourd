package sessions

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
)

type Store interface {
	sessions.Store
}

func NewFileStore(path string, keyPairs ...[]byte) Store {
	return &fileStore{
		gsessions.NewFilesystemStore(path, keyPairs...),
	}
}

type fileStore struct {
	*gsessions.FilesystemStore
}

func (c *fileStore) Options(options sessions.Options) {
	c.FilesystemStore.Options = options.ToGorillaOptions()
}
