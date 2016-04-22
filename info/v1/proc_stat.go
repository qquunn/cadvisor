package v1

import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"
)

type ProcCpuStats struct {
	CpuStats []CpuStat
}

func (self *ProcCpuStats) FindCpuStatsByNumber(cpuNumbers []int) []CpuStat {
	cpuNames := []string{}
	for _, cpuNumber := range cpuNumbers {
		cpuNames = append(cpuNames, "cpu" + strconv.Itoa(cpuNumber))
	}

	return self.FindCpuStatsByName(cpuNames)
}

func (self *ProcCpuStats) FindCpuStatsByName(cpuNames []string) []CpuStat {
	result := []CpuStat{}
	for _, cpustat := range self.CpuStats {
		if stringArrayContains(cpuNames, cpustat.Name) {
			result = append(result, cpustat)
		}
	}
	return result
}

func stringArrayContains(array []string, searchValue string) bool {
	for _, value := range array {
		if value == searchValue {
			return true
		}
	}

	return false
}

func ReadProcCpuStat() (ProcCpuStats, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return ProcCpuStats{}, err
	}
	defer f.Close()


	allCpuStat := ProcCpuStats{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		if !strings.HasPrefix(line, "cpu") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 11 {
			cpuStat := CpuStat {
				Name : fields[0],
				User : toUint64(fields[1]),
				Nice : toUint64(fields[2]),
				System : toUint64(fields[3]),
				Idle : toUint64(fields[4]),
				Iowait : toUint64(fields[5]),
				Irq : toUint64(fields[6]),
				Softirq : toUint64(fields[7]),
				Steal : toUint64(fields[8]),
				Guest : toUint64(fields[9]),
				GuestNice : toUint64(fields[10]),
			}

			allCpuStat.CpuStats = append(allCpuStat.CpuStats, cpuStat)

		}

	}

	return allCpuStat, nil
}

func toUint64(s string) uint64 {
	result, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}else {
		return result
	}
}
