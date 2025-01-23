package service

import (
        "github.com/Pinkman-77/records-restapi/repository"
)

type RecordList struct {
        repo repository.Record
}

func NewRecordList(repo repository.Record) *RecordList {
        return &RecordList{repo: repo}
}


