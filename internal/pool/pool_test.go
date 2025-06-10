package pool

import "testing"

func TestPool_AddRemove(t *testing.T) {
	p := New(2)
	for i := 0; i < 5; i++ {
		p.AddWorker()
	}
	for i := 0; i < 10; i++ {
		p.Submit("x")
	}
	for i := 0; i < 3; i++ {
		p.RemoveWorker()
	}
	p.Shutdown()
}
