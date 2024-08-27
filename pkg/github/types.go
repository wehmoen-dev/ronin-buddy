package github

type ReviewStatus string

const (
	ReviewStatusPending ReviewStatus = "pending"
	ReviewStatusPassed  ReviewStatus = "passed"
	ReviewStatusFailed  ReviewStatus = "failed"
)

type LabelName string

type Label struct {
	Name  LabelName
	Color string
}

type LabelList []Label
