package history

import (
	"fmt"
	"time"

	"github.com/fiwon123/crower/internal/data/operation"
)

type Data struct {
	AllData []Content
}

type Content struct {
	Version     int
	File        string
	Timestemp   string
	Operation   operation.State
	CommandName string
	Note        string
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

func (h *Data) GetLastOperation(op operation.State) *Content {
	lenAllData := len(h.AllData)

	if lenAllData == 0 {
		return nil
	}

	currentIndex := lenAllData - 1
	for currentIndex >= 0 {
		if h.AllData[currentIndex].Operation == op {
			break

		}

		currentIndex -= 1
	}

	if currentIndex < 0 {
		return nil
	}

	return &h.AllData[currentIndex]
}

func (h *Data) Add(op operation.State, commandName string, note string) {

	version := 1
	if len(h.AllData) != 0 {
		lastData := h.AllData[len(h.AllData)-1]
		version = lastData.Version + 1
	}

	data := Content{
		Version:     version,
		File:        fmt.Sprintf("%05d.yaml", version),
		Timestemp:   time.Now().Format(time.RFC3339),
		Operation:   op,
		CommandName: commandName,
		Note:        note,
	}

	h.AllData = append(h.AllData, data)
}

func (h *Data) RemoveLast() {
	if len(h.AllData) == 0 {
		return
	}

	h.AllData = h.AllData[:len(h.AllData)-1]
}

func (h *Data) ListLast(steps int) {

	stop := len(h.AllData) - steps
	printHeader()
	for i := len(h.AllData) - 1; i >= stop; i-- {
		data := h.AllData[i]
		printContent(data.Version, data.File, data.Timestemp, data.Note)
	}
}

func (h *Data) ListGoBack(steps int) {
	start := len(h.AllData) - 1 - steps

	printHeader()
	for i := start; i >= 0; i-- {
		data := h.AllData[i]
		printContent(data.Version, data.File, data.Timestemp, data.Note)
	}
}

func (h *Data) List() {
	printHeader()
	for i := len(h.AllData) - 1; i >= 0; i-- {
		data := h.AllData[i]
		printContent(data.Version, data.File, data.Timestemp, data.Note)
	}
}

func printContent(version int, file string, timestemp string, note string) {
	fmt.Printf("%-8d %-16s %-32s %-3s \n", version, file, timestemp, note)
}

func printHeader() {
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Printf("%-8s %-16s %-32s %-3s \n", "Version", "File", "Timestemp", "Note")
	fmt.Println("----------------------------------------------------------------------------")
}
