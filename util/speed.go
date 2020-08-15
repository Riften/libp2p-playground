package util

import (
	"context"
	"log"
	"sync"
	"time"
)

// Several utils used to analyse speed.

// SpeedStampReport Record and compute speed with timestamp.
type SpeedStampReport struct {
	stamp time.Time
	value int
}

type recorder struct {
	ch chan *SpeedStampReport
	ctx context.Context

	// for speed compute
	reportCache []*SpeedStampReport
	cacheSize int
	cacheHead int
	totalValue int
	cacheLock sync.Mutex
}

type SpeedRecorder struct {
	recorders map[string] *recorder
}

// Please make sure that the cacheSize is larger than 2
func newRecorder(ctx context.Context, cacheSize int) *recorder {
	return &recorder{
		ch:          make(chan *SpeedStampReport, 5),
		ctx:         ctx,
		reportCache: make([]*SpeedStampReport, cacheSize),
		cacheSize:   cacheSize,
		cacheHead:   0,
		totalValue:  0,
	}
}

// Run start() in independent routine
func (r *recorder) start() {
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

			r.cacheLock.Unlock()
		case <-r.ctx.Done():
			log.Println("WARN: recorder end, ", r.ctx.Err())
			return
		}
	}
}

// return value/nanoseconds in float64
func (r *recorder) instantV() float64 {
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
		return float64(lastReport.value)/float64(duration.Nanoseconds())
	} else {
		return 0
	}
}

func (r *recorder) averageV() float64 {
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
		return float64(totalValue) / float64(duration.Nanoseconds())
	} else {
		return 0
	}
}
