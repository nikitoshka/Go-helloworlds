package tempFlag

import (
	"flag"
	"fmt"

	"../tempConv"
)

type tempToCel struct {
	tempConv.Cel
}

func (t *tempToCel) Set(s string) error {
	var value float64
	var measure string

	if _, err := fmt.Sscanf(s, "%f%s", &value, &measure); err != nil {
		return err
	}

	switch measure {
	case "C":
		t.Cel = tempConv.Cel(value)

	case "F":
		t.Cel = tempConv.Far(value).ToCel()

	case "K":
		t.Cel = tempConv.Kel(value).ToCel()

	default:
		return fmt.Errorf("error: wrong measure(%s) for a given value (%f)", measure, value)
	}

	return nil
}

func TempFlag(name string, defaultValue float64, usage string) *tempConv.Cel {
	var t = tempToCel{tempConv.Cel(defaultValue)}
	flag.Var(&t, name, usage)

	return &t.Cel
}
