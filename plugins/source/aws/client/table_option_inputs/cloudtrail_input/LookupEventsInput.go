// Code generated by generator, DO NOT EDIT.
package cloudtrail_input

import "time"

type EventCategory string
type LookupAttributeKey string
type LookupAttribute struct {
	AttributeKey   *LookupAttributeKey `json:"attribute_key,omitempty"`
	AttributeValue *string             `json:"attribute_value,omitempty"`
}
type LookupEventsInput struct {
	EndTime          *time.Time        `json:"end_time,omitempty"`
	EventCategory    *EventCategory    `json:"event_category,omitempty"`
	LookupAttributes []LookupAttribute `json:"lookup_attributes,omitempty"`
	MaxResults       *int32            `json:"max_results,omitempty"`
	NextToken        *string           `json:"next_token,omitempty"`
	StartTime        *time.Time        `json:"start_time,omitempty"`
}
