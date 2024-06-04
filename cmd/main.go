package main

import (
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	cfg := newViperConfig()

	dbCfg := newDatabase(cfg.GetString("DATABASE_URL"))
	var repo repository.RepositoryInterface = repository.NewRepository(repository.Repository{Db: dbCfg})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}

func newDatabase(dbURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	return db
}

// for local development to read the env file
func newViperConfig() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("../../../../params")
	v.AddConfigPath("./params")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	v.ReadInConfig()

	return v
}
