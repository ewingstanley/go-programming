package solution

func IsNotationValid(s string) bool {

	stack := []rune{}
	mapping := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
		}
		if char == ')' || char == '}' || char == ']' {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if mapping[top] != char {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0

}
