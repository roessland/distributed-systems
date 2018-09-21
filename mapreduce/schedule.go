package mapreduce

import (
	"fmt"
	"log"
	"sync"
)

// Gives work to a specific worker
func runWorker(worker string, work chan DoTaskArgs, wg *sync.WaitGroup) {
	for doTaskArgs := range work {
		log.Println("Attempting to make worker", worker, "do stuff")
		successful := call(worker, "Worker.DoTask", doTaskArgs, nil)
		if successful {
			log.Println("TaskNumber", doTaskArgs.TaskNumber, "is done")
			wg.Done()
		} else {
			log.Println("Worker", worker, "failed.")
			// Reschedule task and exit worker communication thread
			work <- doTaskArgs
			log.Println("Worker", worker, "returned work", doTaskArgs.TaskNumber, "to pool")
			return
		}
	}
}

//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan <-chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers. Once all tasks
	// have completed successfully, schedule() should return.
	//
	// Your code here (Part III, Part IV).
	//

	work := make(chan DoTaskArgs)

	// Produce work
	go func() {
		for i := 0; i < ntasks; i++ {
			doTaskArgs := DoTaskArgs{
				JobName:       jobName,
				Phase:         phase,
				TaskNumber:    i,
				NumOtherPhase: n_other,
			}
			if phase == mapPhase {
				doTaskArgs.File = mapFiles[i]
			}
			work <- doTaskArgs
		}
	}()

	// Delegate work
	wg := sync.WaitGroup{}
	wg.Add(ntasks)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case worker := <-registerChan:
				go runWorker(worker, work, &wg)
			}
		}
	}()
	wg.Wait()
	close(work)
	done <- true

	fmt.Printf("Schedule: %v done\n", phase)
}
