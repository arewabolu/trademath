package tradefuncs

import (
	"testing"

	csv "github.com/arewabolu/csvmanager"
	"gonum.org/v1/gonum/stat"
)

func TestMax(t *testing.T) {
	max := max(1.0, 2.0, 3.0)
	if max != 3.0 {
		t.Error("Max failed")
	}
}

func TestAggregate(t *testing.T) {
	records, _ := csv.ReadCsv("./BTCUSDT-1h-2022-11.csv", true)
	high := records.Col("high").Float()
	//	open := records.Col("open").Float() //,,,

	low := records.Col("low").Float()
	close := records.Col("close").Float()
	_ = MACDVs(close, high, low)
}

func TestMeanDiff(t *testing.T) {
	rd, _ := csv.ReadCsv("./csvfiles/MACD.csv", true)
	MACDV := rd.Col("MACDV").Float()
	MACDVmean := stat.Mean(MACDV, nil)
	up, down, _ := meanDiff(MACDVmean, MACDV)
	if up == 0 || down == 0 {
		t.Error("MACDV data must be incorrect")
	}

}
