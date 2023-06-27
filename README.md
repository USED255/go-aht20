# go-aht20

[![Go Report Card](https://goreportcard.com/badge/github.com/used255/go-aht20)](https://goreportcard.com/report/github.com/used255/go-aht20)
[![MIT License](https://img.shields.io/badge/License-MIT-green)](./LICENSE)

https://github.com/tinygo-org/drivers/tree/release/aht20 的 https://github.com/d2r2/go-i2c 移植

[AHT20 产品规格书](http://www.aosong.com/userfiles/files/media/AHT20%E4%BA%A7%E5%93%81%E8%A7%84%E6%A0%BC%E4%B9%A6(%E4%B8%AD%E6%96%87%E7%89%88)%20B1.pdf)
[AHT20 Datasheet](http://www.aosong.com/userfiles/files/media/AHT20%20%E8%8B%B1%E6%96%87%E7%89%88%E8%AF%B4%E6%98%8E%E4%B9%A6%20A0%2020201222.pdf)

```go
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

	bus, err := i2c.NewI2C(0x38, 1)
	if err != nil {
		log.Fatal(err)
	}

	s := aht20.NewAHT20(bus)

	err = s.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("温度:", s.Celsius(), "摄氏度")
	fmt.Println("相对湿度:", s.RelHumidity(), "%")
}
```

## 您也可以看看
- [AHT20 Temperature & Humidity for Python I2C](https://github.com/Chouffy/python_sensor_aht20)
- [I2C-bus interaction of peripheral sensors with Raspberry PI embedded Linux or respective clones](https://github.com/d2r2/go-i2c)
- [TinyGo Drivers](https://github.com/tinygo-org/drivers)
