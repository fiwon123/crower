package history

import (
	"fmt"
	"time"

	"github.com/fiwon123/crower/internal/data/state"
)

type Data struct {
	AllData []Content
}

type Content struct {
	Version     int
	File        string
	Timestemp   string
	Operation   state.OperationEnum
	CommandName string
	Note        string
}

// Create new History
func New() Data {

	return Data{
		AllData: []Content{},
	}
}

// Get last content from lastIndex - steps if possible
func (h *Data) GetBeforeLast(steps int) (*Content, error) {
	if len(h.AllData) == 0 {
		return nil, fmt.Errorf("history is empty")
	}

	if len(h.AllData)-steps < 0 {
		return nil, fmt.Errorf("step value is greater than quantity of history registry")
	}

	return &h.AllData[len(h.AllData)-steps], nil
}

// Get last content from lastIndex
func (h *Data) GetLast() *Content {
	if len(h.AllData) == 0 {
		return nil
	}

	return &h.AllData[len(h.AllData)-1]
}

// Get last content based on the last operation using crower
func (h *Data) GetLastOperation(op state.OperationEnum) *Content {
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

// Add new content history registering current operation, command name and a note
func (h *Data) Add(op state.OperationEnum, commandName string, note string) {

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

// Remove last content
func (h *Data) RemoveLast() {
	if len(h.AllData) == 0 {
		return
	}

	h.AllData = h.AllData[:len(h.AllData)-1]
}

// Get Index from last index - steps
func (h *Data) GetIndexFromLastTo(steps int) int {
	return len(h.AllData) - steps
}

// List from first to steps
func (h *Data) ListFirstHistory(steps int) {

	start := steps - 1
	printHeader()
	for i := start; i >= 0; i-- {
		data := h.AllData[i]
		printContent(data.Version, data.File, data.Timestemp, data.Note)
	}
}

// List from last to steps
func (h *Data) ListLastHistory(steps int) {
	start := len(h.AllData) - 1
	stop := len(h.AllData) - steps
	printHeader()
	for i := start; i >= stop; i-- {
		data := h.AllData[i]
		printContent(data.Version, data.File, data.Timestemp, data.Note)
	}
}

// List all history
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
