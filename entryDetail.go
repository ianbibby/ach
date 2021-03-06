// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"strconv"
	"strings"
)

// EntryDetail contains the actual transaction data for an individual entry.
// Fields include those designating the entry as a deposit (credit) or
// withdrawal (debit), the transit routing number for the entry recipient’s financial
// institution, the account number (left justify,no zero fill), name, and dollar amount.
type EntryDetail struct {
	// ID is a client defined string used as a reference to this record.
	ID string `json:"id"`
	// RecordType defines the type of record in the block. 6
	recordType string
	// TransactionCode if the receivers account is:
	// Credit (deposit) to checking account ‘22’
	// Prenote for credit to checking account ‘23’
	// Debit (withdrawal) to checking account ‘27’
	// Prenote for debit to checking account ‘28’
	// Credit to savings account ‘32’
	// Prenote for credit to savings account ‘33’
	// Debit to savings account ‘37’
	// Prenote for debit to savings account ‘38’
	TransactionCode int `json:"transactionCode"`
	// RDFIIdentification is the RDFI's routing number without the last digit.
	// Receiving Depository Financial Institution
	RDFIIdentification string `json:"RDFIIdentification"`
	// CheckDigit the last digit of the RDFI's routing number
	CheckDigit string `json:"checkDigit"`
	// DFIAccountNumber is the receiver's bank account number you are crediting/debiting.
	// It important to note that this is an alphanumeric field, so its space padded, no zero padded
	DFIAccountNumber string `json:"DFIAccountNumber"`
	// Amount Number of cents you are debiting/crediting this account
	Amount int `json:"amount"`
	// IdentificationNumber an internal identification (alphanumeric) that
	// you use to uniquely identify this Entry Detail Record
	IdentificationNumber string `json:"identificationNumber,omitempty"`
	// IndividualName The name of the receiver, usually the name on the bank account
	IndividualName string `json:"individualName"`
	// DiscretionaryData allows ODFIs to include codes, of significance only to them,
	// to enable specialized handling of the entry. There will be no
	// standardized interpretation for the value of this field. It can either
	// be a single two-character code, or two distinct one-character codes,
	// according to the needs of the ODFI and/or Originator involved. This
	// field must be returned intact for any returned entry.
	//
	// WEB uses the Discretionary Data Field as the Payment Type Code
	DiscretionaryData string `json:"discretionaryData,omitempty"`
	// AddendaRecordIndicator indicates the existence of an Addenda Record.
	// A value of "1" indicates that one ore more addenda records follow,
	// and "0" means no such record is present.
	AddendaRecordIndicator int `json:"addendaRecordIndicator,omitempty"`
	// TraceNumber assigned by the ODFI in ascending sequence, is included in each
	// Entry Detail Record, Corporate Entry Detail Record, and addenda Record.
	// Trace Numbers uniquely identify each entry within a batch in an ACH input file.
	// In association with the Batch Number, transmission (File Creation) Date,
	// and File ID Modifier, the Trace Number uniquely identifies an entry within a given file.
	// For addenda Records, the Trace Number will be identical to the Trace Number
	// in the associated Entry Detail Record, since the Trace Number is associated
	// with an entry or item rather than a physical record.
	TraceNumber int `json:"traceNumber,omitempty"`
	// Addendum a list of Addenda for the Entry Detail
	Addendum []Addendumer `json:"addendum,omitempty"`
	// Category defines if the entry is a Forward, Return, or NOC
	Category string `json:"category,omitempty"`
	// validator is composed for data validation
	validator
	// converters is composed for ACH to golang Converters
	converters
}

const (
	// CategoryForward defines the entry as being sent to the receiving institution
	CategoryForward = "Forward"
	// CategoryReturn defines the entry as being a return of a forward entry back to the originating institution
	CategoryReturn = "Return"
	// CategoryNOC defines the entry as being a notification of change of a forward entry to the originating institution
	CategoryNOC = "NOC"
	// ReturnOrNoc is the description for the  following TransactionCode: 21, 31, 41, 51, 26, 36, 46, 56
)

