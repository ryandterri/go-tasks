package lib

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

func loadFile(filepath string, flag int, mode os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(filepath, flag, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func closeFile(f *os.File) error {
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if err != nil {
		panic(err)
	}
	return f.Close()
}

func ReadList() List {
	file, err := loadFile("list.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := closeFile(file); err != nil {
			panic(err)
		}
	}()
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var list List
	for _, record := range records {
		var item ListItem
		item.Id, _ = strconv.Atoi(record[0])
		item.Description = record[1]
		item.CreatedAt, _ = time.Parse(time.RFC3339, record[2])
		item.CompletedAt, _ = time.Parse(time.RFC3339, record[3])
		list.Items = append(list.Items, item)
	}
	return list
}

func WriteList(list *List) {
	file, err := loadFile("list.csv", os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := closeFile(file); err != nil {
			panic(err)
		}
	}()

	writer := csv.NewWriter(file)
	// header := []string{"Id", "Description", "CreatedAt", "IsComplete"}
	// _ = writer.Write(header)
	for _, item := range list.Items {
		completedAt := ""
		if !item.CompletedAt.IsZero() {
			completedAt = item.CompletedAt.Format(time.RFC3339)
		}
		row := []string{strconv.Itoa(item.Id), item.Description, item.CreatedAt.Format(time.RFC3339), completedAt}
		_ = writer.Write(row)
	}
	writer.Flush()
}
