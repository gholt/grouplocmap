package grouplocmap

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"testing"

	"gopkg.in/gholt/brimutil.v1"
)

func TestNewRoots(t *testing.T) {
	glm := New(&Config{Roots: 16}).(*groupLocMap)
	if len(glm.roots) < 16 {
		t.Fatal(len(glm.roots))
	}
	glm = New(&Config{Roots: 17}).(*groupLocMap)
	if len(glm.roots) < 17 {
		t.Fatal(len(glm.roots))
	}
}

func TestSetNewKeyOldTimestampIs0AndNewKeySaved(t *testing.T) {
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp := uint64(2)
	blockID := uint32(1)
	offset := uint32(0)
	length := uint32(0)
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp, blockID, offset, length, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 + 2
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(4)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 - 2
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := blockID1 + 1
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, true)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(&Config{Roots: 1, PageSize: 1}).(*groupLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	oldTimestamp := glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA1, keyB1, 789)
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
	oldTimestamp = glm.Set(keyA2, keyB2, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = glm.Get(keyA2, keyB2, 789)
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
	glm := New(&Config{Roots: 1, PageSize: 1}).(*groupLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	oldTimestamp := glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA1, keyB1, 789)
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
	oldTimestamp = glm.Set(keyA2, keyB2, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = glm.Get(keyA2, keyB2, 789)
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
	oldTimestamp = glm.Set(keyA2, keyB2, 789, timestamp2, uint32(0), offset2, length2, true)
	if oldTimestamp != timestamp2 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = glm.Get(keyA2, keyB2, 789)
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
	oldTimestamp = glm.Set(keyA3, keyB3, 789, timestamp3, blockID3, offset3, length3, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet = glm.Get(keyA3, keyB3, 789)
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
	if glm.roots[0].used != 2 {
		t.Fatal(glm.roots[0].used)
	}
}

func TestSetOverflowingKeysLots(t *testing.T) {
	glm := New(&Config{Roots: 1, PageSize: 1, SplitMultiplier: 1000}).(*groupLocMap)
	keyA := uint64(0)
	timestamp := uint64(2)
	blockID := uint32(1)
	offset := uint32(2)
	length := uint32(3)
	for keyB := uint64(0); keyB < 100; keyB++ {
		glm.Set(keyA, keyB, 789, timestamp, blockID, offset, length, false)
		blockID++
		offset++
		length++
	}
	if glm.roots[0].used != 100 {
		t.Fatal(glm.roots[0].used)
	}
	if len(glm.roots[0].overflow) != 25 {
		t.Fatal(len(glm.roots[0].overflow))
	}
	blockID = uint32(1)
	offset = uint32(2)
	length = uint32(3)
	for keyB := uint64(0); keyB < 100; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
		timestampSet := glm.Set(keyA, keyB, 789, timestamp2, blockID, offset, length, false)
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
		timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	if glm.roots[0].used != 100 {
		t.Fatal(glm.roots[0].used)
	}
	timestamp3 := timestamp2 + 2
	for keyB := uint64(0); keyB < 50; keyB++ {
		timestampSet := glm.Set(keyA, keyB, 789, timestamp3, uint32(0), uint32(0), uint32(0), false)
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
		timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
		glm.Set(keyA, keyB, 789, timestamp4, blockID, offset, length, false)
		blockID++
		offset++
		length++
	}
	if glm.roots[0].used != 150 {
		t.Fatal(glm.roots[0].used)
	}
	blockID = uint32(1)
	offset = uint32(2)
	length = uint32(3)
	for keyB := uint64(0); keyB < 100; keyB++ {
		timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
		timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp := uint64(2)
	blockID := uint32(0)
	offset := uint32(4)
	length := uint32(5)
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp, blockID, offset, length, false)
	if oldTimestamp != 0 {
		t.Fatal(oldTimestamp)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 + 2
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(4)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1 - 2
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, false)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(2)
	blockID1 := uint32(1)
	offset1 := uint32(0)
	length1 := uint32(0)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	timestamp2 := timestamp1
	blockID2 := uint32(0)
	offset2 := offset1 + 1
	length2 := length1 + 1
	oldTimestamp := glm.Set(keyA, keyB, 789, timestamp2, blockID2, offset2, length2, true)
	if oldTimestamp != timestamp1 {
		t.Fatal(oldTimestamp, timestamp1)
	}
	timestampGet, blockIDGet, offsetGet, lengthGet := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	glm.Discard(0, math.MaxUint64, 2)
	timestamp2, blockID2, offset2, length2 := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA := uint64(0)
	keyB := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA, keyB, 789, timestamp1, blockID1, offset1, length1, false)
	glm.Discard(0, math.MaxUint64, 1)
	timestamp2, blockID2, offset2, length2 := glm.Get(keyA, keyB, 789)
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
	glm := New(nil).(*groupLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := glm.ScanCallback(0, math.MaxUint64, 0, 0, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		if groupKeyA2 == keyA1 && groupKeyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		}
		return true
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
	glm := New(nil).(*groupLocMap)
	keyA1 := uint64(100)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := glm.ScanCallback(101, math.MaxUint64, 0, 0, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		return false
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
	stopped, more = glm.ScanCallback(0, 99, 0, 0, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		return false
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
	glm := New(nil).(*groupLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := glm.ScanCallback(0, math.MaxUint64, 1, 0, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		if groupKeyA2 == keyA1 && groupKeyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		}
		return true
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
	stopped, more = glm.ScanCallback(0, math.MaxUint64, 2, 0, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		return false
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
	glm := New(nil).(*groupLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(1)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := glm.ScanCallback(0, math.MaxUint64, 0, 2, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		if groupKeyA2 == keyA1 && groupKeyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		}
		return true
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
	stopped, more = glm.ScanCallback(0, math.MaxUint64, 0, 1, math.MaxUint64, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		return false
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
	glm := New(nil).(*groupLocMap)
	keyA1 := uint64(0)
	keyB1 := uint64(0)
	timestamp1 := uint64(123)
	blockID1 := uint32(1)
	offset1 := uint32(2)
	length1 := uint32(3)
	glm.Set(keyA1, keyB1, 789, timestamp1, blockID1, offset1, length1, false)
	good := false
	stopped, more := glm.ScanCallback(0, math.MaxUint64, 0, 0, 123, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		if groupKeyA2 == keyA1 && groupKeyB2 == keyB1 {
			if timestamp2 != timestamp1 {
				t.Fatal(timestamp2)
			}
			if length2 != length1 {
				t.Fatal(length2)
			}
			good = true
		} else {
			t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		}
		return true
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
	stopped, more = glm.ScanCallback(0, math.MaxUint64, 0, 0, 122, 100, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		t.Fatalf("%x %x %d %d\n", groupKeyA2, groupKeyB2, timestamp2, length2)
		return false
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
	glm := New(&Config{Roots: 128}).(*groupLocMap)
	keyA := uint64(0)
	for i := 0; i < 4000; i++ {
		keyA += 0x0010000000000000
		glm.Set(keyA, 0, 789, 1, 2, 3, 4, false)
	}
	count := 0
	stopped, more := glm.ScanCallback(0, math.MaxUint64, 0, 0, math.MaxUint64, 50, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		count++
		return true
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
	stopped, more = glm.ScanCallback(0, math.MaxUint64, 0, 0, math.MaxUint64, 5000, func(groupKeyA2 uint64, groupKeyB2 uint64, nameHash2 uint64, timestamp2 uint64, length2 uint32) bool {
		count++
		return true
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

func TestStatsBasic(t *testing.T) {
	// count needs to be high enough to fill all the root pages, hit the
	// overflow of those pages, and some pages below that too.
	count := uint64(100000)
	// seed just provides a repeatable test scenario.
	seed := 1
	// Roots is set low to get deeper quicker.
	// PageSize is set low to cause more page creation and deletion.
	// SplitMultiplier is set low to get splits to happen quicker.
	glm := New(&Config{Workers: 1, Roots: 1, PageSize: 512, SplitMultiplier: 1}).(*groupLocMap)
	keyspace := make([]byte, count*16)
	brimutil.NewSeededScrambled(int64(seed)).Read(keyspace)
	// since scrambled doesn't guarantee uniqueness, we do that in the middle
	// of each key.
	for j := uint32(0); j < uint32(count); j++ {
		binary.BigEndian.PutUint32(keyspace[j*16+4:], j)
	}
	maskedCount := uint64(0)
	for i := len(keyspace) - 16; i >= 0; i -= 16 {
		ka := binary.BigEndian.Uint64(keyspace[i:])
		kb := binary.BigEndian.Uint64(keyspace[i+8:])
		ts := uint64(1)
		if kb&1 != 0 {
			ts = 2
			maskedCount++
		}
		glm.Set(ka, kb, 789, ts, 2, 3, 4, false)
	}
	glm.SetInactiveMask(0)
	stats := glm.Stats(false)
	if stats.ActiveCount != count {
		t.Fatal(stats.ActiveCount)
	}
	if stats.ActiveBytes != count*4 {
		t.Fatal(stats.ActiveBytes)
	}
	if strings.Contains(stats.String(), "roots") {
		t.Fatal("did not expect debug output")
	}
	glm.SetInactiveMask(1)
	stats = glm.Stats(false)
	if stats.ActiveCount != maskedCount {
		t.Fatal(fmt.Sprintf("%d %d", stats.ActiveCount, maskedCount))
	}
	if stats.ActiveBytes != maskedCount*4 {
		t.Fatal(fmt.Sprintf("%d %d", stats.ActiveBytes, maskedCount*4))
	}
	if strings.Contains(stats.String(), "roots") {
		t.Fatal("did not expect debug output")
	}
	glm.SetInactiveMask(0)
	stats = glm.Stats(true)
	if stats.ActiveCount != count {
		t.Fatal(stats.ActiveCount)
	}
	if stats.ActiveBytes != count*4 {
		t.Fatal(stats.ActiveBytes)
	}
	if !strings.Contains(stats.String(), "roots") {
		t.Fatal("should have been debug output")
	}
	glm.SetInactiveMask(1)
	stats = glm.Stats(true)
	if stats.ActiveCount != maskedCount {
		t.Fatal(stats.ActiveCount)
	}
	if stats.ActiveBytes != maskedCount*4 {
		t.Fatal(stats.ActiveBytes)
	}
	if !strings.Contains(stats.String(), "roots") {
		t.Fatal("should have been debug output")
	}
}

func TestExerciseSplitMergeDiscard(t *testing.T) {
	// count needs to be high enough to fill all the root pages, hit the
	// overflow of those pages, and some pages below that too.
	count := 100000
	// seed just provides a repeatable test scenario.
	seed := 1
	// Roots is set low to get deeper quicker.
	// PageSize is set low to cause more page creation and deletion.
	// SplitMultiplier is set low to get splits to happen quicker.
	glm := New(&Config{Workers: 1, Roots: 1, PageSize: 768, SplitMultiplier: 1}).(*groupLocMap)
	// Override the mergeLevel to make it happen more often.
	for i := 0; i < len(glm.roots); i++ {
		glm.roots[i].mergeLevel = glm.roots[i].splitLevel - 2
	}
	if glm.roots[0].mergeLevel < 10 {
		t.Fatal(glm.roots[0].mergeLevel)
	}
	keyspace := make([]byte, count*16)
	brimutil.NewSeededScrambled(int64(seed)).Read(keyspace)
	// since scrambled doesn't guarantee uniqueness, we do that in the middle
	// of each key.
	for j := uint32(0); j < uint32(count); j++ {
		binary.BigEndian.PutUint32(keyspace[j*16+4:], j)
	}
	kt := func(ka uint64, kb uint64, ts uint64, b uint32, o uint32, l uint32) {
		glm.Set(ka, kb, 789, ts, b, o, l, false)
		ts2, b2, o2, l2 := glm.Get(ka, kb, 789)
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
	glm.Discard(0, math.MaxUint64, 2)
	for i := len(keyspace) - 16; i >= 0; i -= 16 {
		kt(binary.BigEndian.Uint64(keyspace[i:]), binary.BigEndian.Uint64(keyspace[i+8:]), 2, 3, 4, 5)
	}
	glm.Discard(0, math.MaxUint64, 1)
	for i := len(keyspace) - 16; i >= 0; i -= 16 {
		kt(binary.BigEndian.Uint64(keyspace[i:]), binary.BigEndian.Uint64(keyspace[i+8:]), 3, 0, 0, 0)
	}
	glm.SetInactiveMask(0)
	stats := glm.Stats(false)
	if stats.ActiveCount != 0 {
		t.Fatal(stats.ActiveCount)
	}
	if stats.ActiveBytes != 0 {
		t.Fatal(stats.ActiveBytes)
	}
}

func TestGetGroup(t *testing.T) {
	glm := New(nil)
	for i := uint64(0); i < 1000; i++ {
		glm.Set(54320, 54321, i, 1, 2, 3, 4, false)
		glm.Set(54321, 54321, i, 1, 2, 3, 4, false)
	}
	locs := glm.GetGroup(54321, 54321)
	if len(locs) != 1000 {
		t.Fatal(len(locs))
	}
	for _, e := range locs {
		if e.GroupKeyA != 54321 {
			t.Fatal(e.GroupKeyA)
		}
		if e.GroupKeyB != 54321 {
			t.Fatal(e.GroupKeyB)
		}
		if e.Timestamp != 1 {
			t.Fatal(e.Timestamp)
		}
		if e.BlockID != 2 {
			t.Fatal(e.BlockID)
		}
		if e.Offset != 3 {
			t.Fatal(e.Offset)
		}
		if e.Length != 4 {
			t.Fatal(e.Length)
		}
	}
}
