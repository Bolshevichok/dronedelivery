package bootstrap

import (
	"fmt"
	"log"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGstorage {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	storage, err := pgstorage.NewPGStorge(connString)
	if err != nil {
		log.Panic(fmt.Sprintf("ошибка инициализации БД, %v", err))
		panic(err)
	}
	return storage
}
