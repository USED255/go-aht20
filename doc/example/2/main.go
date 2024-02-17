package main

import (
	"fmt"
	"log"

	"github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"

	"github.com/used255/go-aht20"
)

func main() {
	//logger.ChangePackageLogLevel("i2c", logger.DebugLevel)
	logger.ChangePackageLogLevel("aht20", logger.DebugLevel)

	i2c, err := i2c.NewI2C(0x38, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	sensor := aht20.NewAHT20(i2c)

	celsius, err := sensor.ReadCelsius()
	if err != nil {
		log.Fatal(err)
	}

	humidity, err := sensor.ReadRelHumidity()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("温度:", celsius, "摄氏度")
	fmt.Println("相对湿度:", humidity, "%")
}
