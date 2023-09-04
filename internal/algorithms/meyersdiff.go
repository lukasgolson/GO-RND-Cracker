package algorithms

func MeyersDifferenceAlgorithm(s1 []byte, s2 []byte) uint16 {

	if len(s1) == 0 {
		return uint16(len(s2)) // Return the length of s2 as the score
	}
	if len(s2) == 0 {
		return uint16(len(s1)) // Return the length of s1 as the score
	}

	score := uint16(len(s2))

	peq := make([]int64, 256)
	var i int

	for i = 0; i < len(peq); i++ {
		peq[i] = 0
	}

	for i = 0; i < len(s2); i++ {
		peq[s2[i]] |= int64(1) << uint(i)
	}

	var mv int64 = 0
	var pv int64 = -1
	var last = int64(1) << uint(len(s2)-1)

	for i = 0; i < len(s1); i++ {
		eq := peq[s1[i]]

		xv := eq | mv
		xh := (((eq & pv) + pv) ^ pv) | eq

		ph := mv | ^(xh | pv)
		mh := pv & xh

		if (ph & last) != 0 {
			score++
		}
		if (mh & last) != 0 {
			score--
		}

		ph = (ph << 1) | 1
		mh = mh << 1

		pv = mh | ^(xv | ph)
		mv = ph & xv
	}

	return score
}