// NewEntryDetail returns a new EntryDetail with default values for non exported fields
func NewEntryDetail() *EntryDetail {
	entry := &EntryDetail{
		recordType: "6",
		Category:   CategoryForward,
	}
	return entry
}

// Parse takes the input record string and parses the EntryDetail values
func (ed *EntryDetail) Parse(record string) {
	// 1-1 Always "6"
	ed.recordType = "6"
	// 2-3 is checking credit 22 debit 27 savings credit 32 debit 37
	ed.TransactionCode = ed.parseNumField(record[1:3])
	// 4-11 the RDFI's routing number without the last digit.
	ed.RDFIIdentification = ed.parseStringField(record[3:11])
	// 12-12 The last digit of the RDFI's routing number
	ed.CheckDigit = ed.parseStringField(record[11:12])
	// 13-29 The receiver's bank account number you are crediting/debiting
	ed.DFIAccountNumber = record[12:29]
	// 30-39 Number of cents you are debiting/crediting this account
	ed.Amount = ed.parseNumField(record[29:39])
	// 40-54 An internal identification (alphanumeric) that you use to uniquely identify this Entry Detail Record
	ed.IdentificationNumber = record[39:54]
	// 55-76 The name of the receiver, usually the name on the bank account
	ed.IndividualName = record[54:76]
	// 77-78 allows ODFIs to include codes of significance only to them
	// For WEB transaction this field is the PaymentType which is either R(reoccurring) or S(single)
	// normally blank
	ed.DiscretionaryData = record[76:78]
	// 79-79 1 if addenda exists 0 if it does not
	ed.AddendaRecordIndicator = ed.parseNumField(record[78:79])
	// 80-94 An internal identification (alphanumeric) that you use to uniquely identify
	// this Entry Detail Record This number should be unique to the transaction and will help identify the transaction in case of an inquiry
	ed.TraceNumber = ed.parseNumField(record[79:94])
}

// String writes the EntryDetail struct to a 94 character string.
func (ed *EntryDetail) String() string {
	var buf strings.Builder
	buf.Grow(94)
	buf.WriteString(ed.recordType)
	buf.WriteString(fmt.Sprintf("%v", ed.TransactionCode))
	buf.WriteString(ed.RDFIIdentificationField())
	buf.WriteString(ed.CheckDigit)
	buf.WriteString(ed.DFIAccountNumberField())
	buf.WriteString(ed.AmountField())
	buf.WriteString(ed.IdentificationNumberField())
	buf.WriteString(ed.IndividualNameField())
	buf.WriteString(ed.DiscretionaryDataField())
	buf.WriteString(fmt.Sprintf("%v", ed.AddendaRecordIndicator))
	buf.WriteString(ed.TraceNumberField())
	return buf.String()
}

// Validate performs NACHA format rule checks on the record and returns an error if not Validated
// The first error encountered is returned and stops that parsing.
func (ed *EntryDetail) Validate() error {
	if err := ed.fieldInclusion(); err != nil {
		return err
	}
	if ed.recordType != "6" {
		msg := fmt.Sprintf(msgRecordType, 6)
		return &FieldError{FieldName: "recordType", Value: ed.recordType, Msg: msg}
	}
	if err := ed.isTransactionCode(ed.TransactionCode); err != nil {
		return &FieldError{FieldName: "TransactionCode", Value: strconv.Itoa(ed.TransactionCode), Msg: err.Error()}
	}
	if err := ed.isAlphanumeric(ed.DFIAccountNumber); err != nil {
		return &FieldError{FieldName: "DFIAccountNumber", Value: ed.DFIAccountNumber, Msg: err.Error()}
	}
	if err := ed.isAlphanumeric(ed.IdentificationNumber); err != nil {
		return &FieldError{FieldName: "IdentificationNumber", Value: ed.IdentificationNumber, Msg: err.Error()}
	}
	if err := ed.isAlphanumeric(ed.IndividualName); err != nil {
		return &FieldError{FieldName: "IndividualName", Value: ed.IndividualName, Msg: err.Error()}
	}
	if err := ed.isAlphanumeric(ed.DiscretionaryData); err != nil {
		return &FieldError{FieldName: "DiscretionaryData", Value: ed.DiscretionaryData, Msg: err.Error()}
	}

	calculated := ed.CalculateCheckDigit(ed.RDFIIdentificationField())

	edCheckDigit, err := strconv.Atoi(ed.CheckDigit)
	if err != nil {
		return &FieldError{FieldName: "CheckDigit", Value: ed.CheckDigit, Msg: err.Error()}
	}

	if calculated != edCheckDigit {
		msg := fmt.Sprintf(msgValidCheckDigit, calculated)
		return &FieldError{FieldName: "RDFIIdentification", Value: ed.CheckDigit, Msg: msg}
	}
	return nil
}

