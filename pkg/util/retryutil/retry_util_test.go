package retryutil

import (
	"errors"
	"testing"
)

type mockTicker struct {
	test    *testing.T
	counter int
}

func (t *mockTicker) Stop() {}

func (t *mockTicker) Tick() {
	t.counter += 1
}

func TestRetryWorkerSuccess(t *testing.T) {
	tick := &mockTicker{t, 0}
	result := RetryWorker(10, 20, func() (bool, error) {
		return true, nil
	}, tick)

	if result != nil {
		t.Errorf("Wrong result, expected: %#v, got: %#v", nil, result)
	}

	if tick.counter != 0 {
		t.Errorf("Ticker was started once, but it shouldn't be")
	}
}

func TestRetryWorkerOneFalse(t *testing.T) {
	var counter = 0

	tick := &mockTicker{t, 0}
	result := RetryWorker(1, 3, func() (bool, error) {
		counter += 1

		if counter <= 1 {
			return false, nil
		} else {
			return true, nil
		}
	}, tick)

	if result != nil {
		t.Errorf("Wrong result, expected: %#v, got: %#v", nil, result)
	}

	if tick.counter != 1 {
		t.Errorf("Ticker was started %#v, but supposed to be just once", tick.counter)
	}
}

func TestRetryWorkerError(t *testing.T) {
	fail := errors.New("Error")

	tick := &mockTicker{t, 0}
	result := RetryWorker(1, 3, func() (bool, error) {
		return false, fail
	}, tick)

	if result != fail {
		t.Errorf("Wrong result, expected: %#v, got: %#v", fail, result)
	}
}
