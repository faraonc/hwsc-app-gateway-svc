package service

import (
	pbauth "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"time"
)

var (
	validCreatedTimestamp    = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	validExpirationTimestamp = time.Unix(validCreatedTimestamp, 0).AddDate(30, 0, 0).UTC().Unix()
	validAuthSecretKey       = "j2Yzh-VcIm-lYUzBuqt8TVPeUHNYB5MP1gWvz3Bolow="

	validAuthSecret = &pbauth.Secret{
		Key:                 validAuthSecretKey,
		CreatedTimestamp:    validCreatedTimestamp,
		ExpirationTimestamp: validExpirationTimestamp,
	}

	expiredAuthSecret = &pbauth.Secret{
		Key:                 validAuthSecretKey,
		CreatedTimestamp:    time.Unix(validCreatedTimestamp, 0).AddDate(-1, 0, 0).UTC().Unix(),
		ExpirationTimestamp: validExpirationTimestamp,
	}
)