package datef

import (
	"log"
	"time"
)

func FormatF(date, dateFormat, needFormat string) (string, error) {
	d, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("date formatF: %v", err)
		return "", err
	}
	return d.Format(needFormat), nil
}

func Format(date, dateFormat string) (time.Time, error) {
	d, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("date format: %v", err)
		return time.Time{}, err
	}
	return d, nil
}
