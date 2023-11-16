package tradefuncs

//mathfin

import (
	"math"

	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/stat"
)

type Iterator func(start int, end int, data []float64) []float64

// formula for average true range indicator
func AverageTrueRange(high, low, prevclose float64) float64 {
	return max((high - low), (low - prevclose), (high - prevclose))
}

func ExPostSqdRet(returns []float64) float64 {
	return stat.Variance(returns, nil)
}

// absRet is the absolute returns(directiion) of
// an asset over a certain period
func AbsRet(returns []float64) float64 {
	return math.Abs(stat.Variance(returns, nil))
}

func Covariance(stockAReturn, AvgStockAReturn, stockBReturn, AvgStockBReturn float64, sampleSize int) float64 {
	return ((stockAReturn - AvgStockAReturn) * (stockBReturn - AvgStockBReturn) / (float64(sampleSize) - 1))
}

func EMA(x []float64) float64 {
	smoothingConstant := func() float64 {
		denum := float64(len(x)) + 1
		return 2 / denum
	}()
	var ema float64

	// Calculate the weighting factor for each closing price using the smoothing constant.
	for i, price := range x {
		weightingFactor := smoothingConstant * math.Pow(1-smoothingConstant, float64(i))
		ema += price * weightingFactor
	}

	return ema
}

func Histogram(x []float64, y []float64) []float64 {
	histogram := make([]float64, 0)
	for i := 0; i < len(x); i++ {
		histogram = append(histogram, x[i]-y[i])
	}
	return histogram
}

func lowRange(x []float64) []float64 {
	percentile := math.Round(float64(len(x)) / float64(4))
	intPercentile := int(percentile)
	return x[:intPercentile]

}

func lowBound(x []float64) float64 {
	return max(x...)
}

func MA(x []float64, window int) []float64 {
	avgs := make([]float64, 0)
	j := window //choose window? 2,3,5,7 game window
	for i := 0; i <= len(x)-window; i++ {
		mn := stat.Mean(x[i:j], nil)
		avgs = append(avgs, mn)
		if j == len(x) && i == len(x)-window {
			break
		}
		j++
	}
	return avgs
}

func MACD(x []float64, xs int) float64 {
	MACD := EMA(x[:xs]) - EMA(x)
	return MACD
}

func MACDV(x []float64, high, low, prevclose float64) float64 {
	MACD := MACD(x, 12)

	numerator := MACD * 100
	return numerator / AverageTrueRange(high, low, prevclose)
}

func MACDVs(close, highs, lows []float64) []float64 {
	MACDVVal := make([]float64, 0)
	for j, i := 0, 26; i < len(close) && j < len(close); {
		val := MACDV(close[j:i], highs[i], lows[i], close[i-1])
		MACDVVal = append(MACDVVal, val)
		j++
		i++
	}
	return MACDVVal
}

func max(input ...float64) float64 {
	var max float64
	for _, i := range input {
		if i > max {
			max = i
		}
	}
	return max
}

func meanDiff(mean float64, x []float64) (above, below, equal float64) {
	var abvCnt, belCnt, eqCnt float64

	for _, v := range x {
		if v > mean {
			abvCnt++
		}
		if v == mean {
			eqCnt++
		}
		if v < mean {
			belCnt++
		}
	}

	percentCalc := func(num, denum float64) float64 {
		multip := num / denum
		return multip * 100
	}

	return percentCalc(abvCnt, float64(len(x))), percentCalc(belCnt, float64(len(x))), percentCalc(eqCnt, float64(len(x)))
}

func min(input ...float64) float64 {
	var min float64
	for _, i := range input {
		if i < min {
			min = i
		}
	}
	return min
}

// It gives the close returns of an asset
func Returns(currcls, prevcls float64) float64 {
	num := currcls - prevcls
	return num / prevcls
}

// It gives the close returns of an asset
func ReturnsMult(closes []float64) []float64 {
	multiReturn := make([]float64, 0, len(closes))
	for index := range closes {
		if index == 0 {
			continue
		}
		ret := Returns(closes[index], closes[index-1])
		multiReturn = append(multiReturn, ret)
	}
	return multiReturn
}

// returns quadratic variation of an asset of a given period
func QuadVar(R []float64) float64 {
	rets := make([]float64, 0)
	if len(R) < 2 {
		return 0
	}
	for i := range R {
		ret := math.Pow(R[i]-R[i-1], 2)
		rets = append(rets, ret)
	}
	return sum(rets)
}

// calculate average of the last 9 MACD's for signal lines
func signal(x []float64) []float64 {
	signal := make([]float64, 0)
	for j, i := 0, 9; i < len(x) && j < len(x); {
		signal = append(signal, EMA(x[j:i]))
		j++
		i++
	}
	return signal
}

func sort(x []float64) []float64 {
	newX := slices.Clone(x)
	slices.Sort(newX)
	return newX
}

func sum(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum = +v
	}
	return sum
}

func upBound(x []float64) float64 {
	return min(x...)
}
func UpRange(x []float64) []float64 {
	percentile := math.Round(float64(len(x)) / float64(4))
	intPercentile := int(percentile)
	return x[intPercentile:]
}

/*
	func iterate(data []float64, windowSize int, iterator Iterator) []float64 {
		result := make([]float64, 0)
		return result
	}

	func FVix(fVix, sVix float64) float64 {
		return fVix - sVix
	}

	func DRL(FVix float64, T float64) float64 {
		return FVix / T
	}

	func dailyReturns(x []float64) []float64 {
		arr := make([]float64, 0)
		for i := 0; i < len(x)-1; i++ {
			rec := x[i] - x[i+1]
			arr = append(arr, rec)
		}
		return arr
	}

	func contrarianStrat(x []float64) float64 {
		sum := func(x []float64) float64 {
			var newVal float64
			for _, val := range x {
				newVal = +val
			}
			return newVal
		}(dailyReturns(x))
		Rm := sum / float64(len(x))
		return Rm
	}
*/

// A wrapper for stat.StdScore to return an Array of standard score
func ZscoreCalc(x []float64, mean, stddev float64) []float64 {
	ZScrArr := make([]float64, 0)
	for _, v := range x {
		ZScr := stat.StdScore(v, mean, stddev)
		ZScrArr = append(ZScrArr, ZScr)
	}
	return ZScrArr
}
