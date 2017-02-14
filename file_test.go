// Copyright 2017 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"strings"
	"testing"
)

// TestBatchCountError tests for to many batch counts
func TestBatchCountError(t *testing.T) {
	r := NewReader(strings.NewReader(" "))
	r.file.addBatch(Batch{})
	r.file.Control.BatchCount = 1
	if err := r.file.Validate(); err != nil {
		t.Errorf("Unexpected File.Validation error: %v", err.Error())
	}
	// More batches than the file control count.
	r.file.addBatch(Batch{})
	if err := r.file.Validate(); err != nil {
		if err != ErrFileBatchCount {
			t.Errorf("Unexpected File.Validation error: %v", err.Error())
		}
	}
}

func TestFileEntryAddendaError(t *testing.T) {
	r := NewReader(strings.NewReader(" "))
	mockBatch := Batch{}
	mockBatch.Control.EntryAddendaCount = 1
	r.file.addBatch(mockBatch)
	r.file.Control.BatchCount = 1
	r.file.Control.EntryAddendaCount = 1
	if err := r.file.Validate(); err != nil {
		t.Errorf("Unexpected File.Validation error: %v", err.Error())
	}

	// more entries than the file control
	r.file.Control.EntryAddendaCount = 5
	if err := r.file.Validate(); err != nil {
		if err != ErrFileEntryCount {
			t.Errorf("Unexpected File.Validation error: %v", err.Error())
		}
	}
}

func TestFileDebitAmount(t *testing.T) {

	r := NewReader(strings.NewReader(" "))
	mockBatch := Batch{}
	mockBatch.Control.EntryAddendaCount = 1
	mockBatch.Control.TotalDebitEntryDollarAmount = 10500

	r.file.addBatch(mockBatch)
	r.file.Control.BatchCount = 1
	r.file.Control.EntryAddendaCount = 1
	r.file.Control.TotalDebitEntryDollarAmountInFile = 10500

	if err := r.file.Validate(); err != nil {
		t.Errorf("Unexpected File.Validation error: %v", err.Error())
	}

	r.file.Control.TotalDebitEntryDollarAmountInFile = 0
	if err := r.file.Validate(); err != nil {
		if err != ErrFileDebitAmount {
			t.Errorf("Unexpected File.Validation error: %v", err.Error())
		}
	}
}

func TestFileCreditAmount(t *testing.T) {

	r := NewReader(strings.NewReader(" "))
	mockBatch := Batch{}
	mockBatch.Control.EntryAddendaCount = 1
	mockBatch.Control.TotalCreditEntryDollarAmount = 10500

	r.file.addBatch(mockBatch)
	r.file.Control.BatchCount = 1
	r.file.Control.EntryAddendaCount = 1
	r.file.Control.TotalCreditEntryDollarAmountInFile = 10500

	if err := r.file.Validate(); err != nil {
		t.Errorf("Unexpected File.Validation error: %v", err.Error())
	}

	r.file.Control.TotalCreditEntryDollarAmountInFile = 0
	if err := r.file.Validate(); err != nil {
		if err != ErrFileCreditAmount {
			t.Errorf("Unexpected File.Validation error: %v", err.Error())
		}
	}
}