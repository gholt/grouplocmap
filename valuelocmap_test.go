package valuelocmap

import (
	"encoding/binary"
	"math"
	"testing"

	"gopkg.in/gholt/brimutil.v1"
)

func TestNewRoots(t *testing.T) {
	vlm := New(OptRoots(16)).(*valueLocMap)
	if len(vlm.roots) < 16 {
		t.Fatal(len(vlm.roots))
	}
	vlm = New(OptRoots(17)).(*valueLocMap)
	if len(vlm.roots) < 17 {
		t.Fatal(len(vlm.roots))
	}
}

func TestSetNewKeyOldTimestampIs0AndNewKeySaved(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp := uint64(2)
	blockID := uint32(1)
	offset := uint32(0)
	length := uint32(0)
	oldTimestamp := vlm.Set(keyA, keyB, timestamp, blockID, offset, length, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp {
		t.Fatal(timestampGet, timestamp)
	}
	if blockIDGet != blockID {
		t.Fatal(blockIDGet, blockID)
	}
	if offsetGet != offset {
		t.Fatal(offsetGet, offset)
	}
	if lengthGet != length {
		t.Fatal(lengthGet, length)
	}
}

func TestSetOverwriteKeyOldTimestampIsOldAndOverwriteWins(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 + 2
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp2 {
		t.Fatal(timestampGet, timestamp2)
	}
	if blockIDGet != blockID2 {
		t.Fatal(blockIDGet, blockID2)
	}
	if offsetGet != offset2 {
		t.Fatal(offsetGet, offset2)
	}
	if lengthGet != length2 {
		t.Fatal(lengthGet, length2)
	}
}

func TestSetOldOverwriteKeyOldTimestampIsPreviousAndPreviousWins(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(4)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 - 2
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp1 {
		t.Fatal(timestampGet, timestamp1)
	}
	if blockIDGet != blockID1 {
		t.Fatal(blockIDGet, blockID1)
	}
	if offsetGet != offset1 {
		t.Fatal(offsetGet, offset1)
	}
	if lengthGet != length1 {
		t.Fatal(lengthGet, length1)
	}
}

func TestSetOverwriteKeyOldTimestampIsSameAndOverwriteIgnored(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp1 {
		t.Fatal(timestampGet, timestamp1)
	}
	if blockIDGet != blockID1 {
		t.Fatal(blockIDGet, blockID1)
	}
	if offsetGet != offset1 {
		t.Fatal(offsetGet, offset1)
	}
	if lengthGet != length1 {
		t.Fatal(lengthGet, length1)
	}
}

func TestSetOverwriteKeyOldTimestampIsSameAndOverwriteWins(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, true)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp2 {
		t.Fatal(timestampGet, timestamp2)
	}
	if blockIDGet != blockID2 {
		t.Fatal(blockIDGet, blockID2)
	}
	if offsetGet != offset2 {
		t.Fatal(offsetGet, offset2)
	}
	if lengthGet != length2 {
		t.Fatal(lengthGet, length2)
	}
}

