package internal

import "github.com/oschwald/geoip2-golang"

type ipDb struct {
	Db *geoip2.Reader
}

var IpDb = ipDb{}

func (i *ipDb) GetIpDbInstance() *geoip2.Reader {
	var err error
	if i.Db == nil {
		i.Db, err = geoip2.Open("data/country.mmdb")
		if err != nil {
			return nil
		}
	}
	return i.Db
}

func (i *ipDb) CloseIpDbInstance() {
	if i.Db != nil {
		_ = i.Db.Close()
	}
}
