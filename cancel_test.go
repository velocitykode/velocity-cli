package prism

import "testing"

type fakeModel struct{ cancelled bool }

func (m fakeModel) Cancelled() bool { return m.cancelled }

func TestExitOnCancel_noopWhenNotCancelled(t *testing.T) {
	called := false
	SetCancelHandler(func() { called = true })
	t.Cleanup(func() { SetCancelHandler(nil) })

	exited := false
	orig := exitFn
	exitFn = func(int) { exited = true }
	t.Cleanup(func() { exitFn = orig })

	ExitOnCancel(fakeModel{cancelled: false})

	if called {
		t.Fatal("handler fired when model was not cancelled")
	}
	if exited {
		t.Fatal("exitFn fired when model was not cancelled")
	}
}

func TestExitOnCancel_runsHandlerAndExits(t *testing.T) {
	called := false
	SetCancelHandler(func() { called = true })
	t.Cleanup(func() { SetCancelHandler(nil) })

	var code int
	exited := false
	orig := exitFn
	exitFn = func(c int) {
		exited = true
		code = c
	}
	t.Cleanup(func() { exitFn = orig })

	ExitOnCancel(fakeModel{cancelled: true})

	if !called {
		t.Fatal("handler did not fire on cancel")
	}
	if !exited {
		t.Fatal("exitFn did not fire on cancel")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestExitOnCancel_nonCancellerIgnored(t *testing.T) {
	exited := false
	orig := exitFn
	exitFn = func(int) { exited = true }
	t.Cleanup(func() { exitFn = orig })

	ExitOnCancel(struct{}{})

	if exited {
		t.Fatal("exitFn fired for non-Canceller value")
	}
}
