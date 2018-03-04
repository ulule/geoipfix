package ipfix

import "github.com/ulule/ipfix/proto"

func recordToProto(r *record) *proto.Location {
	loc := &proto.Location{
		IpAddress: r.IP,
		City:      r.City,
		ZipCode:   r.ZipCode,
		TimeZone:  r.TimeZone,
		Longitude: float32(r.Longitude),
		Latitude:  float32(r.Latitude),
	}

	if r.CountryCode != "" || r.CountryName != "" {
		loc.Country = &proto.Place{
			Code: r.CountryCode,
			Name: r.CountryName,
		}
	}

	if r.RegionCode != "" || r.RegionName != "" {
		loc.Region = &proto.Place{
			Code: r.RegionCode,
			Name: r.RegionName,
		}
	}

	return loc
}
