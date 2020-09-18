# Queue
In-memory job queue. The queue exposes a REST API that producers and consumers perform HTTP requests against in JSON. 

The queue supports the following operations:
1) Enqueue
2) Dequeue
3) Conclude
4) GetJob
5) GetAllJobs
5) Remove

# Thought Process:

The initial thought, Can I use buffered channels as Queue. Although it may satisfy the functionality, didn't
use the approach since channels are the way to communicate between GoRoutines but not convinced to use as a data structure.
And yes, size limit was a main constraint for me to proceed further. For streaming homegrown job queue settling on the
size limit would require consideration of lot other factors.

So settled on the idea to implement Queue data structure, the next though whether it should slice-based or linked-list
based.

Slice based approach might waste a lot of memory, since it doesn't reuse memory occupied by removed items.

Linked-list approach is better for memory reuse. Of course with a overhead of maintaining links. 
It is easy to add/remove items from the middle of the list . It might be useful if the jobs need to be removed based 
on Type(TIME_CRITICAL/TIME_NOT_CRITICAL)

Implemented both the approaches. But In the handler created instance of linked-list queue.

As a next step want to do benchmarking testing using both the implementations for deeper insights.

# Assumptions:
Dequeue just does the peek of the queue with given constraints. But, not pop the item from the queue. 
Ignored Consumer ID sent in the header for deque.

# Improvements:
1) Add benchmark testing and load testing
2) Separate into different packages. Instead of all the files in cmd/server .

# Compile and Run:
make build 
./bin/server