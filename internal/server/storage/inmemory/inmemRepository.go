package inmemory

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var _ storage.MetricRepository = (*inmemMetricRepository)(nil)

type inmemMetricRepository struct {
	gauges          map[string]float64
	counters        map[string]int64
	archiveFileName string
	synchArchive    bool
	metricMx        sync.RWMutex
	archiveMx       sync.RWMutex
}

type InmemArchiveSettings struct {
	StoreInterval   time.Duration
	FileName        string
	RestoreOnCreate bool
}

type archiveMetrics struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

func New(archiveSettings InmemArchiveSettings, ctx context.Context) *inmemMetricRepository {
	synchArchive := archiveSettings.StoreInterval == 0

	repos := &inmemMetricRepository{
		archiveFileName: archiveSettings.FileName,
		synchArchive:    synchArchive,
	}

	if err := repos.initMetrics(archiveSettings.RestoreOnCreate); err != nil {
		log.Fatalf(err.Error())
	}

	if !synchArchive {
		go repos.runArchivator(archiveSettings.StoreInterval, ctx)
	}

	return repos
}

func (r *inmemMetricRepository) GetAll() ([]metric.Counter, []metric.Gauge) {
	r.metricMx.RLock()
	defer r.metricMx.RUnlock()

	counters := make([]metric.Counter, 0, len(r.counters))
	gauges := make([]metric.Gauge, 0, len(r.gauges))

	for n, v := range r.counters {
		counters = append(counters, metric.Counter{
			Name:  n,
			Value: v,
		})
	}

	for n, v := range r.gauges {
		gauges = append(gauges, metric.Gauge{
			Name:  n,
			Value: v,
		})
	}

	return counters, gauges
}

func (r *inmemMetricRepository) GetCounter(name string) (c metric.Counter, ok bool) {
	r.metricMx.RLock()
	defer r.metricMx.RUnlock()

	c.Name = name
	c.Value, ok = r.counters[c.Name]

	return
}

func (r *inmemMetricRepository) GetGauge(name string) (c metric.Gauge, ok bool) {
	r.metricMx.RLock()
	defer r.metricMx.RUnlock()

	c.Name = name
	c.Value, ok = r.gauges[c.Name]

	return
}

func (r *inmemMetricRepository) AddOrUpdateGauge(g metric.Gauge) {
	r.metricMx.Lock()
	r.gauges[g.Name] = g.Value
	r.metricMx.Unlock()

	if r.synchArchive {
		// TODO: Change archive logic to avoid resave whole archive file on every adding
		r.ArchiveAll()
	}
}

func (r *inmemMetricRepository) AddOrUpdateCounter(c metric.Counter) {
	r.metricMx.Lock()
	r.counters[c.Name] = c.Value
	r.metricMx.Unlock()

	if r.synchArchive {
		// TODO: Change archive logic to avoid resave whole archive file on every adding
		r.ArchiveAll()
	}
}

func (r *inmemMetricRepository) ArchiveAll() error {
	if r.archiveFileName == "" {
		return nil
	}

	r.metricMx.RLock()
	metricsToArchive := archiveMetrics{
		Gauges:   r.gauges,
		Counters: r.counters,
	}
	bytesToArchive, err := json.Marshal(metricsToArchive)
	if err != nil {
		return err
	}
	r.metricMx.RUnlock()

	r.archiveMx.Lock()
	defer r.archiveMx.Unlock()

	return os.WriteFile(r.archiveFileName, bytesToArchive, 0777)
}

func (r *inmemMetricRepository) initMetrics(restoreOnCreate bool) error {
	if r.archiveFileName == "" || !restoreOnCreate {
		r.gauges = make(map[string]float64)
		r.counters = make(map[string]int64)

		return nil
	}

	return r.restoreArchive()
}

func (r *inmemMetricRepository) restoreArchive() error {
	fileBytes, err := os.ReadFile(r.archiveFileName)
	isNotExists := err != nil && errors.Is(err, os.ErrNotExist)
	isEmptyFile := err == nil && (len(fileBytes) == 0 || len(fileBytes) == 2) // empty or '{}'

	if isNotExists || isEmptyFile {
		r.gauges = make(map[string]float64)
		r.counters = make(map[string]int64)
		log.Println("nothing to restore from repository archive, skipping")

		return nil
	} else if err != nil {

		return err
	}

	restoredMetrics := &archiveMetrics{}
	err = json.Unmarshal(fileBytes, restoredMetrics)
	if err != nil {
		return err
	}

	r.gauges = restoredMetrics.Gauges
	r.counters = restoredMetrics.Counters

	log.Println("metrics were restored from repository archive")

	return nil
}

func (r *inmemMetricRepository) runArchivator(storeInterval time.Duration, ctx context.Context) {
	ticker := time.NewTicker(storeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.ArchiveAll(); err != nil {
				log.Printf("couldn't archive metrics: %s", err.Error())
			}
		case <-ctx.Done():

			return
		}
	}
}
