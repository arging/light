// Copyright 2014 li. All rights reserved.

package session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/roverli/utils/errors"
	"strconv"
	"time"
)

func Md5Id(sed string) (string, errors.Error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "light/session: Rand package error.")
	}

	timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	id := sed + string(b) + timeStr

	h := md5.New()
	h.Write([]byte(id))

	return hex.EncodeToString(h.Sum(nil)), nil
}
