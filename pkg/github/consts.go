package github

func ActionLabels() *LabelList {
	return &LabelList{
		{
			Name:  LabelValidationPending,
			Color: "ffff00",
		},
		{
			Name:  LabelValidationPassed,
			Color: "00ff00",
		},
		{
			Name:  LabelValidationFailed,
			Color: "ff0000",
		},
		{
			Name:  LabelValidationSkipped,
			Color: "000000",
		},
		{
			Name:  LabelNewProject,
			Color: "add8e6",
		},
	}
}

const (
	LabelValidationPending LabelName = "validation-pending"
	LabelValidationPassed  LabelName = "validation-passed"
	LabelValidationFailed  LabelName = "validation-failed"
	LabelValidationSkipped LabelName = "validation-skipped"
	LabelNewProject        LabelName = "new-project"
)
