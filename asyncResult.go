package main

type asyncTask struct {
	input  []byte
	output []byte
}

func (at *asyncTask) getInput() []byte {
	return at.input
}

func (at *asyncTask) setInput(input []byte) {
	at.input = input
}

type asyncResult struct {
	done     chan []byte
	call     func([]byte, []byte) []byte
	numTasks int
	tasks    []asyncTask
}

func newAsyncResult(call func([]byte, []byte) []byte) *asyncResult {
	ar := new(asyncResult)
	ar.call = call
	ar.done = make(chan []byte)
	ar.tasks = make([]asyncTask, 0, 256)
	return ar
}

func (ar *asyncResult) resetTasks() {
	ar.numTasks = 0
}

func (ar *asyncResult) newTask() *asyncTask {
	if ar.numTasks >= len(ar.tasks) {
		ar.tasks = append(ar.tasks, asyncTask{})
	}
	at := &ar.tasks[ar.numTasks]
	ar.numTasks++
	return at
}

func (ar *asyncResult) runTask(i int) {
	at := ar.tasks[i]
	at.output = ar.call(at.output[:0], at.input)
	ar.done <- at.output
}

func (ar *asyncResult) await(output []byte) []byte {
	for i := 0; i < ar.numTasks; i++ {
		go ar.runTask(i)
	}
	for i := 0; i < ar.numTasks; i++ {
		output = append(output, <-ar.done...)
	}
	return output
}
