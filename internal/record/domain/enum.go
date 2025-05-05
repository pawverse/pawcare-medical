package domain

type RecordType string

const (
	RecordTypeVaccination RecordType = "vaccination"
	RecordTypeTreatment   RecordType = "treatment"
	RecordTypeDeworming   RecordType = "deworming"
	RecordTypeSurgery     RecordType = "surgery"
	RecordTypeCheckup     RecordType = "checkup"
	RecordTypeOther       RecordType = "other"
)
