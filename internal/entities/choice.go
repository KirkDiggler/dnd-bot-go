package entities

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
	OptionTypeMultiple         OptionType = "multiple"
	OptionTypeCountedReference OptionType = "counted_reference"
)

type ChoiceStatus string

const (
	ChoiceStatusUnset    ChoiceStatus = ""
	ChoiceStatusActive   ChoiceStatus = "active"
	ChoiceStatusInactive ChoiceStatus = "inactive"
	ChoiceStatusSelected ChoiceStatus = "selected"
)

type ChoiceType string

const (
	ChoiceTypeUnset       ChoiceType = ""
	ChoiceTypeProficiency ChoiceType = "proficiency"
	ChoiceTypeLanguage    ChoiceType = "language"
	ChoiceTypeEquipment   ChoiceType = "equipment"
)

type Choice struct {
	Name     string       `json:"name"`
	Type     ChoiceType   `json:"type"`
	Key      string       `json:"key"`
	Status   ChoiceStatus `json:"status"`
	Selected bool         `json:"selected"`
	Count    int          `json:"count"`
	Options  []Option     `json:"options"`
}
type Option interface {
	GetOptionType() OptionType
	GetName() string
	GetKey() string
	GetStatus() ChoiceStatus
	SetStatus(ChoiceStatus)
}

type CountedReferenceOption struct {
	Status    ChoiceStatus   `json:"status"`
	Count     int            `json:"count"`
	Reference *ReferenceItem `json:"reference"`
}

func (o *CountedReferenceOption) GetOptionType() OptionType {
	return OptionTypeCountedReference
}

func (o *CountedReferenceOption) GetName() string {
	return o.Reference.Name
}

func (o *CountedReferenceOption) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *CountedReferenceOption) SetStatus(status ChoiceStatus) {
	o.Status = status
}

func (o *CountedReferenceOption) GetKey() string {
	return o.Reference.Key
}

type ReferenceOption struct {
	Status    ChoiceStatus   `json:"status"`
	Reference *ReferenceItem `json:"reference"`
}

func (o *ReferenceOption) GetOptionType() OptionType {
	return OptionTypeReference
}

func (o *ReferenceOption) GetName() string {
	return o.Reference.Name
}

func (o *ReferenceOption) GetKey() string {
	return o.Reference.Key
}

func (o *ReferenceOption) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *ReferenceOption) SetStatus(status ChoiceStatus) {
	o.Status = status
}

func (o *Choice) GetOptionType() OptionType {
	return OptionTypeChoice
}

func (o *Choice) GetName() string {
	return o.Name
}

func (o *Choice) GetKey() string {
	return o.Key
}

func (o *Choice) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *Choice) SetStatus(status ChoiceStatus) {
	o.Status = status
}

type MultipleOption struct {
	Status ChoiceStatus `json:"status"`
	Items  []Option     `json:"items"`
}

func (o *MultipleOption) GetOptionType() OptionType {
	return OptionTypeMultiple
}

func (o *MultipleOption) GetName() string {
	return ""
}

func (o *MultipleOption) GetKey() string {
	return ""
}

func (o *MultipleOption) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *MultipleOption) SetStatus(status ChoiceStatus) {
	o.Status = status
}
