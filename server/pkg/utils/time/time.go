package time

import (
	"fmt"
	"strings"
	"time"
)

const (
	minutesInOneDay  = time.Minute * 24 * 60
	minutesInOneYear = 365 * minutesInOneDay
)

func HumanFriendlyDuration(d time.Duration) string {
	if d <= minutesInOneDay {
		return d.String()
	}

    var b strings.Builder

    if d > minutesInOneYear {
        years := d / minutesInOneYear 
        fmt.Fprintf(&b, "%dy", years)
        d -= years * minutesInOneYear
    }
    days := d / minutesInOneDay
    d -= days * minutesInOneDay
    fmt.Fprintf(&b, "%dd%s", days, d)
	return b.String()

}