// fieldInclusion validate mandatory fields are not default values. If fields are
// invalid the ACH transfer will be returned.
func (ed *EntryDetail) fieldInclusion() error {
	if ed.recordType == "" {
		return &FieldError{FieldName: "recordType", Value: ed.recordType, Msg: msgFieldInclusion}
	}
	if ed.TransactionCode == 0 {
		return &FieldError{FieldName: "TransactionCode", Value: strconv.Itoa(ed.TransactionCode), Msg: msgFieldInclusion}
	}
	if ed.RDFIIdentification == "" {
		return &FieldError{FieldName: "RDFIIdentification", Value: ed.RDFIIdentificationField(), Msg: msgFieldInclusion}
	}
	if ed.DFIAccountNumber == "" {
		return &FieldError{FieldName: "DFIAccountNumber", Value: ed.DFIAccountNumber, Msg: msgFieldInclusion}
	}
	if ed.IndividualName == "" {
		return &FieldError{FieldName: "IndividualName", Value: ed.IndividualName, Msg: msgFieldInclusion}
	}
	if ed.TraceNumber == 0 {
		return &FieldError{FieldName: "TraceNumber", Value: ed.TraceNumberField(), Msg: msgFieldInclusion}
	}
	return nil
}

// AddAddenda appends an Addendumer to the EntryDetail
//
// Note: The order of records here is determined by their insertion order.
// No inspection of SequenceNumbers in Addendas (i.e 05, 17, 18) is done
// to re-order addenda records.
func (ed *EntryDetail) AddAddenda(addenda Addendumer) []Addendumer {
	ed.AddendaRecordIndicator = 1
	// checks to make sure that we only have either or, not both
	switch addenda.(type) {
	case *Addenda99:
		ed.Category = CategoryReturn
		ed.Addendum = nil
		ed.Addendum = append(ed.Addendum, addenda)
		return ed.Addendum
	case *Addenda98:
		ed.Category = CategoryNOC
		ed.Addendum = nil
		ed.Addendum = append(ed.Addendum, addenda)
		return ed.Addendum
	case *Addenda02:
		ed.Category = CategoryForward
		ed.Addendum = nil
		ed.Addendum = append(ed.Addendum, addenda)
		return ed.Addendum
		// default is current *Addenda05
	default:
		ed.Category = CategoryForward
		ed.Addendum = append(ed.Addendum, addenda)
		return ed.Addendum
	}
}

// SetRDFI takes the 9 digit RDFI account number and separates it for RDFIIdentification and CheckDigit
func (ed *EntryDetail) SetRDFI(rdfi string) *EntryDetail {
	s := ed.stringField(rdfi, 9)
	ed.RDFIIdentification = ed.parseStringField(s[:8])
	ed.CheckDigit = ed.parseStringField(s[8:9])
	return ed
}

// SetTraceNumber takes first 8 digits of ODFI and concatenates a sequence number onto the TraceNumber
func (ed *EntryDetail) SetTraceNumber(ODFIIdentification string, seq int) {
	trace := ed.stringField(ODFIIdentification, 8) + ed.numericField(seq, 7)
	ed.TraceNumber = ed.parseNumField(trace)
}

