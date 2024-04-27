package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHumanFriendlyDuration(t *testing.T) {
	t.Run("MoreThanYear", func(t *testing.T) {
		actual := HumanFriendlyDuration(minutesInOneYear + 10 * minutesInOneDay)
		require.Equal(t, "1y10d0s", actual)
	})

	t.Run("MoreThanYearADay", func(t *testing.T) {
		actual := HumanFriendlyDuration(minutesInOneYear + time.Minute * 10 + time.Second * 5)
		require.Equal(t, "1y0d10m5s", actual)
	})

	// t.Run("LessThanYear", func(t *testing.T) {
 //        start_time := time.Now()
 //        timer := time.NewTimer(time.Second * 5)
 //        <-timer.C
 //        latency := minutesInOneYear + time.Since(start_time)
	// 	actual := HumanFriendlyDuration(latency)
	// 	require.Equal(t, "10d0s", actual)
	// })
}
