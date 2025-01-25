package service

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/Pinkman-77/records-restapi/pkg/repository"
)

type RecordList struct {
	repo repository.Record
}

func NewRecordList(repo repository.Record) *RecordList {
        return &RecordList{repo: repo}
}

func (r *RecordList) CreateRecord(record recordsrestapi.Record) (int, error) {
	return r.repo.CreateRecord(record)
}

func (r *RecordList) GetAllRecords() ([]recordsrestapi.Record, error) {
        return r.repo.GetAllRecords()
}

func (r *RecordList) GetRecord(id int) (recordsrestapi.RecordWithArtist, error) {
        return r.repo.GetRecord(id)
}

func (r *RecordList) UpdateRecord(id int, record recordsrestapi.Record) error {
	return r.repo.UpdateRecord(id, record)
}

func (r *RecordList) PatchRecord(id int, updates map[string]interface{}) error {
	return r.repo.PatchRecord(id, updates)
}

func (r *RecordList) DeleteRecord(id int) error {
	return r.repo.DeleteRecord(id)
}
