# go-aht20

https://github.com/tinygo-org/drivers/tree/release/aht20 的 https://github.com/d2r2/go-i2c 移植

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
	aht20 := aht20.AHT20New(bus)
	err = aht20.ReadWithRetry(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("温度:", aht20.Celsius(), "摄氏度")
	fmt.Println("相对湿度:", aht20.RelHumidity(), "%")
}
```
