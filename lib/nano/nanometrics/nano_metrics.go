package nanometrics

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/metrics"
)

// Unified all metrics report inside the framework. Every property starts with "freetk.".
var (
	// -----------------------------server----------------------------- //
	GlobalTaskGauge = metrics.Gauge("nano.global.task_count")
	UserCountGauge  = metrics.Gauge("nano.global.user_count")
)
