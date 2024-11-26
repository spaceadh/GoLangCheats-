package main

import (
	"log"
	"math"
	"math/rand"
)

// There are 2 main operation types: one with leaves and one without leaves (InputsCount() returns 0)
type Operation interface {
	// a number of leaves which produce inputs for this operation
	InputsCount() uint8
	// set children of the operation
	SetInputs([]Operation)
	// always returns 3 elements for RGB
	Eval(x float64, y float64) []float64
}

// OpVarX
type OpVarX struct{}

func (o *OpVarX) InputsCount() uint8 {
	return 0
}

func (o *OpVarX) SetInputs(inputs []Operation) {}

func (o *OpVarX) Eval(x float64, y float64) []float64 {
	return []float64{x, x, x}
}

// OpVarY
type OpVarY struct{}

func (o *OpVarY) InputsCount() uint8 {
	return 0
}

func (o *OpVarY) SetInputs(inputs []Operation) {}

func (o *OpVarY) Eval(x float64, y float64) []float64 {
	return []float64{y, y, y}
}

// OpConstant
type OpConstant struct {
	constant float64
}

func (o *OpConstant) InputsCount() uint8 {
	return 0
}

func (o *OpConstant) SetInputs(inputs []Operation) {}

func (o *OpConstant) Eval(x float64, y float64) []float64 {
	return []float64{o.constant, o.constant, o.constant}
}

// OpColorMix
type OpColorMix struct {
	inputs []Operation
}

func (o *OpColorMix) InputsCount() uint8 {
	return 3
}

func (o *OpColorMix) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func (o *OpColorMix) Eval(x float64, y float64) []float64 {
	r := o.inputs[0].Eval(x, y)[0]
	g := o.inputs[1].Eval(x, y)[1]
	b := o.inputs[2].Eval(x, y)[2]

	return []float64{r, g, b}
}

// OpCircle
type OpCircle struct {
	centerX float64
	centerY float64
}

func (o *OpCircle) InputsCount() uint8 {
	return 0
}

func (o *OpCircle) SetInputs(inputs []Operation) {}

func (o *OpCircle) Eval(x float64, y float64) []float64 {
	h := math.Hypot(x-o.centerX, y-o.centerY)
	return []float64{h, h, h}
}

// OpInverse
type OpInverse struct {
	inputs []Operation
}

func (o *OpInverse) InputsCount() uint8 {
	return 1
}

func (o *OpInverse) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func (o *OpInverse) Eval(x float64, y float64) []float64 {
	v := o.inputs[0].Eval(x, y)
	return []float64{1 - v[0], 1 - v[1], 1 - v[2]}
}

// OpSum
type OpSum struct {
	inputs []Operation
}

func (o *OpSum) InputsCount() uint8 {
	return 2
}

