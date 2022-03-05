package leetcode

//字符串是否是有效括号
func isValid(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	if length%2 != 0 {
		return false
	}
	var (
		leftMap = map[rune]struct{}{
			'(': {},
			'{': {},
			'[': {},
		}
		rightMap = map[rune]rune{
			')': '(',
			'}': '{',
			']': '[',
		}
		temp = make([]rune, 0)
	)
	for _, r := range s {
		if _, exist := leftMap[r]; exist {
			temp = append(temp, r)
		}
		if rt, exist := rightMap[r]; exist {
			if len(temp) == 0 {
				return false
			}
			if rt != temp[len(temp)-1] {
				return false
			}
			temp = temp[:len(temp)-1]
		}
	}
	return len(temp) == 0
}