func TestSetOverflowingKeys(t *testing.T) {
	vlm := New(OptRoots(1), OptPageSize(1)).(*valueLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	oldTimestamp := vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA1, keyB1)
	if timestampGet != timestamp1 {
		t.Fatal(timestampGet, timestamp1)
	}
	if blockIDGet != blockID1 {
		t.Fatal(blockIDGet, blockID1)
	}
	if offsetGet != offset1 {
		t.Fatal(offsetGet, offset1)
	}
	if lengthGet != length1 {
		t.Fatal(lengthGet, length1)
	}
	keyA2 := uint64(0)
	keyB2 := uint64(2)
	timestamp2 := timestamp1 + 2
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp = vlm.Set(keyA2, keyB2, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = vlm.Get(keyA2, keyB2)
	if timestampGet != timestamp2 {
		t.Fatal(timestampGet, timestamp2)
	}
	if blockIDGet != blockID2 {
		t.Fatal(blockIDGet, blockID2)
	}
	if offsetGet != offset2 {
		t.Fatal(offsetGet, offset2)
	}
	if lengthGet != length2 {
		t.Fatal(lengthGet, length2)
	}
}

func TestSetOverflowingKeysReuse(t *testing.T) {
	vlm := New(OptRoots(1), OptPageSize(1)).(*valueLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	oldTimestamp := vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA1, keyB1)
	if timestampGet != timestamp1 {
		t.Fatal(timestampGet, timestamp1)
	}
	if blockIDGet != blockID1 {
		t.Fatal(blockIDGet, blockID1)
	}
	if offsetGet != offset1 {
		t.Fatal(offsetGet, offset1)
	}
	if lengthGet != length1 {
		t.Fatal(lengthGet, length1)
	}
	keyA2 := uint64(0)
	keyB2 := uint64(2)
	timestamp2 := timestamp1 + 2
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp = vlm.Set(keyA2, keyB2, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = vlm.Get(keyA2, keyB2)
	if timestampGet != timestamp2 {
		t.Fatal(timestampGet, timestamp2)
	}
	if blockIDGet != blockID2 {
		t.Fatal(blockIDGet, blockID2)
	}
	if offsetGet != offset2 {
		t.Fatal(offsetGet, offset2)
	}
	if lengthGet != length2 {
		t.Fatal(lengthGet, length2)
	}
	oldTimestamp = vlm.Set(keyA2, keyB2, timestamp2, uint32(0), offset2, length2, true)
	if oldTimestamp != timestamp2 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = vlm.Get(keyA2, keyB2)
	if timestampGet != 0 {
		t.Fatal(timestampGet)
	}
	if blockIDGet != 0 {
		t.Fatal(blockIDGet)
	}
	if offsetGet != 0 {
		t.Fatal(offsetGet)
	}
	if lengthGet != 0 {
		t.Fatal(lengthGet)
	}
	keyA3 := uint64(0)
	keyB3 := uint64(2)
	timestamp3 := timestamp1 + 4
	blockID3 := blockID1 + 2
	offset3 := offset1 + 2
	length3 := length1 + 2
	oldTimestamp = vlm.Set(keyA3, keyB3, timestamp3, blockID3, offset3, length3, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = vlm.Get(keyA3, keyB3)
	if timestampGet != timestamp3 {
		t.Fatal(timestampGet, timestamp3)
	}
	if blockIDGet != blockID3 {
		t.Fatal(blockIDGet, blockID3)
	}
	if offsetGet != offset3 {
		t.Fatal(offsetGet, offset3)
	}
	if lengthGet != length3 {
		t.Fatal(lengthGet, length3)
	}
	if vlm.roots[0].used != 2 {
		t.Fatal(vlm.roots[0].used)
	}
}

func TestSetOverflowingKeysLots(t *testing.T) {
	vlm := New(OptRoots(1), OptPageSize(1), OptSplitMultiplier(1000)).(*valueLocMap)
	keyA := uint64(0)
	timestamp := uint64(2)
	blockID := uint32(1)
	offset := uint32(2)
	length := uint32(3)
	for keyB := uint64(0); keyB < 100; keyB++ {
		vlm.Set(keyA, keyB, timestamp, blockID, offset, length, false)
		blockID++
		offset++
		length++
	}
	if vlm.roots[0].used != 100 {
		t.Fatal(vlm.roots[0].used)
	}
	if len(vlm.roots[0].overflow) != 25 {
		t.Fatal(len(vlm.roots[0].overflow))
	}
	blockID = uint32(1)
	offset = uint32(2)
	length = uint32(3)
	for keyB := uint64(0); keyB < 100; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
		if timestampGet != timestamp {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, timestamp)
		}
		if blockIDGet != blockID {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, blockID)
		}
		if offsetGet != offset {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, offset)
		}
		if lengthGet != length {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, length)
		}
		blockID++
		offset++
		length++
	}
	timestamp2 := timestamp + 2
	blockID = uint32(2)
	offset = uint32(3)
	length = uint32(4)
	for keyB := uint64(0); keyB < 75; keyB++ {
		timestampSet := vlm.Set(keyA, keyB, timestamp2, blockID, offset, length, false)
		if timestampSet != timestamp {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampSet, timestamp)
		}
		blockID++
		offset++
		length++
	}
	blockID = uint32(2)
	offset = uint32(3)
	length = uint32(4)
	for keyB := uint64(0); keyB < 75; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
		if timestampGet != timestamp2 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, timestamp2)
		}
		if blockIDGet != blockID {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, blockID)
		}
		if offsetGet != offset {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, offset)
		}
		if lengthGet != length {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, length)
		}
		blockID++
		offset++
		length++
	}
	if vlm.roots[0].used != 100 {
		t.Fatal(vlm.roots[0].used)
	}
	timestamp3 := timestamp2 + 2
	for keyB := uint64(0); keyB < 50; keyB++ {
		timestampSet := vlm.Set(keyA, keyB, timestamp3, uint32(0), uint32(0), uint32(0), false)
		if timestampSet != timestamp2 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampSet, timestamp2)
		}
		blockID++
		offset++
		length++
	}
	blockID = uint32(2)
	offset = uint32(3)
	length = uint32(4)
	for keyB := uint64(0); keyB < 50; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
		if timestampGet != 0 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, 0)
		}
		if blockIDGet != 0 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, 0)
		}
		if offsetGet != 0 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, 0)
		}
		if lengthGet != 0 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, 0)
		}
		blockID++
		offset++
		length++
	}
	timestamp4 := timestamp3 + 2
	blockID = uint32(7)
	offset = uint32(8)
	length = uint32(9)
	for keyB := uint64(200); keyB < 300; keyB++ {
		vlm.Set(keyA, keyB, timestamp4, blockID, offset, length, false)
		blockID++
		offset++
		length++
	}
	if vlm.roots[0].used != 150 {
		t.Fatal(vlm.roots[0].used)
	}
	blockID = uint32(1)
	offset = uint32(2)
	length = uint32(3)
	for keyB := uint64(0); keyB < 100; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
		if keyB < 50 {
			if timestampGet != 0 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, 0)
			}
			if blockIDGet != 0 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, 0)
			}
			if offsetGet != 0 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, 0)
			}
			if lengthGet != 0 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, 0)
			}
		} else if keyB < 75 {
			if timestampGet != timestamp2 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, timestamp2)
			}
			if blockIDGet != blockID+1 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, blockID+1)
			}
			if offsetGet != offset+1 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, offset+1)
			}
			if lengthGet != length+1 {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, length+1)
			}
		} else {
			if timestampGet != timestamp {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, timestamp)
			}
			if blockIDGet != blockID {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, blockID)
			}
			if offsetGet != offset {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, offset)
			}
			if lengthGet != length {
				t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, length)
			}
		}
		blockID++
		offset++
		length++
	}
	blockID = uint32(7)
	offset = uint32(8)
	length = uint32(9)
	for keyB := uint64(200); keyB < 300; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
		if timestampGet != timestamp4 {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, timestampGet, timestamp4)
		}
		if blockIDGet != blockID {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, blockIDGet, blockID)
		}
		if offsetGet != offset {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, offsetGet, offset)
		}
		if lengthGet != length {
			t.Fatalf("%016x %016x %d %d", keyA, keyB, lengthGet, length)
		}
		blockID++
		offset++
		length++
	}
}

