package main

import "time"

type aws_access_key struct {
	AccessKeyId string
	CreateDate  time.Time
	Status      string
	Username    string
}