func (o *OpSum) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func (o *OpSum) Eval(x float64, y float64) []float64 {
	a := o.inputs[0].Eval(x, y)
	b := o.inputs[0].Eval(x, y)
	return []float64{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

// OpProduct
type OpProduct struct {
	inputs []Operation
}

func (o *OpProduct) InputsCount() uint8 {
	return 2
}

func (o *OpProduct) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func (o *OpProduct) Eval(x float64, y float64) []float64 {
	a := o.inputs[0].Eval(x, y)
	b := o.inputs[1].Eval(x, y)
	return []float64{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}

// OpMod
type OpMod struct {
	inputs []Operation
}

func (o *OpMod) InputsCount() uint8 {
	return 2
}

func (o *OpMod) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func (o *OpMod) Eval(x float64, y float64) []float64 {
	a := o.inputs[0].Eval(x, y)
	b := o.inputs[1].Eval(x, y)
	return []float64{math.Mod(a[0], b[0]), math.Mod(a[1], b[1]), math.Mod(a[2], b[2])}
}

// OpPerChannelMask
type OpPerChannelMask struct {
	inputs   []Operation
	constant float64
}

func (o *OpPerChannelMask) InputsCount() uint8 {
	return 3
}

func (o *OpPerChannelMask) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func (o *OpPerChannelMask) Eval(x float64, y float64) []float64 {
	a := o.inputs[0].Eval(x, y)
	b := o.inputs[1].Eval(x, y)
	c := o.inputs[2].Eval(x, y)

	var aa, bb, cc float64
	if a[0] > o.constant {
		aa = b[0]
	} else {
		aa = c[0]
	}
	if a[0] > o.constant {
		bb = b[1]
	} else {
		bb = c[1]
	}
	if a[0] > o.constant {
		cc = b[2]
	} else {
		cc = c[2]
	}
	return []float64{aa, bb, cc}
}

// OpBinaryMask
type OpBinaryMask struct {
	inputs   []Operation
	constant float64
}

func (o *OpBinaryMask) InputsCount() uint8 {
	return 3
}

func (o *OpBinaryMask) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func length(in []float64) float64 {
	var sum float64 = 0
	for _, v := range in {
		sum += v * v
	}
	return math.Sqrt(sum)
}

func (o *OpBinaryMask) Eval(x float64, y float64) []float64 {
	in := o.inputs[0].Eval(x, y)
	a := o.inputs[1].Eval(x, y)
	b := o.inputs[2].Eval(x, y)

	if length(in) > o.constant {
		return a
	}

	return b
}

// OpSmoothMix
type OpSmoothMix struct {
	inputs   []Operation
	constant float64
}

func (o *OpSmoothMix) InputsCount() uint8 {
	return 3
}

func (o *OpSmoothMix) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func smooth(length float64, a float64, b float64) float64 {
	return length*a + (1-length)*b
}

func (o *OpSmoothMix) Eval(x float64, y float64) []float64 {
	in := o.inputs[0].Eval(x, y)
	a := o.inputs[1].Eval(x, y)
	b := o.inputs[2].Eval(x, y)

	l := length(in)

	return []float64{smooth(l, a[0], b[0]), smooth(l, a[1], b[1]), smooth(l, a[2], b[2])}
}

// OpWell
type OpWell struct {
	inputs []Operation
}

func (o *OpWell) InputsCount() uint8 {
	return 1
}

func (o *OpWell) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func well(x float64) float64 {
	return math.Pow(1-2/(1+x*x), 8)
}

func (o *OpWell) Eval(x float64, y float64) []float64 {
	in := o.inputs[0].Eval(x, y)

	return []float64{well(in[0]), well(in[1]), well(in[2])}
}

// OpTent
type OpTent struct {
	inputs []Operation
}

func (o *OpTent) InputsCount() uint8 {
	return 1
}

func (o *OpTent) SetInputs(inputs []Operation) {
	o.inputs = inputs
}

func tent(x float64) float64 {
	return 1 - 2*math.Abs(x)
}

func (o *OpTent) Eval(x float64, y float64) []float64 {
	in := o.inputs[0].Eval(x, y)

	return []float64{tent(in[0]), tent(in[1]), tent(in[2])}
}

// operations factory
func pickOperation(prng *rand.Rand, depth int) Operation {
	opsNoLeaves := []string{"x", "y", "const", "circle"}
	opsWithLeaves := []string{"colormix", "inverse", "sum", "product", "mod", "perchanmask", "binarymask", "smoothmix", "well", "tent"}

	var opID string
	if depth > 1 {
		i := prng.Intn(len(opsWithLeaves) - 1)
		opID = opsWithLeaves[i]
	} else {
		i := prng.Intn(len(opsNoLeaves) - 1)
		opID = opsNoLeaves[i]
	}

	switch opID {
	case "x":
		return &OpVarX{}
	case "y":
		return &OpVarY{}
	case "const":
		return &OpConstant{constant: prng.Float64()}
	case "circle":
		return &OpCircle{centerX: prng.Float64(), centerY: prng.Float64()}
	case "colormix":
		return &OpColorMix{}
	case "inverse":
		return &OpInverse{}
	case "sum":
		return &OpSum{}
	case "product":
		return &OpProduct{}
	case "mod":
		return &OpMod{}
	case "perchanmask":
		return &OpPerChannelMask{constant: prng.Float64()}
	case "binarymask":
		return &OpBinaryMask{constant: prng.Float64()}
	case "smoothmix":
		return &OpSmoothMix{constant: prng.Float64()}
	case "well":
		return &OpWell{}
	case "tent":
		return &OpTent{}
	default:
		log.Fatalf("operation id %s is not valid", opID)
		return nil
	}
}