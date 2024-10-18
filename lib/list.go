package lib

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

type ListItem struct {
	CreatedAt   time.Time
	Description string
	Id          int
	CompletedAt time.Time
}

type List struct {
	Items []ListItem
}

var template string = "%d\t%s\t%s\t%v"

func (list List) Print(all bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Id\tDescription\tCreatedAt\tCompletedAt")
	for _, item := range list.Items {
		time_diff := timediff.TimeDiff(item.CreatedAt)
		completed_time_diff := "Not completed"
		if item.IsComplete() {
			completed_time_diff = timediff.TimeDiff(item.CompletedAt)
		}
		if !item.IsComplete() || all {
			line := fmt.Sprintf(template, item.Id, item.Description, time_diff, completed_time_diff)
			fmt.Fprintln(w, line)
		}
	}
	w.Flush()
}

func (list List) Save() {
	WriteList(list)
}

func (item ListItem) IsComplete() bool {
	return !item.CompletedAt.IsZero()
}
