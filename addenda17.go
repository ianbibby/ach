// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"strings"
)

// Addenda17 is an addenda which provides business transaction information for Addenda Type
// Code 17 in a machine readable format. It is usually formatted according to ANSI, ASC, X12 Standard.
//
// Addenda17 is optional for IAT entries
//
// The Addenda17 record identifies payment-related data. A maximum of two of these Addenda Records
// may be included with each IAT entry.
type Addenda17 struct {
	// ID is a client defined string used as a reference to this record.
	ID string `json:"id"`
	// RecordType defines the type of record in the block. entryAddenda17 Pos 7
	recordType string
	// TypeCode Addenda17 types code '17'
	typeCode string
	// PaymentRelatedInformation
	PaymentRelatedInformation string `json:"paymentRelatedInformation"`
	// SequenceNumber is consecutively assigned to each Addenda17 Record following
	// an Entry Detail Record. The first addenda17 sequence number must always
	// be a "1".
	SequenceNumber int `json:"sequenceNumber,omitempty"`
	// EntryDetailSequenceNumber contains the ascending sequence number section of the Entry
	// Detail or Corporate Entry Detail Record's trace number This number is
	// the same as the last seven digits of the trace number of the related
	// Entry Detail Record or Corporate Entry Detail Record.
	EntryDetailSequenceNumber int `json:"entryDetailSequenceNumber,omitempty"`
	// validator is composed for data validation
	validator
	// converters is composed for ACH to GoLang Converters
	converters
}

// NewAddenda17 returns a new Addenda17 with default values for none exported fields
func NewAddenda17() *Addenda17 {
	addenda17 := new(Addenda17)
	addenda17.recordType = "7"
	addenda17.typeCode = "17"
	return addenda17
}

// Parse takes the input record string and parses the Addenda17 values
func (addenda17 *Addenda17) Parse(record string) {
	// 1-1 Always "7"
	addenda17.recordType = "7"
	// 2-3 Always 17
	addenda17.typeCode = record[1:3]
	// 4-83 Based on the information entered (04-83) 80 alphanumeric
	addenda17.PaymentRelatedInformation = strings.TrimSpace(record[3:83])
	// 84-87 SequenceNumber is consecutively assigned to each Addenda17 Record following
	// an Entry Detail Record
	addenda17.SequenceNumber = addenda17.parseNumField(record[83:87])
	// 88-94 Contains the last seven digits of the number entered in the Trace Number field in the corresponding Entry Detail Record
	addenda17.EntryDetailSequenceNumber = addenda17.parseNumField(record[87:94])
}

// String writes the Addenda17 struct to a 94 character string.
func (addenda17 *Addenda17) String() string {
	var buf strings.Builder
	buf.Grow(94)
	buf.WriteString(addenda17.recordType)
	buf.WriteString(addenda17.typeCode)
	buf.WriteString(addenda17.PaymentRelatedInformationField())
	buf.WriteString(addenda17.SequenceNumberField())
	buf.WriteString(addenda17.EntryDetailSequenceNumberField())
	return buf.String()
}

// Validate performs NACHA format rule checks on the record and returns an error if not Validated
// The first error encountered is returned and stops that parsing.
func (addenda17 *Addenda17) Validate() error {
	if err := addenda17.fieldInclusion(); err != nil {
		return err
	}
	if addenda17.recordType != "7" {
		msg := fmt.Sprintf(msgRecordType, 7)
		return &FieldError{FieldName: "recordType", Value: addenda17.recordType, Msg: msg}
	}
	if err := addenda17.isTypeCode(addenda17.typeCode); err != nil {
		return &FieldError{FieldName: "TypeCode", Value: addenda17.typeCode, Msg: err.Error()}
	}
	// Type Code must be 17
	if addenda17.typeCode != "17" {
		return &FieldError{FieldName: "TypeCode", Value: addenda17.typeCode, Msg: msgAddendaTypeCode}
	}
	if err := addenda17.isAlphanumeric(addenda17.PaymentRelatedInformation); err != nil {
		return &FieldError{FieldName: "PaymentRelatedInformation", Value: addenda17.PaymentRelatedInformation, Msg: err.Error()}
	}

	return nil
}

// fieldInclusion validate mandatory fields are not default values. If fields are
// invalid the ACH transfer will be returned.
func (addenda17 *Addenda17) fieldInclusion() error {
	if addenda17.recordType == "" {
		return &FieldError{FieldName: "recordType", Value: addenda17.recordType, Msg: msgFieldInclusion}
	}
	if addenda17.typeCode == "" {
		return &FieldError{FieldName: "TypeCode", Value: addenda17.typeCode, Msg: msgFieldInclusion}
	}
	if addenda17.SequenceNumber == 0 {
		return &FieldError{FieldName: "SequenceNumber", Value: addenda17.SequenceNumberField(), Msg: msgFieldInclusion}
	}
	if addenda17.EntryDetailSequenceNumber == 0 {
		return &FieldError{FieldName: "EntryDetailSequenceNumber", Value: addenda17.EntryDetailSequenceNumberField(), Msg: msgFieldInclusion}
	}
	return nil
}

// PaymentRelatedInformationField returns a zero padded PaymentRelatedInformation string
func (addenda17 *Addenda17) PaymentRelatedInformationField() string {
	return addenda17.alphaField(addenda17.PaymentRelatedInformation, 80)
}

// SequenceNumberField returns a zero padded SequenceNumber string
func (addenda17 *Addenda17) SequenceNumberField() string {
	return addenda17.numericField(addenda17.SequenceNumber, 4)
}

// EntryDetailSequenceNumberField returns a zero padded EntryDetailSequenceNumber string
func (addenda17 *Addenda17) EntryDetailSequenceNumberField() string {
	return addenda17.numericField(addenda17.EntryDetailSequenceNumber, 7)
}

// TypeCode Defines the specific explanation and format for the addenda17 information
func (addenda17 *Addenda17) TypeCode() string {
	return addenda17.typeCode
}
