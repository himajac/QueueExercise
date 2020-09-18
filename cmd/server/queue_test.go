package main

import (
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobQueue_Enqueue(t *testing.T) {
	var q = NewQueue()
	item1 := job{Type: "TIME_CRITICAL"}
	item2 := job{Type: "NOT_TIME_CRITICAL"}
	q.Enqueue(&item1)
	q.Enqueue(&item2)
	if q.count != 2 {
		t.Errorf("got %d expected %d \n", q.count, 2)
	}
	if q.items[0].Status != "QUEUED" {
		t.Errorf("got %s expected %s \n", q.items[0].Status, "QUEUED")
	}
	v, ok := q.m.Load(item1.Id)
	if !ok {
		t.Errorf("item %d should be present in the list", item1.Id)
	}
	if v.(int) != q.count-2 {
		t.Errorf("item %d should be mapped to index %d", item1.Id, q.count-2)
	}
}

func TestJobQueue_Dequeue(t *testing.T) {
	var q = NewQueue()
	item1 := job{Type: "TIME_CRITICAL"}
	item2 := job{Type: "NOT_TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)
	id2, _ := q.Enqueue(&item2)

	firstItem, _ := q.Dequeue()
	if firstItem.Id != id1 {
		t.Errorf("got %d expected %d \n", firstItem.Id, id1)
	}
	if firstItem.Status != "IN_PROGRESS" {
		t.Errorf("got %s expected %s \n", firstItem.Status, "IN_PROGRESS")
	}
	secondItem, _ := q.Dequeue()
	if secondItem.Id != id2 {
		t.Errorf("got %d expected %d \n", secondItem.Id, id2)
	}
	if secondItem.Status != "IN_PROGRESS" {
		t.Errorf("got %s expected %s \n", secondItem.Status, "IN_PROGRESS")
	}
}

func TestJobQueue_Conclude(t *testing.T) {
	var q = NewQueue()
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	q.Conclude(id1)
	if q.items[0].Status != "CONCLUDED" {
		t.Errorf("got %s expected %s \n", q.items[0].Status, "CONCLUDED")
	}
}

func TestJobQueue_GetJob(t *testing.T) {
	var q *JobQueue = NewQueue()
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	item, _ := q.GetJob(id1)
	if item.Id != id1 {
		t.Errorf("got %d expected %d \n", item.Id, id1)
	}

	var id2 int
	expectedError := "JobId not present in the Queue"
	_, actualError := q.GetJob(id2)
	assert.Equal(t, expectedError, actualError.Error())
}

func TestJobListQueue_Enqueue(t *testing.T) {
	var q = NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	if q.count != 1 {
		t.Errorf("got %d expected %d \n", q.count, 1)
	}
	if q.head.Value.Id != id1 {
		t.Errorf("got %d expected %d \n", q.head.Value.Id, id1)
	}
	v, ok := q.m.Load(item1.Id)
	if !ok {
		t.Errorf("item %d should be present in the list", item1.Id)
	}
	addrOfElement := v.(*Element)
	if addrOfElement.Value.Id != id1 {
		t.Errorf("item %d should be mapped to index %d", item1.Id, addrOfElement.Value.Id)
	}
}

func TestJobListQueue_Dequeue(t *testing.T) {
	var q *JobListQueue = NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)
	item2 := job{Type: "NOT_TIME_CRITICAL"}
	id2, _ := q.Enqueue(&item2)

	firstItem, _ := q.Dequeue()
	if firstItem.Id != id1 {
		t.Errorf("got %d expected %d \n", firstItem.Id, id1)
	}
	secondItem, _ := q.Dequeue()
	if secondItem.Id != id2 {
		t.Errorf("got %d expected %d \n", secondItem.Id, id2)
	}
}

func TestJobListQueue_Conclude(t *testing.T) {
	var q *JobListQueue = NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	q.Conclude(id1)
	if q.head.Value.Status != "CONCLUDED" {
		t.Errorf("got %s expected %s \n", q.head.Value.Status, "CONCLUDED")
	}
}

func TestJobListQueue_GetJob(t *testing.T) {
	var q *JobListQueue = NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	item, _ := q.GetJob(id1)
	if item.Id != id1 {
		t.Errorf("got %d expected %d \n", item.Id, id1)
	}

	var id2 int
	expectedError := "JobId not present in the Queue"
	_, actualError := q.GetJob(id2)
	assert.Equal(t, expectedError, actualError.Error())
}
