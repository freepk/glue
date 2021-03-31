package main

import "testing"

func test(a, b []byte) []byte {
	return append(a, b...)
}

func TestAsyncResult(t *testing.T) {
	ar := newAsyncResult(test)
	for i := 0; i < 64; i++ {
		ar.resetTasks()
		{
			at := ar.newTask()
			at.appendInput(append(at.resetInput(), []byte(`{1}`)...))
		}
		{
			at := ar.newTask()
			at.appendInput(append(at.resetInput(), []byte(`{2}`)...))
		}
		{
			at := ar.newTask()
			at.appendInput(append(at.resetInput(), []byte(`{3}`)...))
		}
		for j := 0; j < 1024; j++ {
			out := ar.await(nil)
			switch {
			case bytesEq(out, []byte(`{1}{2}{3}`)):
			case bytesEq(out, []byte(`{1}{3}{2}`)):
			case bytesEq(out, []byte(`{2}{1}{3}`)):
			case bytesEq(out, []byte(`{2}{3}{1}`)):
			case bytesEq(out, []byte(`{3}{1}{2}`)):
			case bytesEq(out, []byte(`{3}{2}{1}`)):
			default:
				t.Log(string(out))
				t.Fail()
			}
		}
	}
}

func BenchmarkAsyncResult(b *testing.B) {
	ar := newAsyncResult(test)
	ar.resetTasks()
	{
		at := ar.newTask()
		at.appendInput(append(at.resetInput(), []byte(`{1}`)...))
	}
	{
		at := ar.newTask()
		at.appendInput(append(at.resetInput(), []byte(`{2}`)...))
	}
	{
		at := ar.newTask()
		at.appendInput(append(at.resetInput(), []byte(`{3}`)...))
	}
	out := make([]byte, 0, 256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out = ar.await(out[:0])
	}
}
