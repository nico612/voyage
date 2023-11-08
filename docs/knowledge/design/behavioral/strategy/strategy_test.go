package strategy

import (
	"fmt"
	"testing"
)

func TestOperator_calculate(t *testing.T) {

	o := &Operator{}
	o.SetStrategy(&add{})
	result := o.calculate(1, 2)
	fmt.Println("add: ", result)

	o.SetStrategy(&reduce{})
	result = o.calculate(2, 1)
	fmt.Println("reduce: ", result)

}
