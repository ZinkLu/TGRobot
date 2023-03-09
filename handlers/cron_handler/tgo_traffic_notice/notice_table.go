package tgo_traffic_notice

type PercentReport struct {
	Current int
}

func (pr *PercentReport) Clear() {
	pr.Current = 0
}

func (pr *PercentReport) NeedReport(currentThereHold int) bool {
	if currentThereHold <= pr.Current {
		return false
	}
	pr.Current = currentThereHold
	return true
}
