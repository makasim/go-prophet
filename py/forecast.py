import pandas as pd
from prophet import Prophet
import argparse
import sys
import json

import logging
logger = logging.getLogger('cmdstanpy')
logger.addHandler(logging.NullHandler())
logger.propagate = False
logger.setLevel(logging.CRITICAL)

json_data = sys.stdin.read()
data = json.loads(json_data)
df = pd.DataFrame(data)
df.columns = ["ds", "y"]

parser = argparse.ArgumentParser(description='Prophet Time Series Forecasting')
parser.add_argument('--changepoint_prior_scale', type=float, default=0.05)
parser.add_argument('--changepoint_range', type=float, default=0.8)
parser.add_argument('--interval_width', type=float, default=0.80)
parser.add_argument('--future_dataframe_periods', type=int, default=1)
parser.add_argument('--future_dataframe_freq', type=str, default='10s')
args = parser.parse_args()

m = Prophet(
    changepoint_prior_scale=args.changepoint_prior_scale,
    changepoint_range=args.changepoint_range,
    interval_width=args.interval_width)
m.fit(df)

future = m.make_future_dataframe(periods=args.future_dataframe_periods, freq=args.future_dataframe_freq)
forecast = m.predict(future)

# Convert the forecast results to JSON and write to stdout
forecast_json = forecast.to_json(orient='records', lines=True)
print(forecast_json)
