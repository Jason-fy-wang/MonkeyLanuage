package evaluator

import "testing"

func TestIfNestExpression(t *testing.T) {

	inputs := `
		if (10 > 1) {
			if(10>1) {
				return 10;
			}
			return 1;
		}
	`

	obj := testEval(inputs)

	t.Log(obj)
}
