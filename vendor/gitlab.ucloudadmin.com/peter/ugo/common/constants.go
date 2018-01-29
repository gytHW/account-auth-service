package common

import "time"

const (
	ETCD_DIAL_TIMEOUT = time.Minute * 5
)


const (

	COMMON_OK                      = iota
	COMMON_ERROR

	COMMON_ERROR_DB                = 100
	COMMON_ERROR_MONGO_CONNECT

	COMMON_ERROR_SERVICE           = 200
	COMMON_ERROR_SERVICE_NOT_FOUND


	COMMON_ERROR_MAX               = 10000
)
