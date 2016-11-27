package tempConv

import "fmt"

type Cel float64
type Far float64
type Kel float64

func (c Cel) String() string {
	return fmt.Sprintf("%.2f°C", c)
}

func (c Cel) ToFar() Far {
	return Far(c)*1.8 + 32
}

func (c Cel) ToKel() Kel {
	return Kel(c) - 273
}

func (f Far) String() string {
	return fmt.Sprintf("%.2f°F", f)
}

func (f Far) ToCel() Cel {
	return (Cel(f) - 32) / 1.8
}

func (f Far) ToKel() Kel {
	return Kel(f.ToCel()) - 273
}

func (k Kel) String() string {
	return fmt.Sprintf("%.2f°K", k)
}

func (k Kel) ToCel() Cel {
	return Cel(k) + 273
}

func (k Kel) ToFar() Far {
	return k.ToCel().ToFar()
}
