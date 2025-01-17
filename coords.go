package coords

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const SIZEOFREC int = 1024
const SEGMENT_START_TIME int64 = -4734072000
const SEGMENT_LAST_TIME int64 = 4735368000
const TOTAL_SUMMARIES_NUMBER int = 12 // we dont need the last two records 13 and 14
const INIT float64 = -4734072000.0

type FileRecords struct {
	Name               string
	Number             int //просто порядковый номер в файле
	TargetCode         int
	CenterCode         int
	RecordStartAddress int
	RecordLastAddress  int
	IntLen             float64
	RSize              float64
}

var DE440S_FILE_RECORDS = [TOTAL_SUMMARIES_NUMBER + 1]FileRecords{
	{
		Name:               "SSB",
		Number:             0,
		TargetCode:         0,
		CenterCode:         0,
		RecordStartAddress: 0,
		RecordLastAddress:  0,
		IntLen:             0,
		RSize:              0,
	},
	{
		Name:               "MERCURY_BARYCENTER",
		Number:             1,
		TargetCode:         1,
		CenterCode:         0,
		RecordStartAddress: 8065,
		RecordLastAddress:  610868,
		IntLen:             691200.0,
		RSize:              44.0,
	},
	{
		Name:               "VENUS_BARYCENTER",
		Number:             2,
		TargetCode:         2,
		CenterCode:         0,
		RecordStartAddress: 610869,
		RecordLastAddress:  830072,
		IntLen:             1382400.0,
		RSize:              32.0,
	},
	{
		Name:               "EARTH_BARYCENTER",
		Number:             3,
		TargetCode:         3,
		CenterCode:         0,
		RecordStartAddress: 830073,
		RecordLastAddress:  1110926,
		IntLen:             1382400.0,
		RSize:              41.0,
	},
	{
		Name:               "MARS_BARYCENTER",
		Number:             4,
		TargetCode:         4,
		CenterCode:         0,
		RecordStartAddress: 1110927,
		RecordLastAddress:  1230805,
		IntLen:             2764800.0,
		RSize:              35.0,
	},
	{
		Name:               "JUPITER_BARYCENTER",
		Number:             5,
		TargetCode:         5,
		CenterCode:         0,
		RecordStartAddress: 1230806,
		RecordLastAddress:  1319859,
		IntLen:             2764800.0,
		RSize:              26.0,
	},
	{
		Name:               "SATURN_BARYCENTER",
		Number:             6,
		TargetCode:         6,
		CenterCode:         0,
		RecordStartAddress: 1319860,
		RecordLastAddress:  1398638,
		IntLen:             2764800.0,
		RSize:              23.0,
	},
	{
		Name:               "URANUS_BARYCENTER",
		Number:             7,
		TargetCode:         7,
		CenterCode:         0,
		RecordStartAddress: 1398639,
		RecordLastAddress:  1467142,
		IntLen:             2764800.0,
		RSize:              20.0,
	},
	{
		Name:               "NEPTUNE_BARYCENTER",
		Number:             8,
		TargetCode:         8,
		CenterCode:         0,
		RecordStartAddress: 1467143,
		RecordLastAddress:  1535646,
		IntLen:             2764800.0,
		RSize:              20.0,
	},
	{
		Name:               "PLUTO_BARYCENTER",
		Number:             9,
		TargetCode:         9,
		CenterCode:         0,
		RecordStartAddress: 1535647,
		RecordLastAddress:  1604150,
		IntLen:             2764800.0,
		RSize:              20.0,
	},
	{
		Name:               "SUN",
		Number:             10,
		TargetCode:         10,
		CenterCode:         0,
		RecordStartAddress: 1604151,
		RecordLastAddress:  1843904,
		IntLen:             1382400.0,
		RSize:              35.0,
	},
	{
		Name:               "MOON",
		Number:             11,
		TargetCode:         301,
		CenterCode:         3,
		RecordStartAddress: 1843905,
		RecordLastAddress:  2967308,
		IntLen:             345600.0,
		RSize:              41.0,
	},
	{
		Name:               "EARTH",
		Number:             12,
		TargetCode:         399,
		CenterCode:         3,
		RecordStartAddress: 2967309,
		RecordLastAddress:  4090712,
		IntLen:             345600.0,
		RSize:              41.0,
	},
	// we dont need the following records
	// {
	//     Name: "MERCURY",
	//     Number: 13,
	//     TargetCode: 199,
	//     CenterCode: 1,
	//     RecordStartAddress: 4090713,
	//     RecordLastAddress: 4090724,
	// },
	// {
	//     Name: "VENUS",
	//     Number: 14,
	//     TargetCode: 299,
	//     CenterCode: 2,
	//     RecordStartAddress: 4090725,
	//     RecordLastAddress": 4090736,
	// },
}

// Расчёт координат положений планет на указанную дату (секунды с J2000)
// v2024
func GetCoordinates(dateInSeconds int64, targetCode int, centerCode int, bspFile *bytes.Reader) Position {

	if dateInSeconds <= SEGMENT_START_TIME || dateInSeconds >= SEGMENT_LAST_TIME {

		fmt.Println("getCoordinates: Date is out of range")
		return Position{X: 0.0, Y: 0.0, Z: 0.0}
	}

	var i = 0
	if centerCode == 0 {
		i = targetCode
	} else if centerCode == 3 {
		if targetCode == 301 {
			i = 11
			// print("get_coords: Center code is 301")
		} else if targetCode == 399 {
			// print("get_coords: Center code is 399")
			i = 12
		}
	} else {
		fmt.Println("getCoordinates: Date is out of range")
		return Position{X: 0.0, Y: 0.0, Z: 0.0}
	}

	var intLen = DE440S_FILE_RECORDS[i].IntLen
	var rSize = DE440S_FILE_RECORDS[i].RSize

	var internal_offset float64 = math.Floor((float64(dateInSeconds)-INIT)/intLen) * rSize
	var record = 8 * (DE440S_FILE_RECORDS[i].RecordStartAddress + int(internal_offset))

	var order = (int(rSize)-2)/3 - 1

	start_record := int64(record - 8)

	// byteData := make([]byte, (int)(rSize*8))
	// the only line that uses bsp file that is in memory
	// dataReader.ReadAt(byteData, start_record)
	// data := math.Float64frombits(binary.LittleEndian.Uint64(byteData))
	//data = array.array("d", dataReader[record-8:record+int(rSize)*8])
	data := make([]float64, (int)(rSize))
	bspFile.Seek(int64(start_record), io.SeekStart)

	for k := 0; k < int(rSize); k++ {
		var byteArr = make([]byte, 8)
		bspFile.Read(byteArr)
		data[k] = math.Float64frombits(binary.LittleEndian.Uint64(byteArr[0:]))

	}

	var tau = (float64(dateInSeconds) - data[0]) / data[1]
	order = int(order)
	var deg = order + 1

	return Position{
		X: chebyshev(order, tau, data[2:2+deg]),
		Y: chebyshev(order, tau, data[2+deg:2+2*deg]),
		Z: chebyshev(order, tau, data[2+2*deg:2+3*deg]),
	}

}

func chebyshev(order int, x float64, data []float64) float64 {

	// Evaluate a Chebyshev polynomial
	var bk float64
	two_x := 2 * x
	bkp2 := data[order]
	bkp1 := two_x*bkp2 + data[order-1]

	for n := order - 2; n > 0; n-- {
		bk = data[n] + two_x*bkp1 - bkp2
		bkp2 = bkp1
		bkp1 = bk
	}
	return data[0] + x*bkp1 - bkp2
}
