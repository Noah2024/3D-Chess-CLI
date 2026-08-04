// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// test cases for VecToUint and UintToVec
package bitutil_test

import (
	"3DC/util/bitutil"
	"fmt"
	"testing"
)

func TestBitutil(t *testing.T) {
	t.Run("VecToUintBoundry#1", func(t *testing.T) {
		got := bitutil.VecToUint(0, 0, 0)
		expected := uint32(0)
		if got != expected {
			t.Errorf("Expected %d but got %d ", expected, got)
		}

	})
	t.Run("VecToUintBoundry#2", func(t *testing.T) {
		got := bitutil.VecToUint(7, 7, 7)
		expected := uint32(511)
		if got != expected {
			t.Errorf("Expected %d but got %d ", expected, got)
		}

	})

	t.Run("VecToUint", func(t *testing.T) {
		got := bitutil.VecToUint(2, 3, 4)
		expected := uint32(226)
		if got != expected {
			t.Errorf("Expected %d but got %d ", expected, got)
		}
	})

	t.Run("VecToUintOutsideRange#1", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		fmt.Println("!!! A Formatted Error Below Here is Expected!!!")
		bitutil.VecToUint(7, 7, 8)
	})

	t.Run("VecToUintOutsideRange#2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		fmt.Println("!!! A Formatted Error Below Here is Expected!!!")
		bitutil.VecToUint(-1, 0, 0)
	})

	t.Run("UintToVecBoundry#1", func(t *testing.T) {
		gotx, goty, gotz := bitutil.UintToVec(uint32(0))
		expX, expY, expZ := 0, 0, 0
		if gotx != expX || goty != expY || gotz != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotx, goty, gotz)
		}
	})
	t.Run("UintToVecBoundry#2", func(t *testing.T) {
		gotx, goty, gotz := bitutil.UintToVec(uint32(511))
		expX, expY, expZ := 7, 7, 7
		if gotx != expX || goty != expY || gotz != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotx, goty, gotz)
		}
	})
	t.Run("UintToVec", func(t *testing.T) {
		gotx, goty, gotz := bitutil.UintToVec(uint32(226))
		expX, expY, expZ := 2, 3, 4
		if gotx != expX || goty != expY || gotz != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotx, goty, gotz)
		}
	})

	t.Run("UintToVecOutsideRange#1", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		fmt.Println("!!! A Formatted Error Below Here is Expected!!!")
		bitutil.UintToVec(512)
	})

	t.Run("UintToVecOutsideRange#2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		fmt.Println("!!! A Formatted Error Below Here is Expected!!!")
		bitutil.UintToVec(10000)

	})

	t.Run("BitutilInverseTest", func(t *testing.T) {
		expX, expY, expZ := 7, 2, 5
		asULoc := bitutil.VecToUint(expX, expY, expZ)
		gotX, gotY, gotZ := bitutil.UintToVec(asULoc)
		if gotX != expX || gotY != expY || gotZ != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotX, gotY, gotZ)
		}
	})
}
