package json_memory

import (
	"encoding/json"
	"github.com/diraven/secret-server-task/server/models"
	"github.com/diraven/secret-server-task/server/storage"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"time"
)

// NewJSONMemory creates new json/memory based storage to store secrets.
func NewJSONMemory(path string) (storage storage.Secret, err error) {
	// Prepare storage for usage.
	var store = &jsonInMemory{}

	// Load data into the storage.
	if err = store.init(path); err != nil {
		err = errors.Wrap(err, "unable to load the json data from file into storage")
		return
	}

	// Just in case, make sure we can save the data into the storage and return the error early if we are unable to.
	if err = store.sync(); err != nil {
		err = errors.Wrap(err, "unable to save data into the file on disk")
		return
	}

	storage = store
	return
}

// jsonInMemory stores secret on disk in json files while caching all of them in memory.
type jsonInMemory struct {
	Secrets map[string]*models.Secret
	path    string
}

// Put saves secret into the storage.
func (m *jsonInMemory) Put(
	secretData string, // data to be stored
	expireAfterViews int32, // times
	expireAfter int32, // minutes
) (secret *models.Secret, err error) {
	// Make UUID to serve as a hash. It's not cryptographically secure, one might pick some other library/function
	// to generate UUIDs instead if crypto security is necessary.

	// https://github.com/satori/go.uuid claims to conform to the https://tools.ietf.org/html/rfc4122, which
	// states that UUIDs in question are unique for the foreseeable future.
	u2, err := uuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "unable to generate UUID")
		return
	}

	// Get current time.
	var now = time.Now()

	// Generate our secret object.
	secret = &models.Secret{
		CreatedAt: strfmt.DateTime(now),
		// If expireAfter is et to 0, CreatedAt and ExpiresAt will be the same value. This means the secret should
		// never expire.
		ExpiresAt:      strfmt.DateTime(now.Add(time.Minute * time.Duration(expireAfter))),
		Hash:           u2.String(),
		RemainingViews: expireAfterViews,
		SecretText:     secretData,
	}

	// Add our secret to the database.
	m.Secrets[secret.Hash] = secret

	// Save changes.
	if err = m.sync(); err != nil {
		err = errors.Wrap(err, "unable to save changes to disk")
		return
	}

	return
}

func (m *jsonInMemory) Get(hash string) (secret *models.Secret, err error) {
	var syncPending bool

	// Clean up expired secrets.
	// If any secrets were cleaned up:
	if m.clean() > 0 {
		// Make sure to flush changes to the file on disk.
		syncPending = true
	}

	// Look for the secret in our database and return it if found.
	var ok bool
	secret, ok = m.Secrets[hash]

	// If we found the secret - validate and return it.
	if ok {
		// Handle views count.
		secret.RemainingViews--

		// If no views left:
		if secret.RemainingViews == 0 {
			// That was the last time we served this secret.
			delete(m.Secrets, secret.Hash)
		}

		// Make sure to save changes as we need to track the remaining views.
		syncPending = true
	}

	// Save changes if any.
	if syncPending {
		if err = m.sync(); err != nil {
			err = errors.Wrap(err, "unable to save changes to disk")
			return
		}
	}

	// No secrets found.
	return
}

func (m *jsonInMemory) init(path string) (err error) {
	// Save path for later use.
	m.path = path

	// Open JSON file.
	var jsonFile *os.File
	if jsonFile, err = os.OpenFile(m.path, os.O_RDONLY|os.O_CREATE, 0600); err != nil {
		err = errors.Wrap(err, "unable to open json file")
		return
	}
	defer jsonFile.Close()

	// Get JSON data.
	var data []byte
	if data, err = ioutil.ReadAll(jsonFile); err != nil {
		err = errors.Wrap(err, "unable to read json data from file")
		return
	}

	// If data provided:
	if len(data) > 0 {
		// Parse JSON.
		if err = json.Unmarshal(data, &m.Secrets); err != nil {
			err = errors.Wrap(err, "unable to parse json data")
			return
		}
	} else {
		// Initialize empty secrets map otherwise.
		m.Secrets = map[string]*models.Secret{}
	}

	return
}

func (m *jsonInMemory) sync() (err error) {
	// Marshal the secrets data.
	var data []byte
	if data, err = json.Marshal(m.Secrets); err != nil {
		err = errors.Wrap(err, "unable to marshal json data")
		return
	}

	// Write data info the file.
	if err = ioutil.WriteFile(m.path, data, 0600); err != nil {
		err = errors.Wrap(err, "unable to write data to file")
		return
	}

	return
}

// clean removes all expired secrets.
func (m *jsonInMemory) clean() (count int) {
	// Storage for hashes staged for deletion.
	var stagedForDeletion []string

	// Find expired hashes.
	for hash, secret := range m.Secrets {
		// Creation date will be the same as expiration date for items that do not expire.
		// If now is later then expiration date and expiration date is different from creation date:
		if time.Now().After(time.Time(secret.ExpiresAt)) &&
			time.Time(secret.ExpiresAt).After(time.Time(secret.CreatedAt)) {
			stagedForDeletion = append(stagedForDeletion, hash)
		}
	}

	// Deleted items count to be returned.
	count = len(stagedForDeletion)

	// Delete expired hashes, if any.
	for _, hash := range stagedForDeletion {
		delete(m.Secrets, hash)
	}

	return
}
