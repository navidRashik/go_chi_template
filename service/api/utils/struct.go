package utils

import (
	"example_project/internal/config"
	"example_project/internal/database"
	"example_project/internal/leveledlog"
	"sync"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator"
)

type Application struct {
	Config       config.Config
	Db           *database.DB
	Logger       *leveledlog.LogStruct
	Validator    *validator.Validate
	Wg           *sync.WaitGroup
	TokenManager *jwtauth.JWTAuth
}
