package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("version1 NewUUID --")
	for i := 0; i < 10; i++ {
		uuidObj, _ := uuid.NewUUID()
		fmt.Println("  ", uuidObj.String())
	}

	fmt.Println("version2 NewDCESecurity --")
	for i := 0; i < 10; i++ {
		uuidObj, _ := uuid.NewUUID()
		domain := uuidObj.Domain()
		id := uuidObj.ID()

		uuidObj2, _ := uuid.NewDCESecurity(domain, id)
		fmt.Println("  ", uuidObj2.String())
	}

	fmt.Println("version3 NewMD5 --")
	for i := 0; i < 10; i++ {
		uuidObj, _ := uuid.NewUUID()
		data := []byte("wnw8olzvmjp0x6j7ur8vafs4jltjabi0")
		uuidObj2 := uuid.NewMD5(uuidObj, data)
		fmt.Println("  ", uuidObj2.String())
	}

	fmt.Println("version4 NewRandom --")
	for i := 0; i < 10; i++ {
		uuidObj, _ := uuid.NewRandom()
		fmt.Println("  ", uuidObj.String())
	}

	fmt.Println("version5 NewSHA1 --")
	for i := 0; i < 10; i++ {
		uuidObj, _ := uuid.NewUUID()
		data := []byte("wnw8olzvmjp0x6j7ur8vafs4jltjabi0")
		uuidObj2 := uuid.NewSHA1(uuidObj, data)
		fmt.Println("  ", uuidObj2.String())
	}

}