func TestSetNewKeyBlockID0OldTimestampIs0AndNoEffect(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp := uint64(2)
	blockID := uint32(0)
	offset := uint32(4)
	length := uint32(5)
	oldTimestamp := vlm.Set(keyA, keyB, timestamp, blockID, offset, length, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != 0 {
		t.Fatal(timestampGet, 0)
	}
	if blockIDGet != 0 {
		t.Fatal(blockIDGet, 0)
	}
	if offsetGet != 0 {
		t.Fatal(offsetGet, 0)
	}
	if lengthGet != 0 {
		t.Fatal(lengthGet, 0)
	}
}

func TestSetOverwriteKeyBlockID0OldTimestampIsOldAndOverwriteWins(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 + 2
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != 0 {
		t.Fatal(timestampGet, 0)
	}
	if blockIDGet != 0 {
		t.Fatal(blockIDGet, 0)
	}
	if offsetGet != 0 {
		t.Fatal(offsetGet, 0)
	}
	if lengthGet != 0 {
		t.Fatal(lengthGet, 0)
	}
}

func TestSetOldOverwriteKeyBlockID0OldTimestampIsPreviousAndPreviousWins(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(4)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 - 2
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp1 {
		t.Fatal(timestampGet, timestamp1)
	}
	if blockIDGet != blockID1 {
		t.Fatal(blockIDGet, blockID1)
	}
	if offsetGet != offset1 {
		t.Fatal(offsetGet, offset1)
	}
	if lengthGet != length1 {
		t.Fatal(lengthGet, length1)
	}
}

func TestSetOverwriteKeyBlockID0OldTimestampIsSameAndOverwriteIgnored(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != timestamp1 {
		t.Fatal(timestampGet, timestamp1)
	}
	if blockIDGet != blockID1 {
		t.Fatal(blockIDGet, blockID1)
	}
	if offsetGet != offset1 {
		t.Fatal(offsetGet, offset1)
	}
	if lengthGet != length1 {
		t.Fatal(lengthGet, length1)
	}
}

func TestSetOverwriteKeyBlockID0OldTimestampIsSameAndOverwriteWins(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := vlm.Set(keyA, keyB, timestamp2, blockID2, offset2, length2, true)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := vlm.Get(keyA, keyB)
	if timestampGet != 0 {
		t.Fatal(timestampGet, 0)
	}
	if blockIDGet != 0 {
		t.Fatal(blockIDGet, 0)
	}
	if offsetGet != 0 {
		t.Fatal(offsetGet, 0)
	}
	if lengthGet != 0 {
		t.Fatal(lengthGet, 0)
	}
}

func TestDiscardMaskNoMatch(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	vlm.Discard(0, math.MaxUint64, 2)
	timestamp2, blockID2, offset2, length2 := vlm.Get(keyA, keyB)
	if timestamp2 != timestamp1 {
		t.Fatal(timestamp2)
	}
	if blockID2 != blockID1 {
		t.Fatal(blockID2)
	}
	if offset2 != offset1 {
		t.Fatal(offset2)
	}
	if length2 != length1 {
		t.Fatal(length2)
	}
}

func TestDiscardMaskMatch(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA, keyB, timestamp1, blockID1, offset1, length1, false)
	vlm.Discard(0, math.MaxUint64, 1)
	timestamp2, blockID2, offset2, length2 := vlm.Get(keyA, keyB)
	if timestamp2 != 0 {
		t.Fatal(timestamp2)
	}
	if blockID2 != 0 {
		t.Fatal(blockID2)
	}
	if offset2 != 0 {
		t.Fatal(offset2)
	}
	if length2 != 0 {
		t.Fatal(length2)
	}
}

