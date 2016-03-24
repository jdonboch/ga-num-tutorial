package main

import (
	"fmt"
)

type Gene byte

func (g Gene) IsNumber() bool {
	return g < 12
}

func (g Gene) IsOperator() bool {
	return g >= 12
}

func (g Gene) IsPlus() bool {
	return g == 12
}

func (g Gene) IsMinus() bool {
	return g == 13
}

func (g Gene) IsMultiply() bool {
	return g == 14
}

func (g Gene) IsDivide() bool {
	return g == 15
}

func (g Gene) GetValue() int {
	return int(g)
}

func (g Gene) GetFloatValue() float64 {
	return float64(g)
}

func (g Gene) String() string {
	switch g {
	case 12:
		return "+"
	case 13:
		return "-"
	case 14:
		return "/"
	case 15:
		return "*"
	default:
		return fmt.Sprintf("%.2f", g.GetFloatValue())
	}
}

func (g Gene) Operate(x, y float64) float64 {
	switch g {
	case 12:
		return x + y
	case 13:
		return x - y
	case 14:
		return x / y
	case 15:
		return x * y
	default:
		return x
	}
}
