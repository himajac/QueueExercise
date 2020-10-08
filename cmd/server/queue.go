package main

import (
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"math/rand"
	"sync"
	"time"
)

//Queue is an interface used for storing job details.
type Queue interface {
	Enqueue(item *job) (int, error)
	Dequeue(consumerId string) (*job, error)
	Conclude(jobID int, consumerId string) error
	GetJob(jobID int) (*job, error)
	GetJobs() (*[]job, error)
	Remove() (int, error)
}

//JobQueue is a concrete implementation of the Queue Interface using slice.
//sync.Map is used to store jobId and index as key,value pair. This helps to get the job in constant time.
type JobQueue struct {
	items []job
	count int
	m     *sync.Map
	mutex sync.Mutex

}

func NewQueue() *JobQueue {
	items := make([]job, 0)
	count := 0
	return &JobQueue{items, count, &sync.Map{}, sync.Mutex{}}
}

//Adds a job to the queue.And changes the job Status to "QUEUED". Returns JobId and error.
func (q *JobQueue) Enqueue(item *job) (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	//Generate random JobId and add the status.
	item.Id = rand.Int()
	item.Status = "QUEUED"

	q.items = append(q.items, *item)

	q.m.Store(item.Id, q.count) // Used to store the itemId with the index
	q.count++
	return item.Id, nil
}

//Returns a job from the queue . Jobs are considered available for Dequeue if the job has not been concluded or has not been Dequeued already.
func (q *JobQueue) Dequeue(consumerId string) (*job, error) {
	if q.count == 0 {
		return nil, errors.New("Dequeue on empty Job Queue.No jobs to process")
	}

	//Dequeue from the front of the list. But given the constraint need to verify the status.
	for index := 0; index < q.count; index++ {
		if q.items[index].Status == "QUEUED" {
			//Lock on the item to be dequeued instead of entire queue.This allows operations like enqueue & conclude without blocking
			q.items[index].mutex.Lock()
			q.items[index].Status = "IN_PROGRESS"
			q.items[index].mutex.Unlock()
			return &q.items[index], nil
		}
	}

	return nil, errors.New("Jobs are not available to process.")
}

//Given JobID ,finish execution on the job and change the status to CONCLUDED.
func (q *JobQueue) Conclude(jobID int, consumerId string) error {
	if q.count == 0 {
		return errors.New("Empty Job Queue.No jobs to Conclude")
	}

	//Get the index from the map
	v, ok := q.m.Load(jobID)
	if !ok {
		return errors.New("JobId not present in the Queue")
	}

	index := v.(int)
	time.Sleep(2 * time.Millisecond)    //Instead of job execution code
	q.items[index].Status = "CONCLUDED" //Change the status to concluded.
	return nil
}

//Get job details given JobID
func (q *JobQueue) GetJob(jobID int) (*job, error) {
	if q.count == 0 {
		return nil, errors.New("Empty Job Queue.")
	}

	//Get the index from the map
	v, ok := q.m.Load(jobID)
	if !ok {
		return nil, errors.New("JobId not present in the Queue")
	}
	return &q.items[v.(int)], nil
}

//Get all the jobs in the Queue
func (q *JobQueue) GetJobs() (*[]job, error) {
	if q.count == 0 {
		return nil, errors.New("Empty Job Queue")
	}
	arr := make([]job, 0)
	for index := 0; index < q.count; index++ {
		arr = append(arr, q.items[index])
	}
	return &arr, nil
}

//Remove the item from the front of the Queue
func (q *JobQueue) Remove() (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.count == 0 {
		return 0, errors.New("Empty Job Queue.")
	}

	firstItem := q.items[0]
	q.items = q.items[1:]
	q.m.Delete(firstItem.Id)

	return firstItem.Id, nil
}

// LinkedList Node structure.
type Element struct {
	Value job
	Prev  *Element
	Next  *Element
}

//JobListQueue is a concrete implementation of the Queue Interface using LinkedList.
//sync.Map is used to store jobId and index as key,value pair. This helps to get the job in constant time.
type JobListQueue struct {
	head            *Element
	tail            *Element
	count           int
	log             log.Logger
	mutex           sync.Mutex
	m               sync.Map
	consumerDetails sync.Map
//
}

func NewLinkedListQueue(logger log.Logger) *JobListQueue {
	count := 0
	return &JobListQueue{nil, nil, count, logger, sync.Mutex{}, sync.Map{}, sync.Map{}}
}

//Adds a job to the queue.And changes the job Status to "QUEUED". Returns JobId and error.
func (q *JobListQueue) Enqueue(item *job) (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	//Generate random JobId and add the status.
	item.Id = rand.Int()
	item.Status = "QUEUED"
	newElement := Element{Value: *item}

	if q.head == nil {
		q.head = &newElement
		q.tail = &newElement
	} else {
		newElement.Prev = q.tail
		q.tail.Next = &newElement
		q.tail = &newElement
	}
	q.m.Store(item.Id, &newElement) //Used to store the itemId and Address of the item as key,value pair.
	q.count++
	return item.Id, nil
}

