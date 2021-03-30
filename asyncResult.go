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
	join     chan []byte
	call     func([]byte, []byte) []byte
	numTasks int
	tasks    []asyncTask
}

func newAsyncResult(call func([]byte, []byte) []byte) *asyncResult {
	ar := new(asyncResult)
	ar.call = call
	ar.join = make(chan []byte)
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

func (ar *asyncResult) await(output []byte) []byte {
	for i := 0; i < ar.numTasks; i++ {
		go runTask(ar.join, ar.call, ar.tasks[i].output[:0], ar.tasks[i].input)
	}
	for i := 0; i < ar.numTasks; i++ {
		ar.tasks[i].output = <-ar.join
		output = append(output, ar.tasks[i].output...)
	}
	return output
}

func runTask(join chan<- []byte, call func([]byte, []byte) []byte, output, input []byte) {
	join <- call(output, input)
}