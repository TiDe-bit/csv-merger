package main

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

const outFile = "output.csv"
const sourceDir = "source"
const targetDir = "target"

func main() {

	out, e := os.Create(outFile)
	if e != nil {
		log.Panic("\nUnable to create new file: ", e)
	}
	defer out.Close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(log.WithError(err))
	}

	srcDir, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Fatalln(log.WithError(err))
	}

	// cd source
	err = os.Chdir(sourceDir)
	if err != nil {
		log.Fatalln(log.WithError(err))
	}

	sourceRecords := make([][]string, 0, 0)

	for _, srcFile := range srcDir {
		records, err := readData(srcFile.Name())
		if err != nil {
			log.Fatalln(log.WithError(err))
		}
		sourceRecords = append(sourceRecords, records...)
	}

	// cd ..
	log.Info("WD\t", wd)

	err = os.Chdir(wd)
	if err != nil {
		log.Fatalln(log.WithError(err))
	}
	currWD, _ := os.Getwd()
	if wd != currWD {
		log.Fatal(currWD, "not correct WD")
	}
	log.Info("current WD:\t", currWD)

	//---------------------------------------

	log.Info("about to enter target dir")

	trgtDir, err := os.ReadDir(targetDir)
	if err != nil {
		log.Fatalln(log.WithError(err))
	}

	err = os.Chdir(targetDir)
	if err != nil {
		log.Fatalln(log.WithError(err))
	}

	recordsTarget, err := readData(trgtDir[0].Name())
	if err != nil {
		log.Fatalln(log.WithError(err))
	}

	writer := csv.NewWriter(out)
	defer writer.Flush()

	if runtime.GOOS == "windows" {
		log.Info(runtime.GOOS)
	}
	writer.UseCRLF = true

	for _, targetRow := range recordsTarget {
		newRow := targetRow
		for entryIndex, sourceRow := range sourceRecords {
			if targetRow[2] == sourceRow[0] {
				newRow = append(newRow, sourceRow[2], sourceRow[9], sourceRow[10])

				// Delete found row
				if len(sourceRecords) <= 2 {
					break
				}
				sourceRecords = append(sourceRecords[:entryIndex], sourceRecords[entryIndex+1])
				break
			}
		}
		err = writer.Write(newRow)
		if err != nil {
			log.Fatalln(log.WithError(err))
		}
	}
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	log.Info(fileName, "\tREADING ALL\t")

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
