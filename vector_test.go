package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestVectorLength(t *testing.T) {
	Convey("Vector{0,0} length should be 0.", t, func() {
		v := Vector{0, 0}
		So(v.Length(), ShouldBeZeroValue)
	})

	Convey("Vector{1,0} length should be 1.", t, func() {
		v := Vector{1, 0}
		So(v.Length(), ShouldEqual, 1)
	})

	Convey("Vector{0,-1} length should be 1.", t, func() {
		v := Vector{3, 4}
		So(v.Length(), ShouldEqual, 5)
	})

	Convey("Vector{3, 4} normalized length should be 1.", t, func() {
		v := Vector{3, 4}
		v.Normalize()
		So(v.Length(), ShouldEqual, 1)
	})

	Convey("Vector{1, 0} rotated 90 degrees should be Vector{0, 1}.", t, func() {
		v := Vector{1, 0}
		v.SetAngle(math.Pi / 2)
		So(v.Length(), ShouldEqual, 1)
	})

	Convey("Vector{1, 0} set to length 2 should set it and not change the angle.", t, func() {
		v := Vector{1, 0}
		oldAngle := v.Angle()
		v.SetLength(2)
		So(v.Length(), ShouldEqual, 2)
		So(oldAngle, ShouldEqual, v.Angle())
	})

	Convey("With two vectors, rotation and setting angle should work as expected.", t, func() {
		v := Vector{1, 0}
		v1 := Vector{1, 0}
		for i := 0.1; i < math.Pi*2; i = i + 0.1 {
			v1.Rotate(0.1)
			v.SetAngle(i)
			So(v.Angle(), ShouldAlmostEqual, v1.Angle())
			So(v.Length(), ShouldAlmostEqual, v1.Length())
			So(v[0], ShouldAlmostEqual, v1[0])
			So(v[1], ShouldAlmostEqual, v1[1])
		}
	})
}
