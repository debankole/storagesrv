package output

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type FileWriter struct {
	file       *os.File
	writerChan chan string
	wg         sync.WaitGroup
}

// NewFileWriter creates a new FileWriter, which handles concurrent writes safely.
func NewFileWriter() (*FileWriter, error) {
	file, err := os.Create("output.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}

	fw := &FileWriter{
		file:       file,
		writerChan: make(chan string, 100), // Buffered channel to handle multiple writes
	}

	fw.wg.Add(1)
	go fw.startWriter()

	return fw, nil
}

func (fw *FileWriter) startWriter() {
	defer fw.wg.Done()

	writer := bufio.NewWriter(fw.file)

	for v := range fw.writerChan {
		log.Printf("Writing to file: %s\n", v)
		_, err := writer.WriteString(v)
		if err != nil {
			fmt.Printf("failed to write to file: %v\n", err)
		}
	}

	writer.Flush()
}

// Write safely handles concurrent writes to the file using the channel.
func (fw *FileWriter) Write(data string) {
	fw.writerChan <- data
}

// Close safely closes the file, ensuring all writes are completed.
func (fw *FileWriter) Close() error {
	log.Print("Closing file writer...")
	close(fw.writerChan)
	fw.wg.Wait()
	return fw.file.Close()
}
