package indicators

import (
	"errors"

	"github.com/technicalviking/sliceWindow"
)

//Aroon port of Aroon indicator from tulipindicators c lib.
//@todo migrate this, and other ported go versions to my tulipindicators go lib.
func Aroon(highs []float64, lows []float64, period int) ([][]float64, error) {

	if len(highs) != len(lows) {
		return nil, errors.New("inputs must be same length")
	}

	inputLength := len(highs)

	if inputLength < period {
		return nil, errors.New("inputs length must be equal to or greater than period")
	}

	offset := period - 1

	aroonDown := make([]float64, len(lows)-offset)
	aroonUp := make([]float64, len(lows)-offset)

	lowsWindow := sliceWindow.New(period)
	highsWindow := sliceWindow.New(period)

	for i := 0; i < inputLength; i++ {
		lowsWindow.PushBack(lows[i])
		highsWindow.PushBack(highs[i])

		if i < offset {
			continue
		}

		/*
			Aroon-Up = ((Period - Days Since Period High)/14) x 100
			Aroon-Down = ((Period - Days Since Period Low)/14) x 100
		*/

		daysSinceMin := period - lowsWindow.MinPosition()
		daysSinceMax := period - highsWindow.MaxPosition()

		aroonDown[i-offset] = float64(period-daysSinceMin) * (float64(100) / float64(period))
		aroonUp[i-offset] = float64(period-daysSinceMax) * (float64(100) / float64(period))
		//aroonDown[i-offset] = (float64(lowsWindow.MinPosition()) / 14) * 100
		//aroonUp[i-offset] = (float64(highsWindow.MaxPosition()) / 14) * 100

	}

	return [][]float64{
		aroonDown,
		aroonUp,
	}, nil

}

//AroonOsc Aroon Oscillator
func AroonOsc(highs []float64, lows []float64, period int) ([]float64, error) {
	aroonIndicator, err := Aroon(highs, lows, period)

	if err != nil {
		return nil, err
	}

	output := make([]float64, len(aroonIndicator[0]))

	for i := 0; i < len(aroonIndicator[0]); i++ {
		output[i] = aroonIndicator[1][i] - aroonIndicator[0][i]
	}

	return output, nil
}
