package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// jsonデコード用の構造体
type Person struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
}

func main() {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("sample.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	// できるだけPerson型にマッピングしようとする
	var persons []Person
	if err := json.Unmarshal(bytes, &persons); err != nil {
		log.Fatal(err)
	}

	// デコードしたデータを表示
	// 中身の要素をそれぞれ確保する
	for _, p := range persons {
		fmt.Printf("%d : %s %s\n", p.Id, p.Name, p.Birthday)
	}
}
