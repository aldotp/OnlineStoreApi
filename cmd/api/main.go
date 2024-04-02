package main

import (
	"github.com/aldotp/OnlineStore/internal/config"
	"github.com/aldotp/OnlineStore/internal/route"
)

func main() {
	viper := config.NewViper()
	db, err := config.NewDB(viper)
	if err != nil {
		panic(err)
	}

	config := config.NewBoostrapConfig(db, viper)
	r := route.NewRouter(config)
	r.Run()
}
