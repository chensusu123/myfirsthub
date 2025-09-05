package cluster

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/metrics"
)

func getProcessLabels(handleName string, err error) []*metrics.Dimension {
	code := "success"
	if err != nil {
		code = "failed"
	}
	return []*metrics.Dimension{
		{Name: "handleName", Value: handleName},
		//{Name: "processType", Value: processType},
		{Name: "Code", Value: code},
	}
}

// Set the time consumption interval.
// 1 10 20 40 80 160 1000 2000 3000
var (
	clientBounds = metrics.NewValueBounds(1.0, 10.0, 20.0, 40.0, 80.0, 160.0, 1000.0, 2000.0, 3000.0)
	serverBounds = clientBounds
)
