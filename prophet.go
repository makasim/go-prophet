package prophet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Prophet struct {
	changePointPriorScale float64
	changePointRange      float64
	intervalWidth         float64

	df DataFrame
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

func New(opts ...Option) *Prophet {
	p := &Prophet{
		changePointPriorScale: 0.05,
		changePointRange:      0.8,
		intervalWidth:         0.8,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Prophet) Forecast(df DataFrame) ([]Forecast, error) {
	in, err := json.Marshal(df)
	if err != nil {
		return nil, fmt.Errorf("data frame: json marsal: %w", err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	outDec := json.NewDecoder(buf)

	cmd := exec.Command("python3", "forecast.py")
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

type DataFrame []DataPoint

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
