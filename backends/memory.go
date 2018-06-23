package backends

import (
	"sync"
)

type Memory struct {
	Saver map[string]Record
	sync.RWMutex
}

func NewMemory() (*Memory, error){
	return &Memory{
		Saver: map[string]Record{},
	}, nil
}

func (m *Memory) getValue() error {
	//value, ok := m.Saver[domain]
	return nil
}

func (m *Memory) DeleteRecord(record Record) error {
	m.Lock()
	defer m.Unlock()
	delete(m.Saver, record.Name)
	return nil
}

func (m *Memory) SaveRecord(record Record) error {
	m.Lock()
	defer m.Unlock()
	m.Saver[record.Name] = record
	return nil
}

func (m *Memory) GetRecord(name string) (Record, bool) {
	m.Lock()
	defer m.Unlock()
	rec, ok := m.Saver[name]
	if ok && rec.Valid() {
		return rec, true
	}
	return Record{}, false
}

func (m *Memory) UpdateRecord(record Record) error {
	m.Lock()
	defer m.Unlock()
	m.Saver[record.Name] = record
	return nil
}

func (m *Memory) Close() error {
	return nil
}
