package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/blacktop/go-macho"
	"github.com/blacktop/go-macho/types"
)

var dylibCmdSize = binary.Size(types.DylibCmd{})

func injectLC(fsPath, lcName, tmpdir string) error {
	fat, err := macho.OpenFat(fsPath)
	if err == nil {
		defer fat.Close() // in case of returning early

		var slices []string
		for _, arch := range fat.Arches {
			if err = addDylibCommand(arch.File, lcName); err != nil {
				return err
			}

			tmp, err := os.CreateTemp(tmpdir, "macho_"+arch.File.CPU.String())
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			defer os.Remove(tmp.Name())

			if err = arch.File.Save(tmp.Name()); err != nil {
				return fmt.Errorf("failed to save temp file: %w", err)
			}

			if err = tmp.Close(); err != nil {
				return fmt.Errorf("failed to close temp file: %w", err)
			}

			slices = append(slices, tmp.Name())
		}
		fat.Close()

		// uses os.Create internally, the file will be truncated, everything is fine
		ff, err := macho.CreateFat(fsPath, slices...)
		if err != nil {
			return fmt.Errorf("failed to create fat file: %w", err)
		}
		return ff.Close()
	} else if errors.Is(err, macho.ErrNotFat) {
		m, err := macho.Open(fsPath)
		if err != nil {
			return fmt.Errorf("failed to open MachO file: %w", err)
		}
		defer m.Close()

		if err = addDylibCommand(m, lcName); err != nil {
			return err
		}

		// uses WriteFile internally, it also truncates
		if err = m.Save(fsPath); err != nil {
			return fmt.Errorf("failed to save patched MachO file: %w", err)
		}
		return nil
	}
	return err
}

func addDylibCommand(m *macho.File, name string) error {
	for i := len(m.Loads) - 1; i >= 0; i-- {
		lc := m.Loads[i]
		cmd := lc.Command()
		if cmd != types.LC_LOAD_WEAK_DYLIB && cmd != types.LC_LOAD_DYLIB {
			continue
		}
		if strings.HasPrefix(lc.String(), name) {
			return fmt.Errorf("load command '%s' already exists (already patched)", name)
		}
	}

	var vers types.Version
	vers.Set("0.0.0")

	m.AddLoad(&macho.Dylib{
		DylibCmd: types.DylibCmd{
			LoadCmd:        types.LC_LOAD_WEAK_DYLIB,
			Len:            pointerAlign(uint32(dylibCmdSize + len(name) + 1)),
			NameOffset:     0x18,
			Timestamp:      2, // TODO: I've only seen this value be 2
			CurrentVersion: vers,
			CompatVersion:  vers,
		},
		Name: name,
	})
	return nil
}

func pointerAlign(sz uint32) uint32 {
	if (sz % 8) != 0 {
		sz += 8 - (sz % 8)
	}
	return sz
}
