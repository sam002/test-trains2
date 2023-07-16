package solution

import "trains2/internal/trains2"

type Solution interface {
	Calculate() error
	GetSteps() *[]trains2.Move
}
