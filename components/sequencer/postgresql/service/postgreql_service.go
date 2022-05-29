package service

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

/**
* @Author: azh
* @Date: 2022/5/13 22:18
* @Context:
 */

import (
	"context"
	"errors"
	"github.com/zouyx/agollo/v4/component/log"

	//_ "github.com/bmizerany/pq"
	"sync"
	"sync/atomic"
	"time"

	"mosn.io/layotto/components/sequencer/postgresql/dao"
	"mosn.io/layotto/components/sequencer/postgresql/model"
)

const (
	RETRY_MAX    = 3
	DEFAULT_STEP = 1
	MAX_STEP     = 10e6
)

type PostgresqlService struct {
	dao        *dao.PostgresqlDao
	overAllSeq *model.PostgresqlSeq
	mutex      sync.Mutex
}

func NewPostgresqlService(dao *dao.PostgresqlDao, seq *model.PostgresqlSeq) *PostgresqlService {
	return &PostgresqlService{
		dao:        dao,
		overAllSeq: seq,
	}
}

func (p *PostgresqlService) GetId(ctx context.Context, bizTag string) (uint64, error) {
	p.mutex.Lock()
	var err error
	seqList := p.overAllSeq.GetId(bizTag)
	if seqList == nil {
		// if it doesn't exist in memory, just initialize it
		seqList, err = p.InitCache(ctx, bizTag)
		if err != nil {
			p.mutex.Unlock()
			return 0, nil
		}
	}
	p.mutex.Unlock()

	var id uint64
	id, err = p.NextId(seqList)
	if err != nil {
		return 0, err
	}
	p.overAllSeq.Update(bizTag, seqList)
	return id, nil
}

func (p *PostgresqlService) Create(ctx context.Context, model *model.PostgresqlModel) error {
	if model.Step > MAX_STEP {
		return errors.New("step must less than MAX_STEP")
	}
	if len(model.Description) == 0 || len(model.BizTag) == 0 {
		return errors.New("description invalid and biz_tag invalid")
	}
	if model.Step == 0 {
		model.Step = DEFAULT_STEP
	}
	if model.MaxID == 0 {
		model.MaxID = 1
	}
	return p.dao.Add(ctx, model)
}

func (p *PostgresqlService) InitCache(ctx context.Context, bizTag string) (*model.PostgresqlAlloc, error) {
	postgresql, err := p.dao.NextSegment(ctx, bizTag)
	if err != nil {
		log.Errorf("initCache error, err: %v\n", err)
		return nil, err
	}
	alloc := model.NewPostgresqlAlloc(postgresql)
	alloc.Buffer = append(alloc.Buffer, model.NewPostgresqlSegment(postgresql))
	_ = p.overAllSeq.Add(alloc)
	return alloc, nil
}

func (p *PostgresqlService) NextId(current *model.PostgresqlAlloc) (uint64, error) {
	current.Lock()
	defer current.Unlock()

	var id uint64
	currentBuf := current.Buffer[current.CurrentPos]
	// judge whether the current buffer can be used
	if current.HasSeq() {
		id = atomic.AddUint64(&current.Buffer[current.CurrentPos].Cursor, 1)
		current.UpdateTime = time.Now()
	}

	// When 50% of the current number segment has been distributed,
	//if the next number segment has not been updated and loaded,
	//another thread will be added to update the number segment
	if currentBuf.Max-id < uint64(float32(current.Step) * 0.5) && len(current.Buffer) <= 1 && !current.IsPreload {
		current.IsPreload = true
		cancel, _ := context.WithTimeout(context.Background(), 3*time.Second)
		go p.preLoadBuf(cancel, current.Key, current)
	}

	// If the current buffer is used, switch to the next buffer and remove the current buffer
	if id == currentBuf.Max {
		//Determine whether the buffer is ready
		if len(current.Buffer) > 1 && current.Buffer[current.CurrentPos+1].InitOk {
			current.Buffer = append(current.Buffer[:0], current.Buffer[1:]...)
		}
	}
	if current.HasID(id) {
		return id, nil
	}

	//When you get here, it means that the current buffer has no available ID,
	//and the compensation thread has not finished running
	waitChan := make(chan byte, 1)
	current.Waiting[current.Key] = append(current.Waiting[current.Key], waitChan)

	// Other clients cannot be blocked while waiting
	current.Unlock()
	timer := time.NewTimer(500 * time.Millisecond)
	select {
	case <-waitChan:
	case <-timer.C:
	}
	current.Lock()
	if len(current.Buffer) <= 1 {
		return 0, errors.New("get id errror")
	}
	// switch uffer
	current.Buffer = append(current.Buffer[:0], current.Buffer[1:]...)
	if current.HasSeq() {
		id = atomic.AddUint64(&current.Buffer[current.CurrentPos].Cursor, 1)
		current.UpdateTime = time.Now()
	}
	return id, nil
}

func (p *PostgresqlService) preLoadBuf(ctx context.Context, bizTag string, current *model.PostgresqlAlloc) error {
	for i := 0; i < RETRY_MAX; i++ {
		pModel, err := p.dao.NextSegment(ctx, bizTag)
		if err != nil {
			log.Errorf("preLoadBuffer error, bizTag: %s, err: %v", bizTag, err)
			continue
		}
		segment := model.NewPostgresqlSegment(pModel)
		current.Buffer = append(current.Buffer, segment)
		p.overAllSeq.Update(bizTag, current)
		current.WakeUp()
		break
	}
	current.IsPreload = false
	return nil
}

func (p *PostgresqlService) InitMaxId(ctx context.Context, bizTag string, maxId int64, step int64) error {
	err := p.dao.InitMaxId(ctx, bizTag, maxId, step)
	if err != nil {
		return err
	}
	return nil
}