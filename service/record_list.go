package service

import (
        recordsrestapi "github.com/Pinkman-77/records-restapi"
        "github.com/Pinkman-77/records-restapi/repository"
)

type RecordList struct {
        repo repository.Record
}

func NewRecordList(repo repository.Record) *RecordList {
        return &RecordList{repo: repo}
}

func (r *RecordList) CreateRecord(record recordsrestapi.Record) (recordsrestapi.Record, error) {
        return r.repo.CreateRecord(record)
}

