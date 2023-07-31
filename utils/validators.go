package utils

import (
	"errors"
	"fmt"
)

func VerifyUrlParamsNullabity(strs ...string) error {
	for _, s := range strs {
		if s == "" {
			strError := fmt.Sprintf("resquest needs to contains: {%v}", s)
			return errors.New(strError)
		}
	}

	return nil
}
