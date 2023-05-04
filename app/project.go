package app

import (
	"crudly/model"
	"crudly/util"
	"crudly/util/result"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type projectCreator interface {
	CreateProject(id model.ProjectId, authInfo model.ProjectAuthInfo) error
}

type projectAuthInfoFetcher interface {
	FetchProjectAuthInfo(id model.ProjectId) result.R[model.ProjectAuthInfo]
}

type projectManager struct {
	projectCreator         projectCreator
	projectAuthInfoFetcher projectAuthInfoFetcher
}

func NewProjectManager(projectCreator projectCreator, projectAuthInfoFetcher projectAuthInfoFetcher) projectManager {
	return projectManager{
		projectCreator,
		projectAuthInfoFetcher,
	}
}

func (p projectManager) GetProjectAuthInfo(id model.ProjectId) result.R[model.ProjectAuthInfo] {
	return p.projectAuthInfoFetcher.FetchProjectAuthInfo(id)
}

func (p projectManager) CreateProject() result.R[model.CreateProjectResponse] {
	id := model.ProjectId(uuid.New())

	key := generateKey()
	keySalt := generateKeySalt()
	keySaltedHash := util.StringHash(string(key) + keySalt)

	err := p.projectCreator.CreateProject(id, model.ProjectAuthInfo{
		Salt:       keySalt,
		SaltedHash: keySaltedHash,
	})

	if err != nil {
		return result.Err[model.CreateProjectResponse](fmt.Errorf("error creating project: %w", err))
	}

	return result.Ok(model.CreateProjectResponse{
		Key: key,
		Id:  id,
	})
}

const PROJECT_KEY_SIZE uint = 80

func generateKey() model.ProjectKey {
	rand.Seed(time.Now().UnixNano())

	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyString := ""

	for i := 0; i < int(PROJECT_KEY_SIZE); i++ {
		randomIndex := rand.Intn(len(alphabet))
		keyString += string(alphabet[randomIndex])
	}

	return model.ProjectKey(keyString)
}

const PROJECT_KEY_SALT_SIZE uint = 10

func generateKeySalt() string {
	rand.Seed(time.Now().UnixNano())

	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := ""

	for i := 0; i < int(PROJECT_KEY_SALT_SIZE); i++ {
		randomIndex := rand.Intn(len(alphabet))
		salt += string(alphabet[randomIndex])
	}

	return salt
}
