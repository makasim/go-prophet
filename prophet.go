package prophet

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

//go:embed py/forecast.py
var forecastPy []byte

type Prophet struct {
	changePointPriorScale  float64
	changePointRange       float64
	intervalWidth          float64
	futureDataFramePeriods int
	futureDataFrameFreq    string
}

type Option func(p *Prophet)

func WithChangePointPriorScale(val float64) Option {
	return func(p *Prophet) {
		p.changePointPriorScale = val
	}
}

func WithChangePointRange(val float64) Option {
	return func(p *Prophet) {
		p.changePointRange = val
	}
}

func WithIntervalWidth(val float64) Option {
	return func(p *Prophet) {
		p.intervalWidth = val
	}
}

func WithFutureDataFramePeriods(val int) Option {
	return func(p *Prophet) {
		p.futureDataFramePeriods = val
	}
}

func WithFutureDataFrameFreq(val string) Option {
	return func(p *Prophet) {
		p.futureDataFrameFreq = val
	}
}

func New(opts ...Option) *Prophet {
	p := &Prophet{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Prophet) Forecast(df []DataPoint) ([]Forecast, error) {
	in, err := json.Marshal(df)
	if err != nil {
		return nil, fmt.Errorf("data frame: json marsal: %w", err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	outDec := json.NewDecoder(buf)

	args := []string{`-c`, string(forecastPy)}
	if p.changePointPriorScale != 0 {
		args = append(args, fmt.Sprintf(`--changepoint_prior_scale=%f`, p.changePointPriorScale))
	}
	if p.changePointRange != 0 {
		args = append(args, fmt.Sprintf(`--changepoint_range=%f`, p.changePointPriorScale))
	}
	if p.intervalWidth != 0 {
		args = append(args, fmt.Sprintf(`--interval_width=%f`, p.intervalWidth))
	}
	if p.futureDataFramePeriods != 0 {
		args = append(args, fmt.Sprintf(`--future_dataframe_periods=%d`, p.futureDataFramePeriods))
	}
	if p.futureDataFrameFreq != `` {
		args = append(args, fmt.Sprintf(`--future_dataframe_freq=%s`, p.futureDataFrameFreq))
	}

	cmd := exec.Command("python3", args...)
	cmd.Env = os.Environ()
	cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("exec: python: %w", err)
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		return nil, fmt.Errorf("python exit code: %d", exitCode)
	}

	fs := make([]Forecast, 0)
	for outDec.More() {
		f := Forecast{}
		if err := outDec.Decode(&f); err != nil {
			return nil, fmt.Errorf("output: json: decode:%w", err)
		}

		fs = append(fs, f)
	}

	return fs, nil
}

type DataPoint struct {
	Ds string  `json:"ds"`
	Y  float64 `json:"y"`
}

type Forecast struct {
	Ds                       int64   `json:"ds"`
	Trend                    float64 `json:"trend"`
	YhatLower                float64 `json:"yhat_lower"`
	YhatUpper                float64 `json:"yhat_upper"`
	TrendLower               float64 `json:"trend_lower"`
	TrendUpper               float64 `json:"trend_upper"`
	AdditiveTerms            float64 `json:"additive_terms"`
	AdditiveTermsLower       float64 `json:"additive_terms_lower"`
	AdditiveTermsUpper       float64 `json:"additive_terms_upper"`
	Daily                    float64 `json:"daily"`
	DailyLower               float64 `json:"daily_lower"`
	DailyUpper               float64 `json:"daily_upper"`
	MultiplicativeTerms      float64 `json:"multiplicative_terms"`
	MultiplicativeTermsLower float64 `json:"multiplicative_terms_lower"`
	MultiplicativeTermsUpper float64 `json:"multiplicative_terms_upper"`
	Yhat                     float64 `json:"yhat"`
}
