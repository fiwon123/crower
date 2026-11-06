// Command winres writes the Windows resource object a Go build links into a
// binary: the application icon, a DPI-aware manifest and the version block
// Explorer shows in a file's Properties.
//
// Go has no flag for this. The only hook is a .syso file sitting beside the
// package being built, so a Windows GUI build has to generate one first:
//
//	go run ./internal/tools/winres -icon internal/fyneui/logo.png \
//	    -out cmd/crowergui/rsrc_windows_amd64.syso -arch amd64 -version v1.2.3
//
// The object is plain COFF with a .rsrc section; no Windows toolchain is
// involved, so it can be produced on any host.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	icon := flag.String("icon", "", "PNG image to use as the application icon (required)")
	out := flag.String("out", "", "path of the .syso file to write (required)")
	arch := flag.String("arch", "amd64", "target architecture: amd64, arm64 or 386")
	version := flag.String("version", "", "version to record, such as v0.20.15")
	name := flag.String("name", "crower", "product and file description")
	verify := flag.String("verify", "", "check a built GUI executable instead of writing one")
	verifyConsole := flag.String("verify-console", "", "check a built console executable instead of writing one")

	flag.Parse()

	if *verify != "" {
		run(*verify, true)
		return
	}

	if *verifyConsole != "" {
		run(*verifyConsole, false)
		return
	}

	if *icon == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "winres: -icon and -out are required")
		flag.Usage()
		os.Exit(2)
	}

	object, err := build(*icon, *name, *version, *arch)
	if err != nil {
		fmt.Fprintln(os.Stderr, "winres:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*out, object, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "winres:", err)
		os.Exit(1)
	}
}

// run verifies one built executable and exits non-zero when it does not match
// what the name promises.
func run(path string, gui bool) {
	if err := verifyPath(path, gui); err != nil {
		fmt.Fprintln(os.Stderr, "winres:", err)
		os.Exit(1)
	}

	kind := "console"
	if gui {
		kind = "GUI with icon, manifest and version"
	}

	fmt.Printf("ok   %s (%s)\n", path, kind)
}

// build turns the icon and metadata into the bytes of a COFF object.
func build(iconPath, name, version, arch string) ([]byte, error) {
	machine, err := machineFor(arch)
	if err != nil {
		return nil, err
	}

	source, err := os.Open(iconPath)
	if err != nil {
		return nil, fmt.Errorf("open icon: %w", err)
	}
	defer source.Close()

	images, err := iconImages(source)
	if err != nil {
		return nil, err
	}

	section, relocs, err := resources(images, manifest(name, version), versionInfo(name, version))
	if err != nil {
		return nil, err
	}

	return coff(machine, section, relocs), nil
}

func machineFor(arch string) (uint16, error) {
	switch arch {
	case "amd64":
		return machineAMD64, nil
	case "arm64":
		return machineARM64, nil
	case "386":
		return machineI386, nil
	default:
		return 0, fmt.Errorf("unsupported architecture %q", arch)
	}
}
