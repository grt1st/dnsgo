package backends

type Backend interface {
	Init() (error)
	getValue() (error)
	DeleteRecord(domain string) (error)
	SaveRecord(domain, address string) (error)
	GetRecord(domain string) (string, bool)
	UpdateRecord(domain, address string) (error)
	Close() (error)
}
