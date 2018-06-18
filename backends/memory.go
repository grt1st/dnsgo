package backends

import "sync"

type Memory struct {
	Saver map[string]string
	//Rtype dns.Type
	//Class dns.Class
	//Ttl int
	sync.RWMutex
}

func (m *Memory) Init() (error) {
	m.Saver = map[string]string{
		"google.com.": "1.2.3.4",
		"ja.com.":     "104.198.14.52",
		"grt1st.cn.":  "123.206.60.140",
	}
	return nil
}

func (m *Memory) getValue() (error) {
	//value, ok := m.Saver[domain]
	return nil
}

func (m *Memory) DeleteRecord(domain string) (error){
	m.Lock()
	defer m.Unlock()
	delete(m.Saver, domain)
	return nil
}

func (m *Memory) SaveRecord(domain, address string) (error) {
	m.Lock()
	defer m.Unlock()
	m.Saver[domain] = address
	return nil
}

func (m *Memory) GetRecord(domain string) (string, bool) {
	m.Lock()
	defer m.Unlock()
	address, ok := m.Saver[domain]
	return address, ok
}

func (m *Memory) UpdateRecord(domain, address string) (error) {
	m.Lock()
	defer m.Unlock()
	m.Saver["domain"] = address
	return nil
}

func (m *Memory) Close() (error) {
	return nil
}