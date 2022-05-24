//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mysql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"testing"
)

const (
	MySQLUrl = "localhost:xxx"
	Value    = 1
	Key      = "sequenceKey"
	Size     = 50
)

func TestMySQLSequencer_Init(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	comp := NewMySQLSequencer(log.DefaultLogger)
	comp.db = db

	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(0)
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(rows)
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))

	cfg.Properties["tableNameKey"] = "layotto_sequencer"
	cfg.Properties["connectionString"] = "connectionString"
	cfg.Properties["dataBaseName"] = "layotto"
	err = comp.Init(cfg)

	assert.Nil(t, err)
}

func TestMySQLSequencer_GetNextId(t *testing.T) {

	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO").WithArgs("layotto_sequencer", "sequenceKey", 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	req := &sequencer.GetNextIdRequest{Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	comp.db = db
	_, err = comp.GetNextId(req)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	assert.NoError(t, err)
}

func TestMySQLSequencer_GetSegment(t *testing.T) {

	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO").WithArgs("layotto_sequencer", Key, Size).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	req := &sequencer.GetSegmentRequest{Size: Size, Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	comp.db = db
	_, _, err = comp.GetSegment(req)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	assert.NoError(t, err)
}

func TestMySQLSequencer_Close(t *testing.T) {

	var MySQLUrl = MySQLUrl

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	comp := NewMySQLSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["MySQLHost"] = MySQLUrl

	comp.db = db
	_ = comp.Init(cfg)

	comp.Close(db)
}

func TestMySQLSequencer_Segment_Insert(t *testing.T) {
	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WithArgs("layotto_sequencer", Key, Size).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	segmentReq := &sequencer.GetSegmentRequest{Size: Size, Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}
	comp.db = db
	_, _, err = comp.GetSegment(segmentReq)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	assert.NoError(t, err)
}

func TestMySQLSequencer_GetNextId_Insert(t *testing.T) {
	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WithArgs("layotto_sequencer", "sequenceKey", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	req := &sequencer.GetNextIdRequest{Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	comp.db = db
	_, err = comp.GetNextId(req)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	assert.NoError(t, err)
}

func TestMySQLSequencer_GetNextId_InsertError(t *testing.T) {
	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WithArgs("layotto_sequencer", "sequenceKey", 1).WillReturnError(errors.New("insert error"))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	req := &sequencer.GetNextIdRequest{Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	comp.db = db
	_, err = comp.GetNextId(req)

	assert.Error(t, err)
}

func TestMySQLSequencer_Segment_InsertError(t *testing.T) {
	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WithArgs("layotto_sequencer", Key, Size).WillReturnError(errors.New("insert error"))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	segmentReq := &sequencer.GetSegmentRequest{Size: Size, Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}
	comp.db = db
	_, _, err = comp.GetSegment(segmentReq)

	assert.Error(t, err)
}

func TestMySQLSequencer_GetNextId_UpdateError(t *testing.T) {
	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnError(errors.New("update error"))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	req := &sequencer.GetNextIdRequest{Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	comp.db = db
	_, err = comp.GetNextId(req)

	assert.Error(t, err)
}

func TestMySQLSequencer_Segment_UpdateError(t *testing.T) {
	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnError(errors.New("update error"))
	mock.ExpectCommit()

	properties := make(map[string]string)

	properties["tableName"] = "layotto_sequencer"
	properties["connectionString"] = "connectionString"
	properties["dataBaseName"] = "layotto"

	req := &sequencer.GetSegmentRequest{Size: Size, Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	comp.db = db
	_, _, err = comp.GetSegment(req)

	assert.Error(t, err)
}
