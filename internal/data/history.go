package data

import (
	"fmt"
	"time"
)

type History struct {
	AllData []HistoryData
}

type HistoryData struct {
	Version   int
	File      string
	Timestemp string
	Note      string
}

func NewHistory() History {

	return History{
		AllData: []HistoryData{},
	}
}

func (h *History) GetBeforeLast() *HistoryData {
	if len(h.AllData) == 0 {
		return nil
	}

	if len(h.AllData) < 2 {
		return nil
	}

	return &h.AllData[len(h.AllData)-2]
}

func (h *History) GetLast() *HistoryData {
	if len(h.AllData) == 0 {
		return nil
	}

	return &h.AllData[len(h.AllData)-1]
}

func (h *History) Add(note string) {

	version := 1
	if len(h.AllData) != 0 {
		lastData := h.AllData[len(h.AllData)-1]
		version = lastData.Version + 1
	}

	data := HistoryData{
		Version:   version,
		File:      fmt.Sprintf("%05d.yaml", version),
		Timestemp: time.Now().Format(time.RFC3339),
		Note:      note,
	}

	h.AllData = append(h.AllData, data)
}

func (h *History) RemoveLast() {
	if len(h.AllData) == 0 {
		return
	}

	h.AllData = h.AllData[:len(h.AllData)-1]
}

func (h *History) List() {

	fmt.Println("----------------------------------------------------------------------------")
	fmt.Printf("%-8s %-16s %-32s %-3s \n", "Version", "File", "Timestemp", "Note")
	fmt.Println("----------------------------------------------------------------------------")
	for i := len(h.AllData) - 1; i >= 0; i-- {
		data := h.AllData[i]
		fmt.Printf("%-8d %-16s %-32s %-3s \n", data.Version, data.File, data.Timestemp, data.Note)
	}
}
