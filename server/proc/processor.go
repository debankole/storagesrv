package proc

import (
	"fmt"
	"log"
	"storage/collections"
	"storage/common"
	"sync"
)

type MessageReader interface {
	ReadMessages() <-chan common.Command
}

type Output interface {
	Write(string)
}

type CommandProcessor struct {
	reader MessageReader
	out    Output
	omap   *collections.OrderedMap[string, any]
	wg     sync.WaitGroup
}

func NewCommandProcessor(r MessageReader, o Output) *CommandProcessor {
	return &CommandProcessor{
		reader: r,
		out:    o,
		omap:   collections.NewOrderedMap[string, any](),
	}
}

// Start starts the command processor with the specified number of workers.
func (s *CommandProcessor) Start(workers int) {
	for i := 0; i < workers; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// Stop stops the command processor and waits for all workers to finish.
func (s *CommandProcessor) Stop() {
	log.Printf("Stopping command processor...")
	s.wg.Wait()
}

func (s *CommandProcessor) worker() {
	defer s.wg.Done()

	for cmd := range s.reader.ReadMessages() {
		s.executeCommand(cmd)
	}
}

func (s *CommandProcessor) executeCommand(cmd common.Command) {
	log.Printf("Processing command: %s", cmd.Type)
	switch cmd.Type {
	case common.AddItem:
		s.omap.AddItem(cmd.Key, cmd.Value)
		log.Printf("add-item: Key: %s, Value: %v", cmd.Key, cmd.Value)
	case common.DeleteItem:
		s.omap.RemoveItem(cmd.Key)
		log.Printf("delete-item: Key: %s", cmd.Key)
	case common.GetItem:
		if v, ok := s.omap.GetItem(cmd.Key); ok {
			s.out.Write(fmt.Sprintf("Key: %s, Value: %s\n", cmd.Key, v))
			log.Printf("get-item: Key: %s, Value: %s", cmd.Key, v)
		} else {
			log.Printf("item not found")
		}
	case common.GetAllItems:
		items := s.omap.GetAllItems()
		for _, item := range items {
			s.out.Write(fmt.Sprintf("Key: %s, Value: %v\n", item.Key, item.Value))
		}

		log.Printf("get-all-items: %d items", len(items))

	default:
		log.Printf("unknown command: %s", cmd.Type)
	}
}
