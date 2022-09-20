package helpers

import (
	"fmt"
	"os"
	"time"
)

func Fatal(str string, err error) {
	if err != nil {
		fmt.Println(str)
		fmt.Println(err.Error())
		fmt.Println("** Programm quits in 10 seconds!")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
}
