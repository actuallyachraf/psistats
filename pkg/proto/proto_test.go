package proto

import "testing"

func TestProtocolPrimitives(t *testing.T) {
	t.Run("TestShuffle", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		aShuffled := Shuffle(a)
		t.Log(aShuffled)
	})
}
