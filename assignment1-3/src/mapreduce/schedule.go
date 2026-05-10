package mapreduce

import "sync"

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}
	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)
	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//

	var wg sync.WaitGroup

	for i := 0; i < ntasks; i++ {
		wg.Add(1)
		go func(taskNum int) {
			defer wg.Done()
			for {
				// get an available worker
				worker := <-mr.registerChannel

				// build the task argumentss
				args := DoTaskArgs{
					JobName:       mr.jobName,
					File:          mr.files[taskNum],
					Phase:         phase,
					TaskNumber:    taskNum,
					NumOtherPhase: nios,
				}

				// call the worker via RPC
				ok := call(worker, "Worker.DoTask", &args, new(struct{}))
				if ok {
					// task succeeded so put worker back and exit retry loop
					go func() { mr.registerChannel <- worker }()
					break
				}
				// task failed so loop again to get another worker
			}
		}(i)
	}

	wg.Wait()
	debug("Schedule: %v phase done\n", phase)
}