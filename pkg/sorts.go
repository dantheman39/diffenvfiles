package pkg

type sortableByFirst [][]string

func (ss sortableByFirst) Len() int {
	return len(ss)
}

func (ss sortableByFirst) Less(i, j int) bool {
	iFirst := ""
	jFirst := ""
	if len(ss[i]) > 0 {
		iFirst = ss[i][0]
	}
	if len(ss[j]) > 0 {
		jFirst = ss[j][0]
	}
	isLess := iFirst < jFirst
	return isLess
}

func (ss sortableByFirst) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
