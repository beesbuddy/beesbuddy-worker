package worker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/constants"
	"github.com/beesbuddy/beesbuddy-worker/internal/log"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nakabonne/tstorage"
)

func (w *workerComponent) storageWorker() {
	for msq := range w.queue {
		log.Debug.Println("try to persist metrics: ", msq)
		err := w.persist(msq)
		if err != nil {
			log.Error.Print(err)
		}
	}
}

func (w *workerComponent) persist(m metrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.WorkerTimeout)*time.Second)
	defer cancel()

	labels := []tstorage.Label{
		{Name: "apriaryId", Value: m.ApiaryId},
		{Name: "hiveId", Value: m.HiveId},
	}

	if err := w.writeLocalMetric("temperature", m.Temperature, labels); err != nil {
		log.Error.Println("unable to write temperature: ", err)
	}

	if err := w.writeLocalMetric("humidity", m.Humidity, labels); err != nil {
		log.Error.Println("unable to write humidity: ", err)
	}

	if err := w.writeLocalMetric("weight", m.Humidity, labels); err != nil {
		log.Error.Println("unable to write Weight: ", err)
	}

	if err := w.writeInfluxDbMetric(context.Background(), "temperature", m.Temperature, labels); err != nil {
		log.Error.Println("unable to write temperature: ", err)
	}

	if err := w.writeInfluxDbMetric(context.Background(), "humidity", m.Humidity, labels); err != nil {
		log.Error.Println("unable to write humidity: ", err)
	}

	if err := w.writeInfluxDbMetric(context.Background(), "weight", m.Weight, labels); err != nil {
		log.Error.Println("unable to write weight: ", err)
	}

	ctx.Done()
	return nil
}

func (w *workerComponent) writeLocalMetric(name, value string, labels []tstorage.Label) error {
	log.Debug.Printf("[LocalStorage] storing [%v] metric [%s] with value [%s]", labels, name, value)

	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return fmt.Errorf("unable to parse metric for: %v", w)
	} else {
		point := tstorage.Row{
			Metric:    name,
			Labels:    labels,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: v},
		}

		if err = w.storage.InsertRows([]tstorage.Row{
			point,
		}); err != nil {
			return fmt.Errorf("unable to write point: %s", err)
		}

		row := tstorage.Row{
			Metric:    name,
			Labels:    labels,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: v},
		}

		if err := w.storage.InsertRows([]tstorage.Row{row}); err != nil {
			return fmt.Errorf("unable to write row: %v", row)
		}
	}

	return nil
}

func (w *workerComponent) writeInfluxDbMetric(ctx context.Context, name, value string, labels []tstorage.Label) error {
	log.Debug.Printf("[InfluxDB] storing [%v] metric [%s] with value [%s]", labels, name, value)

	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return fmt.Errorf("unable to parse metric for: %v", w)
	} else {
		org := w.appCtx.Pref.GetConfig().InfluxDbOrg
		bucket := w.appCtx.Pref.GetConfig().InfluxDbBucket
		writeAPI := w.influxDbClient.WriteAPIBlocking(org, bucket)
		point := influxdb2.NewPointWithMeasurement(name)

		for _, label := range labels {
			point.AddTag(label.Name, label.Value)
		}

		point.AddField("value", v)

		writeAPI.WritePoint(ctx, point)

		if err := writeAPI.WritePoint(ctx, point); err != nil {
			return fmt.Errorf("unable to write point: %s", err)
		}
	}

	return nil
}
