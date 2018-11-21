package storage

import "github.com/diraven/secret-server-task/server/models"

type putter interface {
	Put(
		secretData string, // data to be stored
		expireAfterViews int32, // times
		expireAfter int32, // minutes
	) (secret *models.Secret, err error)
}

type getter interface {
	Get(hash string) (secret *models.Secret, err error)
}

// SecretStorage allows to store and retrieve secrets.
type Secret interface {
	putter
	getter
}
