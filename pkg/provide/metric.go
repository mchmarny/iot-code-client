package provide

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// NewMetricProvider creates nee MetricProvider
func NewMetricProvider(template event.MetricTemplate) MetricProvider {
	return MetricProvider{
		template: template,
	}
}

// MetricProvider generates metric readers based on dynamic value
type MetricProvider struct {
	template event.MetricTemplate
}

// GetParam returns local template
func (p *MetricProvider) GetParam() event.MetricTemplate {
	return p.template
}

// Provide provides os process events
func (p *MetricProvider) Provide(ctx context.Context, src string, h func(e *event.MetricReading)) error {
	ticker := time.NewTicker(p.template.Frequency)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			e, err := makeMetric(src, p.template)
			if err != nil {
				return err
			}
			h(e)
		}
	}
}

func makeMetric(src string, rp event.MetricTemplate) (e *event.MetricReading, err error) {
	v, ge := getRandomValue(&rp.Template)
	if ge != nil {
		return nil, errors.Wrap(ge, "error generating rundom value")
	}

	e = &event.MetricReading{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: rp.Label,
		Unit:  rp.Unit,
		Data:  v,
	}
	return
}

func getRandomValue(arg *event.ValueTemplate) (val interface{}, err error) {
	switch arg.Type {
	case "int", "int8", "int32", "int64":
		return getRandomIntValue(arg)
	case "float", "float32", "float64":
		return getRandomFloatValue(arg)
	case "bool":
		return getRandomBoolValue(), nil
	default:
		return nil, fmt.Errorf("invalid data type in template: %s", arg.Type)
	}
}

func getRandomIntValue(arg *event.ValueTemplate) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	min, err := toInt64(arg.Min)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid min int: %v", arg.Min)
	}
	max, err := toInt64(arg.Max)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid max int: %v", arg.Max)
	}
	return int64(rand.Intn(int(max)-int(min)) + int(min)), nil
}

func toInt64(v interface{}) (int64, error) {
	s := fmt.Sprintf("%v", v)
	return strconv.ParseInt(s, 10, 64)
}

func getRandomFloatValue(arg *event.ValueTemplate) (float64, error) {
	rand.Seed(time.Now().UnixNano())
	min, err := toFloat64(arg.Min)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid min int: %v", arg.Min)
	}
	max, err := toFloat64(arg.Max)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid max int: %v", arg.Max)
	}
	return min + rand.Float64()*(max-min), nil
}

func toFloat64(v interface{}) (float64, error) {
	s := fmt.Sprintf("%v", v)
	return strconv.ParseFloat(s, 64)
}

func getRandomBoolValue() bool {
	rand.Seed(time.Now().UnixNano())
	return (rand.Intn(100-1) + 1) < 50
}
