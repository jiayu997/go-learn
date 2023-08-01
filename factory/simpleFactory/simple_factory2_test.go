package simpleFactory

import (
	"fmt"
	"testing"
)

type Girl interface {
	Weigh() string
}

type FatGirl struct{}
type ThinGirl struct{}

func (g FatGirl) Weigh() string {
	return "200kg"
}

func (g ThinGirl) Weigh() string {
	return "70kg"
}

func GirlFactory(name string) Girl {
	switch name {
	case "fat":
		return new(FatGirl)
	case "thin":
		return new(ThinGirl)
	default:
		return new(ThinGirl)
	}
}

func TestMain1(t *testing.T) {
	girl := GirlFactory("fat")
	fmt.Println(girl.Weigh())
}
