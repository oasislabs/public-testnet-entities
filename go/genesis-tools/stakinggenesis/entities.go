package stakinggenesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/oasislabs/oasis-core/go/common/entity"
	"github.com/oasislabs/oasis-core/go/common/logging"
	registry "github.com/oasislabs/oasis-core/go/registry/api"
)

var (
	logger = logging.GetLogger("stakinggenesis")
)

type Entities interface {
	All() map[string]*entity.Entity
	ResolveEntity(name string) (*entity.Entity, error)
}

// EntitiesDirectory is a directory of unpacked entities packages.
type EntitiesDirectory struct {
	path string

	// A map of Entity Names to the Entity object
	entities map[string]*entity.Entity
}

// LoadEntitiesDirectory loads a directory of unpacked entity packages.
func LoadEntitiesDirectory(dirPath string) (*EntitiesDirectory, error) {
	dir := &EntitiesDirectory{path: dirPath}

	dir.Load()

	return dir, nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (e *EntitiesDirectory) All() map[string]*entity.Entity {
	return e.entities
}

// Load loads a directory of entities. This should a directory of unpacked
// entity packages.
func (e *EntitiesDirectory) Load() error {
	files, err := ioutil.ReadDir(e.path)
	if err != nil {
		logger.Error("failed to load the entities directory",
			"err", err,
		)
	}
	entities := make(map[string]*entity.Entity)
	for _, fileInfo := range files {
		// Only process directories.
		if !fileInfo.IsDir() {
			continue
		}
		entityName := fileInfo.Name()
		ent, err := e.loadEntityDir(entityName)
		if err != nil {
			return err
		}
		entities[entityName] = ent
	}
	e.entities = entities
	return nil
}

// ResolveEntity resolves an entity name to an Entity.
func (e *EntitiesDirectory) ResolveEntity(name string) (*entity.Entity, error) {
	ent, ok := e.entities[name]
	if !ok {
		return nil, fmt.Errorf("Entity %s does not exist", name)
	}
	return ent, nil
}

func (e *EntitiesDirectory) loadEntityDir(entityName string) (*entity.Entity, error) {
	entityGenesisPath := path.Join(e.path, entityName, "entity/entity_genesis.json")
	logger.Debug("loading entity directory", "dir", entityGenesisPath)
	if !isFile(entityGenesisPath) {
		return nil, fmt.Errorf("Entity for \"%s\" does not exist", entityName)
	}

	b, err := ioutil.ReadFile(entityGenesisPath)
	if err != nil {
		return nil, err
	}

	var signedEntity entity.SignedEntity
	if err = json.Unmarshal(b, &signedEntity); err != nil {
		return nil, err
	}

	var ent entity.Entity
	if err := signedEntity.Open(registry.RegisterGenesisEntitySignatureContext, &ent); err != nil {
		return nil, err
	}

	return &ent, nil
}
