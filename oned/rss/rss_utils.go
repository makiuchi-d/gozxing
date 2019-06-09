package rss

/** Adapted from listings in ISO/IEC 24724 Appendix B and Appendix G. */

func RSSUtils_getRSSvalue(widths []int, maxWidth int, noNarrow bool) int {
	n := 0
	for _, width := range widths {
		n += width
	}
	val := 0
	narrowMask := uint(0)
	elements := len(widths)
	for bar := 0; bar < elements-1; bar++ {
		elmWidth := 1
		narrowMask |= 1 << uint(bar)
		for ; elmWidth < widths[bar]; elmWidth, narrowMask = elmWidth+1, narrowMask&^(1<<uint(bar)) {
			subVal := combins(n-elmWidth-1, elements-bar-2)
			if noNarrow && (narrowMask == 0) && (n-elmWidth-(elements-bar-1) >= elements-bar-1) {
				subVal -= combins(n-elmWidth-(elements-bar), elements-bar-2)
			}
			if elements-bar-1 > 1 {
				lessVal := 0
				for mxwElement := n - elmWidth - (elements - bar - 2); mxwElement > maxWidth; mxwElement-- {
					lessVal += combins(n-elmWidth-mxwElement-1,
						elements-bar-3)
				}
				subVal -= lessVal * (elements - 1 - bar)
			} else if n-elmWidth > maxWidth {
				subVal--
			}
			val += subVal
		}
		n -= elmWidth
	}
	return val
}

func combins(n, r int) int {
	maxDenom := n - r
	minDenom := r
	if n-r > r {
		minDenom = r
		maxDenom = n - r
	}
	val := 1
	j := 1
	for i := n; i > maxDenom; i-- {
		val *= i
		if j <= minDenom {
			val /= j
			j++
		}
	}
	for j <= minDenom {
		val /= j
		j++
	}
	return val
}
