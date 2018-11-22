package storage

import "github.com/diraven/secret-server-task/server/models"

// Secret allows to Put secrets to storage and Get them out of there.
type Secret interface {
	Put(
		secretData string, // data to be stored
		expireAfterViews int32, // times
		expireAfter int32, // minutes
	) (secret *models.Secret, err error)
	Get(hash string) (secret *models.Secret, err error)
}
