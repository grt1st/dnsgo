package parser

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	g "github.com/grt1st/dnsgo/global"
	"github.com/boltdb/bolt"
	"log"
)


func getRecord(domain string, rType uint16) (dns.RR, error) {
	rA := new(dns.A)

	address, ok := g.DB.GetRecord(domain)
	if ok {
		fmt.Println("ok")
		rHeader := dns.RR_Header{
			Name:   domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60,
		}
		rA.Hdr = rHeader
		rA.A = net.ParseIP(address)
	}

	return rA, nil
}


func deleteRecord(domain string, rtype uint16) (err error) {
	key, _ := getKey(domain, rtype)
	err = bdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(rr_bucket))
		err := b.Delete([]byte(key))

		if err != nil {
			e := errors.New( "Delete record failed for domain:  " + domain)
			log.Println(e.Error())

			return e
		}

		return nil
	})

	return err
}