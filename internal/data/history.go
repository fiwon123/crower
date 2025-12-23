package data

import (
	"fmt"
	"time"
)

type History struct {
	allData []HistoryData
}

type HistoryData struct {
	Version   int
	File      string
	Timestemp time.Time
	Note      string
}

func (h *History) Add(data HistoryData) {
	h.allData = append(h.allData, data)
}

func (h *History) List() {

	fmt.Println("-------------------------------")
	fmt.Printf("%-8s %-16s %-16s %-3s \n", "Version", "File", "Timestemp", "Note")
	fmt.Println("-------------------------------")
	for _, data := range h.allData {
		fmt.Printf("%-8d %-16s %-16s %-3s \n", data.Version, data.File, data.Timestemp, data.Note)
	}
}
