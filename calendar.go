package calendar

import (
	"sort"
	"time"
)

type TimeRange struct {
	StartTime time.Time
	EndTime   time.Time
}

type ByEarliest []TimeRange

func (a ByEarliest) Len() int           { return len(a) }
func (a ByEarliest) Less(i, j int) bool { return a[i].StartTime.Before(a[j].StartTime) }
func (a ByEarliest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func IsTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	//same time is overlap. Eg. End at 2pm and Start at 2pm is overlap
	//return (start1.Before(end2)||start1.Equal(end2)) && (start2.Before(end1)||start2.Equal(end1)

	return (start1.Before(end2)) && (start2.Before(end1))
}

func IsTime1WithinTime2(start1, end1, start2, end2 time.Time) bool {

	//same time is within as well
	return (start1.After(start2) || start1.Equal(start2)) && (end1.Before(end2) || end1.Equal(end2))
}

//Assume both rangelist is sorted
func MergeTimeRangeList(ranges1, ranges2 []TimeRange) []TimeRange {
	if len(ranges1) == 0 {
		return ranges2
	} else if len(ranges2) == 0 {
		return ranges1
	}

	return mergeTimeRangeList(ranges1, ranges2)
}

//Assume ranges1 is earliest
func mergeTimeRangeList(ranges1, ranges2 []TimeRange) []TimeRange {

	len1 := len(ranges1)
	len2 := len(ranges2)

	merged := make([]TimeRange, 0, len1+len2)

	merged = append(ranges1, ranges2...)

	sort.Sort(ByEarliest(merged))

	//log.Println("--------------------SORTED")
	//for _,e:=range merged{
	//	log.Println(e.StartTime.Format("03:04pm"),e.EndTime.Format("03:04pm"))
	//}
	//log.Println("--------------------")

	mergedFlattened := make([]TimeRange, 0, len1+len2)

	mergedFlattened = append(mergedFlattened, merged[0])
	currentTimeRange := merged[0]

	for i := 0; i < len(merged); i++ {

		//log.Println("@@@@@@@@@@@@@@@@")

		currentEnd := currentTimeRange.EndTime

		nextStart := merged[i].StartTime
		nextEnd := merged[i].EndTime

		if currentEnd.After(nextStart) {
			//log.Println("1. After",currentEnd.Format("03:04pm"),nextStart.Format("03:04pm"))

			currentTimeRange.EndTime = maxTime(currentEnd, nextEnd)
		} else if currentEnd.Equal(nextStart) {
			//log.Println("2. Equal",currentEnd.Format("03:04pm"),nextStart.Format("03:04pm"))

			mergedFlattened[len(mergedFlattened)-1].EndTime = maxTime(currentEnd, nextEnd)
		} else {
			//log.Println("3. add to output")

			currentTimeRange = merged[i]

			if mergedFlattened[len(mergedFlattened)-1].EndTime == currentTimeRange.StartTime {
				mergedFlattened[len(mergedFlattened)-1].EndTime = currentTimeRange.EndTime
			} else {
				mergedFlattened = append(mergedFlattened, currentTimeRange)

			}
		}
	}

	return mergedFlattened
}

func ArrangeOverlap(merged []TimeRange) []TimeRange {
	if len(merged) == 0 {
		return []TimeRange{}
	}

	sort.Sort(ByEarliest(merged))

	//log.Println("--------------------SORTED")
	//for _,e:=range merged{
	//	log.Println(e.StartTime.Format("03:04pm"),e.EndTime.Format("03:04pm"))
	//}
	//log.Println("--------------------")

	mergedFlattened := make([]TimeRange, 0, len(merged))
	mergedFlattened = append(mergedFlattened, merged[0])

	for _, t := range merged {
		if len(mergedFlattened) == 0 || mergedFlattened[len(merged)-1].EndTime.Before(t.StartTime) {
			mergedFlattened = append(mergedFlattened, t)
		} else {
			mergedFlattened[len(merged)-1].EndTime = maxTime(mergedFlattened[len(merged)-1].EndTime, t.EndTime)
		}
	}

	//Previous Code (To be remove).
	//currentTimeRange := merged[0]
	//
	//for i := 0; i < len(merged); i++ {
	//
	//	//log.Println("@@@@@@@@@@@@@@@@")
	//
	//	currentEnd := currentTimeRange.EndTime
	//
	//	nextStart := merged[i].StartTime
	//	nextEnd := merged[i].EndTime
	//
	//	if currentEnd.After(nextStart) {
	//		//log.Println("1. After",currentEnd.Format("03:04pm"),nextStart.Format("03:04pm"))
	//
	//		currentTimeRange.EndTime = maxTime(currentEnd, nextEnd)
	//	} else if currentEnd.Equal(nextStart) {
	//		//log.Println("2. Equal",currentEnd.Format("03:04pm"),nextStart.Format("03:04pm"))
	//
	//		mergedFlattened[len(mergedFlattened)-1].EndTime = maxTime(currentEnd, nextEnd)
	//	} else {
	//		//log.Println("3. add to output")
	//
	//		currentTimeRange = merged[i]
	//
	//		if mergedFlattened[len(mergedFlattened)-1].EndTime == currentTimeRange.StartTime {
	//			mergedFlattened[len(mergedFlattened)-1].EndTime = currentTimeRange.EndTime
	//		} else {
	//			mergedFlattened = append(mergedFlattened, currentTimeRange)
	//
	//		}
	//	}
	//}

	return mergedFlattened
}

func maxTime(time1, time2 time.Time) time.Time {
	if time1.After(time2) {
		return time1
	}
	return time2
}
