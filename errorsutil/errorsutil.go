package errorsutil

import "fmt"

func HandleError(err error) {
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}
