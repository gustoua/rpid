package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/parMaster/rpid/config"
)

type Response struct {
	TimeInState map[string]int
	LoadAvg     map[string][]float64
}

type SystemReporter struct {
	data Response
	dbg  bool
	mx   sync.Mutex
}

func LoadSystemReporter(cfg config.System, dbg bool) (*SystemReporter, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("SystemReporter is not enabled")
	}

	return &SystemReporter{
		dbg: dbg,
		data: Response{
			TimeInState: map[string]int{},
			LoadAvg:     map[string][]float64{},
		},
	}, nil
}

func (r *SystemReporter) Name() string {
	return "system"
}

func (r *SystemReporter) Collect() error {
	var err error
	r.mx.Lock()
	defer r.mx.Unlock()

	r.data.TimeInState, err = r.getCPUTimeInState(r.dbg)
	if err != nil {
		log.Printf("[ERROR] failed to get cpu time in state: %v", err)
	}

	la, err := r.getLoadAvg(r.dbg)
	if err != nil {
		log.Printf("[ERROR] failed to get load avg: %v", err)
	} else {
		r.data.LoadAvg["1m"] = append(r.data.LoadAvg["1m"], la["1m"])
	}
	return nil
}

func (r *SystemReporter) Report() (interface{}, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	return r.data, nil
}

func (r *SystemReporter) getCPUTimeInState(dbg bool) (map[string]int, error) {
	var (
		out  = map[string]int{}
		data []byte
		err  error
	)

	if dbg {
		data, err = os.ReadFile("cpu_time_in_state.txt")
	} else { // https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-class-thermal
		data, err = os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/stats/time_in_state")
	}

	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			continue
		}
		parts[0] = parts[0][0 : len(parts[0])-3]
		out[parts[0]], _ = strconv.Atoi(parts[1])
		out[parts[0]] /= 100 // tens of milliseconds to seconds
	}
	return out, nil
}

// Load average for last 1, 5 and 15 minutes
func (r *SystemReporter) getLoadAvg(dbg bool) (map[string]float64, error) {
	var (
		out  = map[string]float64{}
		data []byte
		err  error
	)

	if dbg {
		data, err = os.ReadFile("cpu_loadavg.txt")
	} else {
		data, err = os.ReadFile("/proc/loadavg")
	}

	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(data), " ")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid loadavg data")
	}

	for i, v := range []string{"1m", "5m", "15m"} {
		out[v], _ = strconv.ParseFloat(parts[i], 32)
	}
	return out, nil
}