// RDFIIdentificationField get the rdfiIdentification with zero padding
func (ed *EntryDetail) RDFIIdentificationField() string {
	return ed.stringField(ed.RDFIIdentification, 8)
}

// DFIAccountNumberField gets the DFIAccountNumber with space padding
func (ed *EntryDetail) DFIAccountNumberField() string {
	return ed.alphaField(ed.DFIAccountNumber, 17)
}

// AmountField returns a zero padded string of amount
func (ed *EntryDetail) AmountField() string {
	return ed.numericField(ed.Amount, 10)
}

// IdentificationNumberField returns a space padded string of IdentificationNumber
func (ed *EntryDetail) IdentificationNumberField() string {
	return ed.alphaField(ed.IdentificationNumber, 15)
}

// CheckSerialNumberField is used in RCK, ARC, BOC files but returns
// a space padded string of the underlying IdentificationNumber field
func (ed *EntryDetail) CheckSerialNumberField() string {
	return ed.alphaField(ed.IdentificationNumber, 15)
}

// SetCheckSerialNumber setter for RCK, ARC, BOC CheckSerialNumber
// which is underlying IdentificationNumber
func (ed *EntryDetail) SetCheckSerialNumber(s string) {
	ed.IdentificationNumber = s
}

// SetPOPCheckSerialNumber setter for POP CheckSerialNumber
// which is characters 1-9 of underlying CheckSerialNumber \ IdentificationNumber
func (ed *EntryDetail) SetPOPCheckSerialNumber(s string) {
	ed.IdentificationNumber = ed.alphaField(s, 9)
}

// SetPOPTerminalCity setter for POP Terminal City
// which is characters 10-13 of underlying CheckSerialNumber \ IdentificationNumber
func (ed *EntryDetail) SetPOPTerminalCity(s string) {
	ed.IdentificationNumber = ed.IdentificationNumber + ed.alphaField(s, 4)
}

// SetPOPTerminalState setter for POP Terminal State
// which is characters 14-15 of underlying CheckSerialNumber \ IdentificationNumber
func (ed *EntryDetail) SetPOPTerminalState(s string) {
	ed.IdentificationNumber = ed.IdentificationNumber + ed.alphaField(s, 2)
}

// POPCheckSerialNumberField is used in POP, characters 1-9 of underlying BatchPOP
// CheckSerialNumber / IdentificationNumber
func (ed *EntryDetail) POPCheckSerialNumberField() string {
	return ed.parseStringField(ed.IdentificationNumber[0:9])
}

// POPTerminalCityField is used in POP, characters 10-13 of underlying BatchPOP
// CheckSerialNumber / IdentificationNumber
func (ed *EntryDetail) POPTerminalCityField() string {
	return ed.parseStringField(ed.IdentificationNumber[9:13])
}

// POPTerminalStateField is used in POP, characters 14-15 of underlying BatchPOP
// CheckSerialNumber / IdentificationNumber
func (ed *EntryDetail) POPTerminalStateField() string {
	return ed.parseStringField(ed.IdentificationNumber[13:15])
}

// SetSHRCardExpirationDate format MMYY is used in SHR, characters 1-4 of underlying
// IdentificationNumber
func (ed *EntryDetail) SetSHRCardExpirationDate(s string) {
	ed.IdentificationNumber = ed.alphaField(s, 4)
}

// SetSHRDocumentReferenceNumber format int is used in SHR, characters 5-15 of underlying
// IdentificationNumber
func (ed *EntryDetail) SetSHRDocumentReferenceNumber(s string) {
	ed.IdentificationNumber = ed.IdentificationNumber + ed.stringField(s, 11)
}

// SetSHRIndividualCardAccountNumber format int is used in SHR, underlying
// IndividualName
func (ed *EntryDetail) SetSHRIndividualCardAccountNumber(s string) {
	ed.IndividualName = ed.stringField(s, 22)
}

