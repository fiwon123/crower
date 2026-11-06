package main

import (
	"debug/pe"
	"fmt"
	"os"
	"sort"
)

// Windows subsystem values from the PE optional header. A console program gets
// its own window from the loader; a GUI program does not, which is what makes
// double-clicking crower-gui.exe open the window and nothing else.
const (
	subsystemGUI     = 2
	subsystemConsole = 3
)

// verifyGUI checks a built executable is a GUI program carrying the icon and
// manifest this tool generates. It is the gate that keeps a plain `go build`
// from shipping crower-gui.exe as a console binary again: the check inspects
// the built file, not the command that produced it.
func verifyGUI(path string) error {
	file, err := pe.Open(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	defer file.Close()

	if subsystem := subsystemOf(file); subsystem != subsystemGUI {
		return fmt.Errorf("%s: subsystem %d, want %d (a console binary opens a console window on double-click)", path, subsystem, subsystemGUI)
	}

	types, err := resourceTypes(file)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	for _, want := range []uint32{rtIcon, rtGroupIcon, rtManifest, rtVersion} {
		if !types[want] {
			return fmt.Errorf("%s: no resource of type %d (icon 3, group icon 14, manifest 24, version 16); found %v", path, want, sortedTypes(types))
		}
	}

	return nil
}

// verifyConsole checks a built executable is a console program, the way the
// cr/crower CLI pair must be.
func verifyConsole(path string) error {
	file, err := pe.Open(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	defer file.Close()

	if subsystem := subsystemOf(file); subsystem != subsystemConsole {
		return fmt.Errorf("%s: subsystem %d, want %d (a terminal command must keep its console)", path, subsystem, subsystemConsole)
	}

	return nil
}

func subsystemOf(file *pe.File) uint16 {
	switch header := file.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		return header.Subsystem
	case *pe.OptionalHeader64:
		return header.Subsystem
	}

	return 0
}

// resourceTypes lists the resource types (RT_ICON, RT_GROUP_ICON, ...) present
// in the image. The directory is a tree: the root's entries are the types, each
// pointing at a subdirectory of names.
func resourceTypes(file *pe.File) (map[uint32]bool, error) {
	var (
		rva  uint32
		size uint32
	)

	switch header := file.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		if len(header.DataDirectory) > 2 {
			rva, size = header.DataDirectory[2].VirtualAddress, header.DataDirectory[2].Size
		}
	case *pe.OptionalHeader64:
		if len(header.DataDirectory) > 2 {
			rva, size = header.DataDirectory[2].VirtualAddress, header.DataDirectory[2].Size
		}
	}

	if rva == 0 || size == 0 {
		return nil, fmt.Errorf("no resource directory")
	}

	offset, err := fileOffset(file, rva)
	if err != nil {
		return nil, err
	}

	data, err := sectionData(file, rva)
	if err != nil {
		return nil, err
	}

	if offset+16 > uint32(len(data)) {
		return nil, fmt.Errorf("resource directory is truncated")
	}

	count := int(le16(data[offset+14:]))
	types := make(map[uint32]bool, count)

	for index := 0; index < count; index++ {
		at := offset + 16 + uint32(8*index)
		if at+8 > uint32(len(data)) {
			break
		}

		types[le32(data[at:])] = true
	}

	return types, nil
}

// fileOffset maps an RVA to its offset inside the containing section's data.
func fileOffset(file *pe.File, rva uint32) (uint32, error) {
	for _, section := range file.Sections {
		if rva >= section.VirtualAddress && rva < section.VirtualAddress+section.Size {
			return rva - section.VirtualAddress, nil
		}
	}

	return 0, fmt.Errorf("address %#x is not in any section", rva)
}

// sectionData returns the raw bytes of the section holding rva.
func sectionData(file *pe.File, rva uint32) ([]byte, error) {
	for _, section := range file.Sections {
		if rva >= section.VirtualAddress && rva < section.VirtualAddress+section.Size {
			return section.Data()
		}
	}

	return nil, fmt.Errorf("address %#x is not in any section", rva)
}

func le16(data []byte) uint16 {
	if len(data) < 2 {
		return 0
	}

	return uint16(data[0]) | uint16(data[1])<<8
}

func le32(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}

	return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
}

func sortedTypes(types map[uint32]bool) []uint32 {
	out := make([]uint32, 0, len(types))

	for kind := range types {
		out = append(out, kind)
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })

	return out
}

// verifyPath is the CLI entry for -verify and -verify-console.
func verifyPath(path string, gui bool) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	if gui {
		return verifyGUI(path)
	}

	return verifyConsole(path)
}
