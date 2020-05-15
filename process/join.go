package process

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jasontconnell/partfile/data"
	"github.com/pkg/errors"
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
				return errors.Wrapf(err, "couldn't parse partition number from %s", name)
			}

			partition := data.Part{Filename: name, Index: pnum}
			partitions = append(partitions, partition)
		}

		return nil
	})

	if len(partitions) == 0 {
		return errors.New("no parition files found")
	}

	if err != nil {
		return errors.Wrapf(err, "walking directory %s", dir)
	}

	sort.Slice(partitions, func(i, j int) bool {
		return partitions[i].Index < partitions[j].Index
	})

	err = doJoin(dir, filename, partitions)
	if err != nil {
		return errors.Wrap(err, "joining files")
	}

	// maybe delete files after

	return nil
}

func doJoin(dir, filename string, parts []data.Part) error {
	fullPath := filepath.Join(dir, filename)
	outfile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "opening file to join to - %s", fullPath)
	}
	defer outfile.Close()

	var pos int64
	for _, p := range parts {
		partfile := filepath.Join(dir, p.Filename)
		b, err := ioutil.ReadFile(partfile)

		if err != nil {
			return errors.Wrapf(err, "opening part file %s %d", partfile, p.Index)
		}

		_, err = outfile.WriteAt(b, pos)
		pos += int64(len(b))

		if err != nil {
			return errors.Wrapf(err, "writing file %s %d", fullPath, p.Index)
		}
	}
	return nil
}
