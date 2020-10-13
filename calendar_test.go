package calendar

import (
	"testing"
	"time"
)

//cd to this folder
//go test -v
// -v for verbose to show "t.Log("whatever msg")"

func TestIsTime1WithinTime2(t *testing.T) {
	start1 := time.Now()
	end1 := start1.Add(time.Hour)

	start2 := start1
	end2 := end1

	if !IsTime1WithinTime2(start1, end1, start2, end2) {
		t.Error("Expect true but got ", false)
	}
}

func TestMergeThreeTimesIntoOne(t *testing.T) {

	timeNow := time.Now()

	ranges1 := []TimeRange{
		{StartTime: timeNow, EndTime: timeNow.Add(20 * time.Minute)},
		{StartTime: timeNow.Add(21 * time.Minute), EndTime: timeNow.Add(30 * time.Minute)},
	}
	ranges2 := []TimeRange{
		{StartTime: timeNow.Add(20 * time.Minute), EndTime: timeNow.Add(21 * time.Minute)},
	}

	//log.Println("--------------------Range1")
	//
	//for _,e:=range ranges1{
	//	log.Println(e.StartTime.Format("03:04pm"),e.EndTime.Format("03:04pm"))
	//}
	//
	//log.Println("--------------------Range2")
	//
	//for _,e:=range ranges2{
	//	log.Println(e.StartTime.Format("03:04pm"),e.EndTime.Format("03:04pm"))
	//}

	merged := MergeTimeRangeList(ranges1, ranges2)

	if len(merged) != 1 {
		t.Error("Merged fail: should be len 1 but get len ", len(merged))
	}

	//for _, e := range merged {
	//	log.Println(e.StartTime.Format("03:04pm"), e.EndTime.Format("03:04pm"))
	//}
}

func TestMergeThreeTimesIntoTwo(t *testing.T) {

	timeNow := time.Now()

	ranges1 := []TimeRange{
		{StartTime: timeNow, EndTime: timeNow.Add(20 * time.Minute)},
		{StartTime: timeNow.Add(21 * time.Minute), EndTime: timeNow.Add(30 * time.Minute)},
	}
	ranges2 := []TimeRange{
		{StartTime: timeNow.Add(20*time.Minute + 1*time.Second), EndTime: timeNow.Add(21 * time.Minute)},
	}

	merged := MergeTimeRangeList(ranges1, ranges2)

	if len(merged) != 2 {
		t.Error("Merged fail: should be len 2 but get len ", len(merged))
	}
}
