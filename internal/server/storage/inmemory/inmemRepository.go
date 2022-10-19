package inmemory

import (
	"encoding/json"
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
	archiveEnabled  bool
	archiveFilePath string
	mx              sync.RWMutex // корректно ли, что одна структура используется
	// для лока файла-архива и мап репозитория одновременно
}

type InmemArchiveSettings struct {
	StoreInterval   time.Duration
	FileName        string
	RestoreOnCreate bool
}

type archiveMetrics struct {
	gauges   map[string]float64
	counters map[string]int64
}

func New(archiveSettings InmemArchiveSettings) *inmemMetricRepository {
	repos := &inmemMetricRepository{
		archiveEnabled:  archiveSettings.StoreInterval == 0,
		archiveFilePath: archiveSettings.FileName,
	}

	if err := repos.initMetrics(archiveSettings.RestoreOnCreate); err != nil {
		log.Fatalf(err.Error())
	}

	// TODO: запускать в отдельной горутине сохранение в архив по таймеру (если не синхронно)
	// TODO: доработать синхронное сохранение архива

	return repos
}

func (r *inmemMetricRepository) GetAll() ([]metric.Counter, []metric.Gauge) {
	r.mx.RLock()
	defer r.mx.RUnlock()

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
	r.mx.RLock()
	defer r.mx.RUnlock()

	c.Name = name
	c.Value, ok = r.counters[c.Name]

	return
}

func (r *inmemMetricRepository) GetGauge(name string) (c metric.Gauge, ok bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	c.Name = name
	c.Value, ok = r.gauges[c.Name]

	return
}

func (r *inmemMetricRepository) AddOrUpdateGauge(g metric.Gauge) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.gauges[g.Name] = g.Value
}

func (r *inmemMetricRepository) AddOrUpdateCounter(c metric.Counter) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.counters[c.Name] = c.Value
}

func (r *inmemMetricRepository) ArchiveAll() error {
	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.archiveFilePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	metricsToArchive := archiveMetrics{
		gauges:   r.gauges,
		counters: r.counters,
	}
	bytesToArchive, err := json.Marshal(metricsToArchive)
	if err != nil {
		return err
	}

	return os.WriteFile(r.archiveFilePath, bytesToArchive, 0777)
}

func (r *inmemMetricRepository) initMetrics(restoreOnCreate bool) error {
	if !r.archiveEnabled || !restoreOnCreate {
		r.gauges = make(map[string]float64)
		r.counters = make(map[string]int64)

		return nil
	}

	return r.restoreArchive()
}

func (r *inmemMetricRepository) restoreArchive() error {
	fileBytes, err := os.ReadFile(r.archiveFilePath)
	if err != nil {
		return err
	}

	if len(fileBytes) == 0 {
		return nil
	}

	restoredMetrics := &archiveMetrics{}
	err = json.Unmarshal(fileBytes, restoredMetrics)
	if err != nil {
		return err
	}

	r.gauges = restoredMetrics.gauges
	r.counters = restoredMetrics.counters

	return nil
}
