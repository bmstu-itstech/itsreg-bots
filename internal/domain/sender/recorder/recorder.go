package recorder

import (
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"sync"
)

type Record struct {
	Receiver entity.ParticipantId
	Text     string
}

type Recorder struct {
	records []Record
	sync.Mutex
}

func New() *Recorder {
	return &Recorder{
		records: make([]Record, 0),
	}
}

func (r *Recorder) SendMessage(receiver entity.ParticipantId, msg string, _ []entity.Button) error {
	record := Record{
		Receiver: receiver,
		Text:     msg,
	}

	r.Lock()
	r.records = append(r.records, record)
	r.Unlock()

	return nil
}

func (r *Recorder) GetLastRecords() []Record {
	records := make([]Record, len(r.records))
	for i, record := range r.records {
		records[i] = record
	}

	r.records = make([]Record, 0)

	return records
}
