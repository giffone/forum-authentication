package datef

import (
	"time"
)

func Expire(checkDate time.Time) bool {
	today := time.Now()
	//log.Println(today)
	//log.Println(checkDate)
	if today.After(checkDate) && today != checkDate {
		return true
	}
	return false
}

func ExpireF(checkDate, checkDateFormat string) (bool, error) {
	d, err := Format(checkDate, checkDateFormat)
	if err != nil {
		return true, err
	}
	today := time.Now()
	if today.After(d) && today != d {
		return true, nil
	}
	return false, nil
}