func TestScanCallbackBasic(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := vlm.ScanCallback(0, math.MaxUint64, 0, 0, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		if keyA2 == keyA1 && keyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
		}
	})
	if !good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
}

func TestScanCallbackRangeMiss(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA1 := uint64(100)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := vlm.ScanCallback(101, math.MaxUint64, 0, 0, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
	})
	if good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
	good = false
	stopped, more = vlm.ScanCallback(0, 99, 0, 0, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
	})
	if good {
		t.Fatal("failed")
	}
	if stopped != 99 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
}

func TestScanCallbackMask(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := vlm.ScanCallback(0, math.MaxUint64, 1, 0, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		if keyA2 == keyA1 && keyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
		}
	})
	if !good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
	good = false
	stopped, more = vlm.ScanCallback(0, math.MaxUint64, 2, 0, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
	})
	if good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
}

func TestScanCallbackNotMask(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := vlm.ScanCallback(0, math.MaxUint64, 0, 2, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		if keyA2 == keyA1 && keyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
		}
	})
	if !good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
	good = false
	stopped, more = vlm.ScanCallback(0, math.MaxUint64, 0, 1, math.MaxUint64, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
	})
	if good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
}

func TestScanCallbackCutoff(t *testing.T) {
	vlm := New().(*valueLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(123)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	vlm.Set(keyA1, keyB1, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := vlm.ScanCallback(0, math.MaxUint64, 0, 0, 123, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		if keyA2 == keyA1 && keyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
		}
	})
	if !good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
	good = false
	stopped, more = vlm.ScanCallback(0, math.MaxUint64, 0, 0, 122, 100, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		t.Fatalf("%x %x %d %d\n", keyA2, keyB2, timestamp2, length2)
	})
	if good {
		t.Fatal("failed")
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
}

func TestScanCallbackMax(t *testing.T) {
	vlm := New(OptPageSize(128)).(*valueLocMap)
	keyA := uint64(0)
	for i := 0; i < 4000; i++ {
		keyA += 0x0010000000000000
		vlm.Set(keyA, 0, 1, 2, 3, 4, false)
	}
	count := 0
	stopped, more := vlm.ScanCallback(0, math.MaxUint64, 0, 0, math.MaxUint64, 50, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		count++
	})
	if count != 50 {
		t.Fatal(count)
	}
	if stopped == math.MaxUint64 {
		t.Fatal(stopped)
	}
	if !more {
		t.Fatal("should have been more")
	}
	count = 0
	stopped, more = vlm.ScanCallback(0, math.MaxUint64, 0, 0, math.MaxUint64, 5000, func(keyA2 uint64, keyB2 uint64, timestamp2 uint64, length2 uint32) {
		count++
	})
	if count != 4000 {
		t.Fatal(count)
	}
	if stopped != math.MaxUint64 {
		t.Fatal(stopped)
	}
	if more {
		t.Fatal("should not have been more")
	}
}