//Returns a job from the queue . Jobs are considered available for Dequeue if the job has not been concluded or has not been Dequeued already.
func (q *JobListQueue) Dequeue(consumerId string) (*job, error) {
	if q.head == nil {
		return nil, errors.New("Dequeue on empty Job Queue.No jobs to process.")
	}

	//Dequeue from the front of the list. But given the constraint need to verify the status.
	curr := q.head
	for curr != nil {
		if curr.Value.Status == "QUEUED" {
			//Check the type of the JOb. If time_critical return that. Otherwise move towards the end till you find the time_Critical.
			if curr.Value.Type == "TIME_CRITICAL" {
				//Found the job.
				//Lock on the item to be dequeued instead of entire queue.This allows operations like enqueue & conclude without blocking
				curr.Value.mutex.Lock()
				curr.Value.Status = "IN_PROGRESS"
				//Add the details to map
				q.consumerDetails.Store(curr.Value.Id, consumerId) //Used to store the itemId and Address of the item as key,value pair.
				curr.Value.mutex.Unlock()

				return &curr.Value, nil
			} else {
				for curr != nil {
					curr = curr.Next
					if curr.Value.Type == "TIME_CRITICAL" {
						//Found the job. Then do the return .
						//Lock on the item to be dequeued instead of entire queue.This allows operations like enqueue & conclude without blocking
						curr.Value.mutex.Lock()
						curr.Value.Status = "IN_PROGRESS"
						//Add the details to map
						q.consumerDetails.Store(curr.Value.Id, consumerId) //Used to store the itemId and Address of the item as key,value pair.
						curr.Value.mutex.Unlock()

						return &curr.Value, nil
					}
				}
				//It means there is no time_critical job in the queue. So i will return the head of the queue.
				//Lock on the item to be dequeued instead of entire queue.This allows operations like enqueue & conclude without blocking
				curr.Value.mutex.Lock()
				curr.Value.Status = "IN_PROGRESS"
				//Add the details to map
				q.consumerDetails.Store(curr.Value.Id, consumerId) //Used to store the itemId and Address of the item as key,value pair.
				curr.Value.mutex.Unlock()

				return &curr.Value, nil
			}
		}
		curr = curr.Next
	}
	return nil, errors.New("None of the jobs are available to deque")
}

//For the jobId provided ,finishes execution on the job and change the status to CONCLUDED
func (q *JobListQueue) Conclude(jobID int, consumerId string) error {
	if q.head == nil {
		return errors.New("Empty Job Queue.No jobs to Conclude.")
	}

	//Get the consumerId used to Dequeue the Job
	cId, ok := q.consumerDetails.Load(jobID)
	if !ok {
		return errors.New("JobId not present in the Queue")
	}

	if cId != consumerId {
		return errors.New("Not Valid consumer to Conclude the Job")
	}

	//Get the index from the map
	v, ok := q.m.Load(jobID)
	if !ok {
		return errors.New("JobId not present in the Queue")
	}

	addrOfElement := v.(*Element)

	addrOfElement.Value.mutex.Lock()
	time.Sleep(2 * time.Millisecond)         //Instead of job execution code
	addrOfElement.Value.Status = "CONCLUDED" //Change the status to concluded.
	addrOfElement.Value.mutex.Unlock()

	return nil
}

func (q *JobListQueue) Cancel(jobID int) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.head == nil {
		return errors.New("Empty Job Queue.No jobs to Conclude.")
	}

	//Get the index from the map
	v, ok := q.m.Load(jobID)
	if !ok {
		return errors.New("JobId not present in the Queue")
	}

	addrOfElement := v.(*Element)

	prev := addrOfElement.Prev
	next := addrOfElement.Next
	prev.Next = next
	next.Prev = prev

	addrOfElement = nil

	return nil
}

//Given a job ID, returns details about the job
func (q *JobListQueue) GetJob(jobID int) (*job, error) {
	if q.head == nil {
		return nil, errors.New("Empty Job Queue.")
	}

	//Get the index from the map
	v, ok := q.m.Load(jobID)
	if !ok {
		return nil, errors.New("JobId not present in the Queue")
	}
	addrOfElement := v.(*Element)
	return &addrOfElement.Value, nil
}

//Removes the job in front of the queue.
func (q *JobListQueue) Remove() (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.head == nil {
		return 0, errors.New("Empty Job Queue.")
	}
	//To store the next element in the list
	curr := q.head
	next := q.head.Next

	q.m.Delete(q.head.Value.Id) //Delete the key from the map

	//To remove connecting pointers
	curr.Next = nil
	q.head = next
	q.head.Prev = nil

	idDeleted := curr.Value.Id
	curr = nil
	return idDeleted, nil
}

//Returns info about all the jobs in the Queue.
func (q *JobListQueue) GetJobs() (*[]job, error) {
	if q.head == nil {
		return nil, errors.New("Empty Job Queue")
	}
	arr := make([]job, 0)
	curr := q.head
	for curr != nil {
		arr = append(arr, curr.Value)
		curr = curr.Next
	}
	return &arr, nil
}
