package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/parts/data"
	"github.com/jasontconnell/parts/process"
)

func main() {
	start := time.Now()
	size := flag.Int64("s", 10_000_000, "size of the partitions")
	file := flag.String("f", "", "the file")
	mode := flag.String("m", "partition", "mode of operation (join partitions or partition a file)")
	flag.Parse()

	m := data.Partition
	if *mode == "join" {
		m = data.Join
	}

	if file == nil || *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir, wderr := os.Getwd()
	if wderr != nil {
		log.Fatalf("getting current directory %s", dir)
	}

	var err error
	if m == data.Partition {
		err = process.Partition(dir, *file, *size)
	} else {
		err = process.Join(dir, *file)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished", time.Since(start))
}
