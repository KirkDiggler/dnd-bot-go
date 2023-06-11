package combat_log

type ListByEncounterInput struct {
	EncounterID string
	Limit       int64
	Offset      int64
	Reverse     bool
}
