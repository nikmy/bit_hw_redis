package main

import (
	"encoding/json"
	"io"
	"os"
	"runtime"

	"go.uber.org/zap"
)

type Entry struct {
	Caption    string `json:"caption"`
	PhotoID    string `json:"photo_id"`
	BusinessID string `json:"business_id"`
	Label      string `json:"label"`
}

func main() {
	log, _ := zap.NewDevelopment()

	file, err := os.Open("yelp_photos.json")
	if err != nil {
		log.Sugar().Fatal(err)
	}
	defer file.Close()

	data := make([]Entry, 0, 200_000)
	for d := json.NewDecoder(file); ; {
		var next Entry
		err = d.Decode(&next)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Sugar().Fatal(err)
		}
		data = append(data, next)
	}
	log.Sugar().Info("Successfully read test data")

	testSingle(log, data)
	runtime.GC()
	testCluster(log, data)
}