// SHRCardExpirationDateField format MMYY is used in SHR, characters 1-4 of underlying
// IdentificationNumber
func (ed *EntryDetail) SHRCardExpirationDateField() string {
	return ed.parseStringField(ed.IdentificationNumber[0:4])
}

// SHRDocumentReferenceNumberField format int is used in SHR, characters 5-15 of underlying
// IdentificationNumber
func (ed *EntryDetail) SHRDocumentReferenceNumberField() string {
	return ed.stringField(ed.IdentificationNumber[4:15], 11)
}

// SHRIndividualCardAccountNumberField format int is used in SHR, underlying
// IndividualName
func (ed *EntryDetail) SHRIndividualCardAccountNumberField() string {
	return ed.stringField(ed.IndividualName, 22)
}

// IndividualNameField returns a space padded string of IndividualName
func (ed *EntryDetail) IndividualNameField() string {
	return ed.alphaField(ed.IndividualName, 22)
}

// ReceivingCompanyField is used in CCD files but returns the underlying IndividualName field
func (ed *EntryDetail) ReceivingCompanyField() string {
	return ed.IndividualNameField()
}

// SetReceivingCompany setter for CCD ReceivingCompany which is underlying IndividualName
func (ed *EntryDetail) SetReceivingCompany(s string) {
	ed.IndividualName = s
}

// SetCTXAddendaRecords setter for CTX AddendaRecords characters 1-4 of underlying IndividualName
func (ed *EntryDetail) SetCTXAddendaRecords(i int) {
	ed.IndividualName = ed.numericField(i, 4)
}

// SetCTXReceivingCompany setter for CTX ReceivingCompany characters 5-20 underlying IndividualName
// Position 21-22 of underlying Individual Name are reserved blank space for CTX "  "
func (ed *EntryDetail) SetCTXReceivingCompany(s string) {
	ed.IndividualName = ed.IndividualName + ed.alphaField(s, 16) + "  "
}

// CTXAddendaRecordsField is used in CTX files, characters 1-4 of underlying IndividualName field
func (ed *EntryDetail) CTXAddendaRecordsField() string {
	return ed.parseStringField(ed.IndividualName[0:4])
}

// CTXReceivingCompanyField is used in CTX files, characters 5-20 of underlying IndividualName field
func (ed *EntryDetail) CTXReceivingCompanyField() string {
	return ed.parseStringField(ed.IndividualName[4:20])
}

// CTXReservedField is used in CTX files, characters 21-22 of underlying IndividualName field
func (ed *EntryDetail) CTXReservedField() string {
	return ed.IndividualName[20:22]
}

// DiscretionaryDataField returns a space padded string of DiscretionaryData
func (ed *EntryDetail) DiscretionaryDataField() string {
	return ed.alphaField(ed.DiscretionaryData, 2)
}

// PaymentTypeField returns the DiscretionaryData field used in WEB batch files
func (ed *EntryDetail) PaymentTypeField() string {
	// because DiscretionaryData can be changed outside of PaymentType we reset the value for safety
	ed.SetPaymentType(ed.DiscretionaryData)
	return ed.DiscretionaryData
}

// SetPaymentType as R (Recurring) all other values will result in S (single)
func (ed *EntryDetail) SetPaymentType(t string) {
	t = strings.ToUpper(strings.TrimSpace(t))
	if t == "R" {
		ed.DiscretionaryData = "R"
	} else {
		ed.DiscretionaryData = "S"
	}
}

// TraceNumberField returns a zero padded TraceNumber string
func (ed *EntryDetail) TraceNumberField() string {
	return ed.numericField(ed.TraceNumber, 15)
}

// CreditOrDebit returns a "C" for credit or "D" for debit based on the entry TransactionCode
func (ed *EntryDetail) CreditOrDebit() string {
	tc := strconv.Itoa(ed.TransactionCode)
	// take the second number in the TransactionCode
	switch tc[1:2] {
	case "1", "2", "3", "4":
		return "C"
	case "5", "6", "7", "8", "9":
		return "D"
	default:
	}
	return ""
}
