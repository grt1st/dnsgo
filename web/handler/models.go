package handler

import "github.com/grt1st/dnsgo/backends"

type ListParams struct {
	Keywords string `query:"keyword"`
	Limit    int    `query:"limit"`
	Offset   int    `query:"offset"`
}

type DNSConfigRecord struct {
	ID    uint   `json:"id" query:"id"`
	Key   string `json:"key" query:"key"`
	Value string `json:"value" query:"value"`
	Kind  string `json:"kind" query:"value"`
}

func (req DNSConfigRecord) ToDBStruct() *backends.DNSConfig {
	return &backends.DNSConfig{
		Name:  req.Key,
		Value: req.Value,
		Kind:  req.Kind,
	}
}

type DNSLookupRecord struct {
	ID        uint   `json:"id" query:"id"`
	CreatedAt string `json:"created_at"`
	IP        string `json:"ip"`
	Key       string `json:"key" query:"key"`
	Value     string `json:"value" query:"value"`
	Kind      string `json:"kind" query:"value"`
}

type HTTPRecord struct {
	ID         uint              `json:"id" query:"id"`
	CreatedAt  string            `json:"created_at"`
	IP         string            `json:"ip"`
	ReqHeaders map[string]string `json:"req_headers"`
	ReqBody    string            `json:"req_body"`
	URL        string            `json:"url"`
}
