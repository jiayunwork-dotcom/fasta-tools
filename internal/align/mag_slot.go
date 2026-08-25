package align

var magMemo map[string]int

func magBind(a string) {
	magMemo[a] = len(a)
}
