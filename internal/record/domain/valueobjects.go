package domain

import (
	"fmt"
	"time"
)

type RecordInfo struct {
	Type        RecordType
	Description string
	Date        time.Time
}

func NewRecordInfo(recordType RecordType, description string, date time.Time) RecordInfo {
	return RecordInfo{
		Type:        recordType,
		Description: description,
		Date:        date,
	}
}

func (r RecordInfo) String() string {
	return fmt.Sprintf("RecordInfo(Type=%s,Description=%s,Date=%s)", r.Type, r.Description, r.Date)
}
