package cron

import "testing"

func TestRetryDelay(t *testing.T) {
	tests := []struct {
		retry uint32
		want  int64
	}{
		{1, 10},
		{2, 30},
		{3, 60},
		{4, 300},
		{5, 600},
		{20, 600},
	}
	for _, tt := range tests {
		if got := retryDelay(tt.retry); got != tt.want {
			t.Fatalf("retryDelay(%d) = %d, want %d", tt.retry, got, tt.want)
		}
	}
}
