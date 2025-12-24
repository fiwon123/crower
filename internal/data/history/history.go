package history

import (
	"fmt"
	"time"
)

type Data struct {
	AllData []Content
}

type Content struct {
	Version   int
	File      string
	Timestemp string
	Note      string
}

func New() Data {

	return Data{
		AllData: []Content{},
	}
}

func (h *Data) GetBeforeLast() *Content {
	if len(h.AllData) == 0 {
		return nil
	}

	if len(h.AllData) < 2 {
		return nil
	}

	return &h.AllData[len(h.AllData)-2]
}

func (h *Data) GetLast() *Content {
	if len(h.AllData) == 0 {
		return nil
	}

	return &h.AllData[len(h.AllData)-1]
}

func (h *Data) Add(note string) {

	version := 1
	if len(h.AllData) != 0 {
		lastData := h.AllData[len(h.AllData)-1]
		version = lastData.Version + 1
	}

	data := Content{
		Version:   version,
		File:      fmt.Sprintf("%05d.yaml", version),
		Timestemp: time.Now().Format(time.RFC3339),
		Note:      note,
	}

	h.AllData = append(h.AllData, data)
}

func (h *Data) RemoveLast() {
	if len(h.AllData) == 0 {
		return
	}

	h.AllData = h.AllData[:len(h.AllData)-1]
}

func (h *Data) List() {

	fmt.Println("----------------------------------------------------------------------------")
	fmt.Printf("%-8s %-16s %-32s %-3s \n", "Version", "File", "Timestemp", "Note")
	fmt.Println("----------------------------------------------------------------------------")
	for i := len(h.AllData) - 1; i >= 0; i-- {
		data := h.AllData[i]
		fmt.Printf("%-8d %-16s %-32s %-3s \n", data.Version, data.File, data.Timestemp, data.Note)
	}
}
