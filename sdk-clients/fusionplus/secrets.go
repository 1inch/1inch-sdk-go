package fusionplus

import "sort"

// PendingSecretIndexes returns the secret indexes a maker must reveal for the
// given ready-to-accept fills and has not revealed yet.
//
// Each ready fill names the secret index the relayer expects in its Idx field.
// A maker must reveal exactly that secret, once per index. It must not reveal
// secrets in submission order: a fill can become ready for any index. For
// example, a single resolver that fills the whole order needs the last secret,
// so the first ready fill can name a non-zero index.
//
// The result is sorted and deduplicated. It drops any index that is out of
// range for secretCount or is already marked in revealed (revealed may be nil).
// The caller submits each returned index, then records it in revealed only
// after a successful submission.
func PendingSecretIndexes(fills []ReadyToAcceptSecretFill, revealed map[int]bool, secretCount int) []int {
	seen := make(map[int]bool)
	var pending []int
	for _, fill := range fills {
		idx := int(fill.Idx)
		if idx < 0 || idx >= secretCount {
			continue
		}
		if revealed[idx] || seen[idx] {
			continue
		}
		seen[idx] = true
		pending = append(pending, idx)
	}
	sort.Ints(pending)
	return pending
}
