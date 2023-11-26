## go-prophet

A tiny wrapper for Meta's [prophet](https://facebook.github.io/prophet/) library.

### Installation

Python is required to run the prophet library.

```bash
pip install -r requirements.txt
python3 --version
```

### Usage

```go
package main

import "github.com/rylans/go-prophet"

func main() {
	p := prophet.New()

	fs, err := p.Forecast(df)
	if err != nil {
		panic(err)
	}

	log.Println(fs)
}
```

### License

MIT
