package main

import (
	"avito-tech/config"
	"avito-tech/internal/platform/db"
	"avito-tech/internal/server"
	"fmt"
	"log"
)


func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Ошибка считывания конфига (main.go)")
	}
	
	dsn := cfg.GetDSN()

	database, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА ПОДКЛЮЧЕНИЯ К БД: %v", err) 
	}

	defer database.Close()

	s := server.NewServer(cfg, database)
	if err := s.Run(); err != nil{
		log.Fatalln("Just FF")
	}
}

