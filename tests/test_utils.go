package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"zenapi/utils"
)

func TestEmailValid(t *testing.T) {

	Convey("Test Email Valid\n", t, func() {
		So(utils.IsValidEmail("asdf@example.com"), ShouldEqual, true)
		So(utils.IsValidEmail("aABB123@eCCxample.com"), ShouldEqual, true)
		So(utils.IsValidEmail("@example.com"), ShouldEqual, false)
		So(utils.IsValidEmail("asdf@example"), ShouldEqual, false)
	})
}

func TestStringInSlice(t *testing.T) {
	Convey("Test StringInSlice", t, func() {
		So(utils.StringInSlice("asdf", []string{"asdf", "", "qwer"}), ShouldEqual, true)
	})
}
