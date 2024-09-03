package common

import (
	"fmt"
	"strconv"
	"time"
)

// transform the format "2006-01-02 15:04:05" to time.Time
func GetTimeFromFormatString(formatTimeString string) (time.Time, error) {
	var formatTime time.Time
	var err error
	formatTime, err = time.Parse("2006-01-02 15:04:05", formatTimeString)
	if err != nil {
		fmt.Println(err)
		return formatTime, err
	}
	return formatTime, nil
}

//
func CheckTimeHourAndMinute() (int, int) {
	time.LoadLocation("Asia/Shanghai")
	t := time.Now()
	return t.Hour(), t.Minute()
}

//
func PrintCurTime() {
	time.LoadLocation("Asia/Shanghai")
	fmt.Print(time.Now().Format("2006-01-02 15:04:05:\n"))
}

func GetNanoTimeStamp() int64 {
	timeUnixNano := time.Now().UnixNano() / 1000
	return timeUnixNano
}

func GetTime() string {
	now := time.Now()
	year, mon, day := now.UTC().Date()
	hour, min, sec := now.UTC().Clock()
	// zone, _ := now.UTC().Zone()
	time := fmt.Sprintf("%d-%d-%dT%02d:%02d:%02d.000Z", year, mon, day, hour, min, sec)
	return time
}

//
func CalcTimeCounter(period string, interval string) (int, error) {
	periodSecond, err := PeriodToSecond(period)
	if err != nil {
		return 0, fmt.Errorf("CalcTimeCounter error")
	}
	intervalSecont, err := PeriodToSecond(interval)
	if err != nil {
		return 0, fmt.Errorf("CalcTimeCounter error")
	}
	fmt.Printf("%d / %d = %d\n", periodSecond, intervalSecont, int(periodSecond/intervalSecont))
	return int(periodSecond / intervalSecont), nil
}

// get seconds from "*s", "*m", "*h" or "*d"
func PeriodToSecond(period string) (int, error) {
	if period == "" || len(period) <= 1 {
		fmt.Println("period is empty")
		return 0, fmt.Errorf("period is empty")
	}

	unit := period[len(period)-1]
	num, err := strconv.ParseInt(period[0:len(period)-1], 10, 32)
	if err != nil {
		fmt.Printf("%s isn't supported", period)
		return 0, fmt.Errorf("%s isn't supported", period)
	}

	var counter int64 = 0
	switch unit {
	case 's':
		counter = num
		break
	case 'm':
		counter = num * 60
		break
	case 'h':
		counter = num * 60 * 60
		break
	case 'd':
		counter = num * 60 * 60 * 24
		break
	default:
		return 0, fmt.Errorf("%s isn't supported", period)
	}

	return int(counter), nil
}
