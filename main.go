package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ContentDomainLabel = "content_domain"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	started = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iknow_cumulative_items_started",
		Help: "Cumulative started items",
	}, []string{ContentDomainLabel})

	studyTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iknow_cumulative_items_study_time_millis",
		Help: "Cumulative study time",
	}, []string{ContentDomainLabel})

	logHalflifeMillis = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iknow_cumulative_total_log_halflife_millis",
		Help: "Cumulative total log halflife milliseconds",
	}, []string{ContentDomainLabel})

	checkpoint1 = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iknow_cumulative_items_reached_checkpoint_1",
		Help: "Cumulative items that have reached checkpoint 1",
	}, []string{ContentDomainLabel})

	checkpoint2 = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iknow_cumulative_items_reached_checkpoint_2",
		Help: "Cumulative items that have reached checkpoint 2",
	}, []string{ContentDomainLabel})

	checkpoint3 = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iknow_cumulative_items_reached_checkpoint_3",
		Help: "Cumulative items that have reached checkpoint 3",
	}, []string{ContentDomainLabel})
)

type IknowExporter struct {
	Client IknowClient
}

func NewExporter(secret string) IknowExporter {
	return IknowExporter{
		Client: IknowClient{
			Secret: secret,
		},
	}
}

func (i IknowExporter) Update() error {
	stats, err := i.Client.GetCumulativeStats()

	if err != nil {
		return err
	}

	for contentDomain, s := range stats {
		started.WithLabelValues(contentDomain).Set(float64(s.Started))
		studyTime.WithLabelValues(contentDomain).Set(float64(s.TimeMillis))
		logHalflifeMillis.WithLabelValues(contentDomain).Set(float64(s.TotalLogHalflifeMillis))
		checkpoint1.WithLabelValues(contentDomain).Set(float64(s.Checkpoint1))
		checkpoint2.WithLabelValues(contentDomain).Set(float64(s.Checkpoint2))
		checkpoint3.WithLabelValues(contentDomain).Set(float64(s.Checkpoint3))
	}

	return nil
}

func (i IknowExporter) StartCollector() {
	go func() {
		for {
			err := i.Update()
			if err != nil {
				log.Printf("Error updating iknow stats: %v", err)
			}
			time.Sleep(time.Duration(10 * time.Minute))
		}
	}()
}

func main() {
	flag.Parse()

	token := os.Getenv("IKNOW_API_TOKEN")

	exporter := NewExporter(token)
	exporter.StartCollector()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
