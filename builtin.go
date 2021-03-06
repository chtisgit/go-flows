package main

import (
	_ "github.com/chtisgit/go-flows/modules/exporters/csv"
	_ "github.com/chtisgit/go-flows/modules/exporters/ipfix"
	_ "github.com/chtisgit/go-flows/modules/exporters/null"
	_ "github.com/chtisgit/go-flows/modules/exporters/sql"
	_ "github.com/chtisgit/go-flows/modules/features/custom"
	_ "github.com/chtisgit/go-flows/modules/features/iana"
	_ "github.com/chtisgit/go-flows/modules/features/nta"
	_ "github.com/chtisgit/go-flows/modules/features/operations"
	_ "github.com/chtisgit/go-flows/modules/features/staging"
	_ "github.com/chtisgit/go-flows/modules/filters/time"
	_ "github.com/chtisgit/go-flows/modules/keys/header"
	_ "github.com/chtisgit/go-flows/modules/keys/time"
	_ "github.com/chtisgit/go-flows/modules/labels/csv"
	_ "github.com/chtisgit/go-flows/modules/sources/libpcap"
)
