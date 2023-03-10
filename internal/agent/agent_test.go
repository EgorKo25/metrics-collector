package agent

import (
	"testing"
)

func BenchmarkMonitor_Run(b *testing.B) {

	var monitor *Monitor

	for i := 0; i < b.N; i++ {
		b.Run("Run()", func(b *testing.B) {
			monitor.Run()
		})
	}
}
