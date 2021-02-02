package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/beevik/guid"
	"github.com/oguzhn/PlayerRankingService/models"
)

func main() {
	var data models.BulkUserDTO
	data.Count = 10000
	data.List = make(models.UserDTOList, data.Count)
	for i := 0; i < data.Count; i++ {
		data.List[i].ID = guid.NewString()
		data.List[i].Score = float32(rand.Int())
		data.List[i].CountryCode = "tr"
		data.List[i].Name = "gjg" + fmt.Sprint(i)
	}
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("data.json", file, 0644)
}
