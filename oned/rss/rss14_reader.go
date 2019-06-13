package rss

// Decodes RSS-14, including truncated and stacked variants. See ISO/IEC 24724:2006.

var (
	rss14_OUTSIDE_EVEN_TOTAL_SUBSET = []int{1, 10, 34, 70, 126}
	rss14_INSIDE_ODD_TOTAL_SUBSET   = []int{4, 20, 48, 81}
	rss14_OUTSIDE_GSUM              = []int{0, 161, 961, 2015, 2715}
	rss14_INSIDE_GSUM               = []int{0, 336, 1036, 1516}
	rss14_OUTSIDE_ODD_WIDEST        = []int{8, 6, 4, 3, 1}
	rss14_INSIDE_ODD_WIDEST         = []int{2, 4, 6, 8}

	rss14_FINDER_PATTERNS = [][]int{
		{3, 8, 2, 1},
		{3, 5, 5, 1},
		{3, 3, 7, 1},
		{3, 1, 9, 1},
		{2, 7, 4, 1},
		{2, 5, 6, 1},
		{2, 3, 8, 1},
		{1, 5, 7, 1},
		{1, 3, 9, 1},
	}
)
