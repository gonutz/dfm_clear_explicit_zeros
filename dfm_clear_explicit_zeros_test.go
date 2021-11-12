package main

import (
	"strings"
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/dfm"
)

func TestCleanseEmptyObjectDoesNothing(t *testing.T) {
	obj, err := dfm.ParseString(`object Dialog: TDialog
end`)
	check.Eq(t, err, nil)
	check.Eq(t, cleanseObject(obj), false)
	check.Eq(t, obj.String(), strings.Replace(`object Dialog: TDialog
end
`, "\n", "\r\n", -1))
}

func TestCleanseRemovesExplicitZerosFromObject(t *testing.T) {
	obj, err := dfm.ParseString(`object Dialog: TDialog
  ExplicitLeft = 0
  ExplicitTop = 0
  ExplicitWidth = 0
  ExplicitHeight = 0
end`)
	check.Eq(t, err, nil)
	check.Eq(t, cleanseObject(obj), true)
	check.Eq(t, obj.String(), strings.Replace(`object Dialog: TDialog
end
`, "\n", "\r\n", -1))
}

func TestCleanseLeavesZerosIfNotAllAreZero(t *testing.T) {
	obj, err := dfm.ParseString(`object Dialog: TDialog
  ExplicitLeft = 0
  ExplicitTop = 1
  ExplicitWidth = 0
  ExplicitHeight = 0
end`)
	check.Eq(t, err, nil)
	check.Eq(t, cleanseObject(obj), false)
	check.Eq(t, obj.String(), strings.Replace(`object Dialog: TDialog
  ExplicitLeft = 0
  ExplicitTop = 1
  ExplicitWidth = 0
  ExplicitHeight = 0
end
`, "\n", "\r\n", -1))
}

func TestCleanseLeavesZerosIfNotAllAreThere(t *testing.T) {
	obj, err := dfm.ParseString(`object Dialog: TDialog
  ExplicitLeft = 0
  ExplicitWidth = 0
  ExplicitHeight = 0
end`)
	check.Eq(t, err, nil)
	check.Eq(t, cleanseObject(obj), false)
	check.Eq(t, obj.String(), strings.Replace(`object Dialog: TDialog
  ExplicitLeft = 0
  ExplicitWidth = 0
  ExplicitHeight = 0
end
`, "\n", "\r\n", -1))
}

func TestCleanseHandlesSubObject(t *testing.T) {
	obj, err := dfm.ParseString(`object Dialog: TDialog
  ExplicitLeft = 0
  ExplicitTop = 0
  ExplicitWidth = 0
  ExplicitHeight = 0
  object Sub: Thing
    ExplicitLeft = 0
    ExplicitTop = 0
    ExplicitWidth = 0
    ExplicitHeight = 0
  end
end`)
	check.Eq(t, err, nil)
	check.Eq(t, cleanseObject(obj), true)
	check.Eq(t, obj.String(), strings.Replace(`object Dialog: TDialog
  object Sub: Thing
  end
end
`, "\n", "\r\n", -1))
}
