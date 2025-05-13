package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tt := time.Unix(1672531200, 0) // 2023-01-01 00:00:00 UTC
	jsonData, _ := json.Marshal(tt)
	fmt.Println(string(jsonData))

	tt = time.Unix(1672531200, 0).UTC() // 2023-01-01 00:00:00 UTC
	jsonData, _ = json.Marshal(tt)
	fmt.Println(string(jsonData))
}
