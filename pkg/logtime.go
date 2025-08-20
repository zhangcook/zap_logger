package pkg

import (
	"fmt"
	"time"
)

func LogTime() string {
	format := time.Now().Format(time.DateTime)
	sprintf := fmt.Sprintf("%v", format)
	return sprintf[:10]
}
