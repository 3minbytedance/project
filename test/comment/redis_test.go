package comment

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestAddCommentByCommentId(t *testing.T) {
	latestTime := "1691835930"
	latestTimeInt, _ := strconv.ParseInt(latestTime, 10, 64)

	latestTimeUnix := time.Unix(latestTimeInt, 0)
	fmt.Println(latestTimeUnix)
}