func TestGatherStatsBasic(t *testing.T) {
	vlm := New().(*valueLocMap)
	vlm.Set(0, 0, 1, 2, 3, 4, false)
	count, length, _ := vlm.GatherStats(0, false)
	if count != 1 {
		t.Fatal(count)
	}
	if length != 4 {
		t.Fatal(length)
	}
	count, length, _ = vlm.GatherStats(1, false)
	if count != 0 {
		t.Fatal(count)
	}
	if length != 0 {
		t.Fatal(length)
	}
	count, length, _ = vlm.GatherStats(0, true)
	if count != 1 {
		t.Fatal(count)
	}
	if length != 4 {
		t.Fatal(length)
	}
	count, length, _ = vlm.GatherStats(1, true)
	if count != 0 {
		t.Fatal(count)
	}
	if length != 0 {
		t.Fatal(length)
	}
}

func TestExerciseSplitMergeDiscard(t *testing.T) {
	// count needs to be high enough to fill all the root pages, hit the
	// overflow of those pages, and some pages below that too.
	count := 100000
	// seed just provides a repeatable test scenario.
	seed := 1
	// OptRoots is set low to get deeper quicker.
	// OptPageSize is set low to cause more page creation and deletion.
	// OptSplitMultiplier is set low to get splits to happen quicker.
	vlm := New(OptWorkers(1), OptRoots(1), OptPageSize(512), OptSplitMultiplier(1)).(*valueLocMap)
	// Override the mergeLevel to make it happen more often.
	for i := 0; i < len(vlm.roots); i++ {
		vlm.roots[i].mergeLevel = vlm.roots[i].splitLevel - 2
	}
	if vlm.roots[0].mergeLevel < 10 {
		t.Fatal(vlm.roots[0].mergeLevel)
	}
	keyspace := make([]byte, count*16)
	brimutil.NewSeededScrambled(int64(seed)).Read(keyspace)
	// since scrambled doesn't guarantee uniqueness, we do that in the middle
	// of each key.
	for j := uint32(0); j < uint32(count); j++ {
		binary.BigEndian.PutUint32(keyspace[j*16+4:], j)
	}
	kt := func(ka uint64, kb uint64, ts uint64, b uint32, o uint32, l uint32) {
		vlm.Set(ka, kb, ts, b, o, l, false)
		ts2, b2, o2, l2 := vlm.Get(ka, kb)
		if (b != 0 && ts2 != ts) || (b == 0 && ts2 != 0) {
			t.Fatalf("%x %x %d %d %d %d ! %d", ka, kb, ts, b, o, l, ts2)
		}
		if b2 != b {
			t.Fatalf("%x %x %d %d %d %d ! %d", ka, kb, ts, b, o, l, b2)
		}
		if o2 != o {
			t.Fatalf("%x %x %d %d %d %d ! %d", ka, kb, ts, b, o, l, o2)
		}
		if l2 != l {
			t.Fatalf("%x %x %d %d %d %d ! %d", ka, kb, ts, b, o, l, l2)
		}
	}
	for i := len(keyspace) - 16; i >= 0; i -= 16 {
		kt(binary.BigEndian.Uint64(keyspace[i:]), binary.BigEndian.Uint64(keyspace[i+8:]), 1, 2, 3, 4)
	}
	vlm.Discard(0, math.MaxUint64, 2)
	for i := len(keyspace) - 16; i >= 0; i -= 16 {
		kt(binary.BigEndian.Uint64(keyspace[i:]), binary.BigEndian.Uint64(keyspace[i+8:]), 2, 3, 4, 5)
	}
	vlm.Discard(0, math.MaxUint64, 1)
	for i := len(keyspace) - 16; i >= 0; i -= 16 {
		kt(binary.BigEndian.Uint64(keyspace[i:]), binary.BigEndian.Uint64(keyspace[i+8:]), 3, 0, 0, 0)
	}
	endingCount, length, _ := vlm.GatherStats(uint64(0), false)
	if endingCount != 0 {
		t.Fatal(endingCount)
	}
	if length != 0 {
		t.Fatal(length)
	}
}
