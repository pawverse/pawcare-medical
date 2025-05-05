package domain

import "fmt"

type (
	PetId    string
	RecordId string
)

type Record struct {
	Id         RecordId
	PetId      PetId
	RecordInfo RecordInfo
}

func NewRecord(petId PetId, recordInfo RecordInfo) *Record {
	return &Record{
		PetId:      petId,
		RecordInfo: recordInfo,
	}
}

func (r Record) String() string {
	return fmt.Sprintf("Record(id=%s,petId=%s,recordInfo=%s)", r.Id, r.PetId, r.RecordInfo)
}
