package process

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jasontconnell/parts/data"
)

func Join(dir, filename string) error {

	partitions := []data.Part{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != dir {
			return filepath.SkipDir
		}

		name := info.Name()

		if strings.HasPrefix(name, filename) && name != filename {
			ext := strings.TrimPrefix(filepath.Ext(name), ".")
			pnum, err := strconv.Atoi(ext)
			if err != nil {
				return fmt.Errorf("couldn't parse partition number from %s. %w", name, err)
			}

			partition := data.Part{Filename: name, Index: pnum}
			partitions = append(partitions, partition)
		}

		return nil
	})

	if len(partitions) == 0 {
		return fmt.Errorf("no parition files found")
	}

	if err != nil {
		return fmt.Errorf("walking directory %s. %w", dir, err)
	}

	sort.Slice(partitions, func(i, j int) bool {
		return partitions[i].Index < partitions[j].Index
	})

	err = doJoin(dir, filename, partitions)
	if err != nil {
		return fmt.Errorf("joining files. %w", err)
	}

	// maybe delete files after

	return nil
}

func doJoin(dir, filename string, parts []data.Part) error {
	fullPath := filepath.Join(dir, filename)
	outfile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("opening file to join to - %s. %w", fullPath, err)
	}
	defer outfile.Close()

	var pos int64
	for i, p := range parts {
		if i != p.Index {
			return fmt.Errorf("missing part %d - %s", i, fullPath)
		}
		partfile := filepath.Join(dir, p.Filename)
		b, err := os.ReadFile(partfile)

		if err != nil {
			return fmt.Errorf("opening part file %s %d. %w", partfile, p.Index, err)
		}

		_, err = outfile.WriteAt(b, pos)
		pos += int64(len(b))

		if err != nil {
			return fmt.Errorf("writing file %s %d. %w", fullPath, p.Index, err)
		}
	}
	return nil
}
