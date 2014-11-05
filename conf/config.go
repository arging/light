// Copyright 2014 li. All rights reserved.

package conf

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// Config is the representation of configuration settings.
// If the corresponding value not found, find from the environment.
type Config map[string]string

func (c Config) String(key string, defaultv string) string {
	if v, ok := c[key]; ok {
		return v
	}

	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultv
}

func (c Config) Bool(key string, defaultv bool) bool {
	switch strings.ToLower(c.String(key, "")) {
	case "y", "on", "1", "yes":
		return true
	case "n", "off", "0", "no":
		return false
	default:
		return defaultv
	}
}

func (c Config) Float(key string, defaultv float64) float64 {
	v, err := strconv.ParseFloat(c.String(key, ""), 64)
	if err == nil {
		return v
	} else {
		return defaultv
	}
}

func (c Config) Int(key string, defaultv int) int {
	v, err := strconv.Atoi(c.String(key, ""))
	if err == nil {
		return v
	} else {
		return defaultv
	}
}

func read(fname string) (Config, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	config := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		i := strings.Index(line, "=")
		config[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i+1:])
	}

	return Config(config), nil
}
