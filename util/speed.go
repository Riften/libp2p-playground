package util

import (
	"context"
	"log"
	"sync"
	"time"
)

// Several utils used to analyse speed.

// SpeedStampReport Record speed with timestamp.
// To compute the speed between two SpeedStampReport r1, r2:
//		v = r2.value / (r2.stamp - r1.stamp)
type SpeedStampReport struct {
	stamp time.Time
	value int
}

type SpeedRecorder struct {
	ch chan *SpeedStampReport
	ctx context.Context

	// for speed compute
	reportCache []*SpeedStampReport
	cacheSize int
	cacheHead int
	totalValue int
	cacheLock sync.Mutex
}

// Please make sure that the cacheSize is larger than 2
func NewSpeedRecorder(ctx context.Context, cacheSize int) *SpeedRecorder {
	return &SpeedRecorder{
		ch:          make(chan *SpeedStampReport, 5),
		ctx:         ctx,
		reportCache: make([]*SpeedStampReport, cacheSize),
		cacheSize:   cacheSize,
		cacheHead:   0,
		totalValue:  0,
	}
}

// Run start() in independent routine
func (r *SpeedRecorder) Start() {
	var report *SpeedStampReport
	var oldReport *SpeedStampReport
	for{
		select{
		case report = <- r.ch:
			r.cacheLock.Lock()
			oldReport = r.reportCache[r.cacheHead]
			r.totalValue += report.value
			if oldReport != nil {
				r.totalValue -= oldReport.value
			}

			r.reportCache[r.cacheHead] = report
			r.cacheHead = (r.cacheHead + 1) % r.cacheSize
			if len(r.ch) > 1 {
				log.Println("WARN: report stumped in report channel. ")
			}
			r.cacheLock.Unlock()
		case <-r.ctx.Done():
			log.Println("WARN: SpeedRecorder end, ", r.ctx.Err())
			return
		}
	}
}

func (r *SpeedRecorder) AddStamp(stamp time.Time, value int) {
	r.ch <- &SpeedStampReport{
		stamp: stamp,
		value: value,
	}
}

// return value/millisecond in float64
// Note:
//		Nanosecond is toooo small.
func (r *SpeedRecorder) InstantV() float64 {
	r.cacheLock.Lock()
	//firstReport := r.reportCache[r.cacheHead]
	lastReport := r.reportCache[(r.cacheHead + r.cacheSize - 1) % r.cacheSize]
	last2Report := r.reportCache[(r.cacheHead + r.cacheSize - 2) % r.cacheSize]
	r.cacheLock.Unlock()

	if last2Report != nil && lastReport != nil {
		duration := lastReport.stamp.Sub(last2Report.stamp)
		if duration.Nanoseconds() <= 0 {
			log.Println("WARN: speed report with 0 nanoseconds duration")
			return 0
		}
		return (1000000*float64(lastReport.value))/float64(duration.Nanoseconds())
	} else {
		return 0
	}
}

func (r *SpeedRecorder) AverageV() float64 {
	r.cacheLock.Lock()
	firstReport := r.reportCache[r.cacheHead]
	lastReport := r.reportCache[(r.cacheHead + r.cacheSize - 1) % r.cacheSize]
	totalValue := r.totalValue
	r.cacheLock.Unlock()

	if firstReport != nil && lastReport != nil {
		duration := lastReport.stamp.Sub(firstReport.stamp)
		if duration.Nanoseconds() <= 0 {
			log.Println("WARN: speed report with 0 nanoseconds duration")
			return 0
		}
		return (1000000*float64(totalValue)) / float64(duration.Nanoseconds())
	} else {
		return 0
	}
}
