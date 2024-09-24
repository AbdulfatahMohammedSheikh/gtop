package progressparser

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sys-monitor/contsants"
	"sys-monitor/util"
	"unicode"
)

// TODO: order by VmSize
type Process struct {
	PID        string
	Cmd        string
	Scmd       string // shorter cmd
	VmSize     float64
	CpuPercent float64
	ProcUpTime string
	User       string
}

func NewProcessInfo() *[]Process {
	pids, _ := GetPidList()
	var processes []Process

	for _, pid := range *pids {
		cmd, _ := GetCmd(pid)
		t := strings.Split(*cmd, "/")
		sCmd := t[len(t)-1]

		vmSize, _ := GetVmSize(pid)
		cpuPercent, _ := GetCpuPercent(pid)
		temp, _ := GetProcUpTime(pid)
		procUptime := util.FormateTime(*temp)
		user, _ := GetProcUser(pid)
		newProcess := Process{
			PID:        pid,
			Cmd:        *cmd,
			Scmd:       sCmd,
			VmSize:     *vmSize,
			CpuPercent: *cpuPercent,
			ProcUpTime: procUptime,
			User:       *user,
		}
		processes = append(processes, newProcess)
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CpuPercent > processes[j].CpuPercent
	})

	return &processes
}

type SysInfo struct {
	// TODO: add info about cores + memory + swap
	Name                     string
	KernalVersion            string
	NumberOfProcesses        int
	NumberOfRunningProcesses int
	NumberOfCores            string
	UpTime                   string
	Processes                []Process
}

var Sysinfo SysInfo

func init() {
	osName, err := GetOsName()

	if nil != err {
		// TODO: log error
	}

	kernalVersion, err := GetSysKernelVersion()

	if nil != err {
		// TODO: log error
	}
	numberOfProcesses, err := GetTotalNumberOfProcesses()

	if nil != err {
		// TODO: log error
	}

	numberOfRunningProcesses, err := GetNumberOfRunningProcesses()

	if nil != err {
		// TODO: log error
	}

	numbrOfCores, err := GetNumberOfCores()

	if nil != err {
		// TODO: log error
	}

	temp, err := GetSysUpTime()

	if nil != err {
		// TODO: log error
	}
	uptime := util.FormateTime(*temp)

	processes := NewProcessInfo()

	Sysinfo = SysInfo{
		Name:                     *osName,
		KernalVersion:            *kernalVersion,
		NumberOfProcesses:        *numberOfProcesses,
		NumberOfRunningProcesses: *numberOfRunningProcesses,
		NumberOfCores:            *numbrOfCores,
		UpTime:                   uptime,
		Processes:                *processes,
	}
}

// TODO: make the function that get the scanner in utils

func GetCmd(pid string) (*string, error) {
	base := contsants.BASEPATH
	cmd := contsants.CMDLINE

	path := fmt.Sprintf("/%s/%s/%s", base, pid, cmd)
	_ = path
	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result string

	for scanner.Scan() {
		line := scanner.Text()
		result = line
	}

	return &result, nil
}

func GetPidList() (*[]string, error) {
	files, err := os.ReadDir("/proc/")
	pid := []string{}

	if nil != err {
		return nil, err

	}

	for _, file := range files {

		if file.IsDir() {

			name := file.Name()

			if unicode.IsDigit(rune(name[0])) && unicode.IsDigit(rune(name[len(name)-1])) {
				pid = append(pid, name)
			}

		}

	}

	return &pid, nil

}

func GetVmSize(pid string) (*float64, error) {

	// TODO: refactor this
	path := fmt.Sprintf("/proc/%s/status", pid)

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	target := "VmSize"
	var result float64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 2)

		key := strings.TrimSpace(parts[0])

		if key == target {

			// result will be like: 544800 kB

			// we ignore the last 3 elements which is ( KB)
			fields := strings.Fields(line)

			intVmData, err := strconv.Atoi(fields[1])

			if nil != err {
				return nil, err
			}

			// TODO: this retunr a float number with smae hex values like e+11

			// result = (intVmData) * (1024 * 1024) * 10

			result = float64(intVmData) / (1024 * 1024)

			break
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetCpuPercent(pid string) (*float64, error) {

	// FIXME: Calculation Method:
	//    Both your code and htop calculate CPU usage, but they might use different formulas or time intervals.
	//    Your code calculates CPU usage based on user, system, and child process times.
	//    htop might use a different approach, potentially including idle time or other factors

	path := fmt.Sprintf("/proc/%s/stat", pid)
	// First read
	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string

	for scanner.Scan() {
		line = scanner.Text()
	}

	fields := strings.Fields(line)

	utime, err := GetProcUpTime(pid)

	if nil != err {
		return nil, err
	}
	stime, err := strconv.ParseUint(fields[14], 10, 64)

	if nil != err {
		return nil, err
	}
	cutime, err := strconv.ParseUint(fields[15], 10, 64)

	if nil != err {
		return nil, err
	}
	cstime, err := strconv.ParseUint(fields[16], 10, 64)

	if nil != err {
		return nil, err
	}
	startTime, err := strconv.ParseUint(fields[21], 10, 64)
	_ = startTime

	if nil != err {
		return nil, err
	}

	uptime, err := GetSysUpTime()
	_ = uptime

	total_time := *utime + float64(stime) + float64(cutime) + float64(cstime)

	freq := float64(100) // Assuming 100 ticks per second

	seconds := *uptime - (float64(startTime) / freq)
	result := 100.0 * (total_time / freq) / seconds
	if nil != err {
		return nil, err
	}

	r := fmt.Sprintf("%.2f", result)
	feetFloat, _ := strconv.ParseFloat(strings.TrimSpace(r), 64)
	result = math.Abs(feetFloat)

	return &result, nil
}

// to have a formated time use [FormatedTime] function
func GetSysUpTime() (*float64, error) {

	base := contsants.BASEPATH
	upTime := contsants.UPTIME

	path := fmt.Sprintf("/%s/%s", base, upTime)
	_ = path
	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var l string

	for scanner.Scan() {
		line := scanner.Text()
		l = line
	}
	t := strings.Split(l, " ")

	result, err := strconv.ParseFloat(t[0], 64)

	if nil != err {
		return nil, err
	}

	return &result, nil

}

// it return unforamted time so use [formatedTime]  function to be more meaningful
func GetProcUpTime(pid string) (*float64, error) {

	base := contsants.BASEPATH
	stat := contsants.STAT

	path := fmt.Sprintf("/%s/%s/%s", base, pid, stat)

	// TODO:  get start time of the ststem
	sysyemStartTime, err := GetSysUpTime()

	if nil != err {
		return nil, err
	}

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string

	for scanner.Scan() {
		line = scanner.Text()
	}
	fields := strings.Fields(line)

	startTime, err := strconv.ParseUint(fields[21], 10, 64)
	if err != nil {

		return nil, err
	}

	procUptimeJiffies := int64(*sysyemStartTime) - int64(startTime)

	procUptimeSec := float64(procUptimeJiffies) / contsants.SYSTEMClkTck
	return &procUptimeSec, nil
}

func GetProcUser(pid string) (*string, error) {

	// proc/[PID]/status will tell use the uid

	base := contsants.BASEPATH
	stateusPath := contsants.STATUSPATH
	path := fmt.Sprintf("/%s/%s/%s", base, pid, stateusPath)

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var l string

	for scanner.Scan() {
		line := scanner.Text()
		if line[0:3] == "Uid" {

			l = line
			break
		}
	}

	// the result is a string like this : Uid: xxxx xxxx xxxx. To get the feilds of the string and ignore spaces we call feilds
	temp := strings.Fields(l)
	uid := temp[1] // first value after Uid:

	// /etc/passwd -> has the user name

	passFile, err := util.GetStream("/etc/passwd")

	if nil != err {

		return nil, err
	}

	defer file.Close()

	passFileScanner := bufio.NewScanner(passFile)
	target := fmt.Sprintf("x:%s", uid)
	var userName string

	for passFileScanner.Scan() {
		line := passFileScanner.Text()

		if strings.Contains(line, target) {
			n := strings.Split(line, ":")
			userName = n[0]
		}
	}

	return &userName, nil
}

func GetNumberOfCores() (*string, error) {
	// TODO: not sure but i will assume that number of cpu cores is the number
	base := contsants.BASEPATH
	path := fmt.Sprintf("/%s/cpuinfo", base)
	var cores int

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "processor") {
			cores += 1
		}
	}

	result := strconv.Itoa(cores)

	return &result, nil
}

// this function ruturn raw string not foramted
func GetSysCpuPercent(core int) (*string, error) {

	path := fmt.Sprintf("/%s/%s", contsants.BASEPATH, contsants.STAT)
	target := fmt.Sprintf("%s%d", "cpu", core)

	file, err := util.GetStream(path)

	if nil != err {

		return nil, err
	}

	defer file.Close()

	var result string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, target) {
			result = line
		}
	}

	return &result, nil
}

type CpuStats struct {
	User       uint64
	Nice       uint64
	System     uint64
	Idle       uint64
	Iowait     uint64
	Irq        uint64
	Softirq    uint64
	Steal      uint64
	Guest      uint64
	Guest_nice uint64
}

func ReadCpuStatus(status string) (*CpuStats, error) {

	fields := strings.Fields(status)
	cpuStats := *&CpuStats{}

	for i, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing /proc/stat value: %w", err)
		}
		switch i {
		case 0:
			cpuStats.User = value
		case 1:
			cpuStats.Nice = value
		case 2:
			cpuStats.System = value
		case 3:
			cpuStats.Idle = value
		case 4:
			cpuStats.Iowait = value
		case 5:
			cpuStats.Irq = value
		case 6:
			cpuStats.Softirq = value
		case 7:
			cpuStats.Steal = value
		case 8:
			cpuStats.Guest = value
		case 9:
			cpuStats.Guest_nice = value
		}

	}

	return &cpuStats, nil
}

// values cames from GetSysCpuPercent
func GetSysIdleCpuTime(values string) (*uint64, error) {
	// FIXME: this may not be neeeded since the ReadCpuStatus is used in other places which make it just repeting the job we did before here

	stat, err := ReadCpuStatus(values)

	if nil != err {
		return nil, err
	}

	return &stat.Idle, nil
}

// values cames from GetSysCpuPercent
func GetSysActiveCpuTime(values string) (*uint64, error) {
	// FIXME: this may not be neeeded since the ReadCpuStatus is used in other places which make it just repeting the job we did before here
	stat, err := ReadCpuStatus(values)

	if nil != err {
		return nil, err
	}

	result := stat.User + stat.Nice + stat.System + stat.Irq + stat.Softirq + stat.Guest + stat.Nice

	return &result, nil
}

// values cames from GetSysCpuPercent
func GetCpuTotlaUage(values string) (*float64, error) {
	// TODO: make this works
	total, err := GetSysActiveCpuTime(values)

	if nil != err {

	}

	idle, err := GetSysIdleCpuTime(values)

	if nil != err {
	}

	fields := strings.Fields(values)
	if len(fields) < 9 {
		// return 0, fmt.Errorf("invalid /proc/stat format")
	}

	cpuPercent := 100 * (float64(*total-*idle) / float64(*total))
	fmt.Println("prectage : ", cpuPercent)
	return &cpuPercent, nil
}

// static float getSysRamPercent();

type MemoryInfo struct {
	Total        uint64
	FreeMem      uint64
	AvailableMem uint64
	Buffers      uint64
}

func getMemonryInfo(line string, info *MemoryInfo) {

	// FIXME: check for :
	// called here feilds and in the caller also I called feilds

	fields := strings.Fields(line)
	key := fields[0]
	value, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		// return nil, err
	}

	switch key {
	case "MemTotal:":
		info.Total = value
	case "MemFree:":
		info.FreeMem = value
	case "MemAvailable:":
		info.AvailableMem = value
	case "Buffers:":
		info.Buffers = value
	}
}

// make sure to use ("%.2f") so the output will be nicly formated
func GetSysRamPercent() (*float64, error) {

	path := fmt.Sprintf("/%s/%s", contsants.BASEPATH, contsants.MEMINFO)

	file, err := util.GetStream(path)

	if nil != err {

		return nil, err
	}

	var totalMem uint64
	var memFree uint64
	var buffers uint64

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		key := fields[0]

		if key == "MemAvailable:" {
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if nil != err {
				return nil, err
			}
			totalMem = value
		} else if key == "MemFree:" {

			value, err := strconv.ParseUint(fields[1], 10, 64)
			if nil != err {
				return nil, err
			}
			memFree = (value)

		} else if key == "Buffers:" {

			value, err := strconv.ParseUint(fields[1], 10, 64)
			if nil != err {
				return nil, err
			}
			buffers = (value)

		}

	}

	result := 100.0 * (float64(memFree) / (float64(totalMem - buffers)))

	return &result, err
}

// static std::string getSysKernelVersion();
func GetSysKernelVersion() (*string, error) {

	path := fmt.Sprintf("/%s/%s", contsants.BASEPATH, contsants.VERION)

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var version string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Linux version") {
			feilds := strings.Fields(line)
			version = feilds[2]
			break
		}
	}

	return &version, nil
}

// static int getTotalThreads();

func GetTotalThreads() (*int, error) {
	target := "Threads:"
	_ = target
	threadCount := 0
	_ = threadCount
	pids, err := GetPidList()

	if nil != err {
		return nil, err
	}

	for _, pid := range *pids {

		path := fmt.Sprintf("/%s/%s/%s", contsants.BASEPATH, pid, contsants.STATUSPATH)

		file, err := util.GetStream(path)

		if nil != err {
			return nil, err
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, target) {
				fields := strings.Fields(line)

				result, err := strconv.Atoi(fields[1])
				if nil != err {

				}
				threadCount += result
				break
			}

		}
	}
	return &threadCount, nil

}

// static int getTotalNumberOfProcesses();
func GetTotalNumberOfProcesses() (*int, error) {

	target := "processes"
	path := fmt.Sprintf("/%s/%s", contsants.BASEPATH, contsants.STAT)
	var processesCount int

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, target) {
			fields := strings.Fields(line)
			result, err := strconv.Atoi(fields[1])
			if nil != err {

				return nil, err
			}

			processesCount = result
			break
		}
	}
	return &processesCount, nil
}

// static int getNumberOfRunningProcesses();
func GetNumberOfRunningProcesses() (*int, error) {
	// TODO: refactor this so getting the tottal GetNumberOfRunningProcesses( ) and GetTotalNumberOfProcesses is done in one function
	target := "procs_running"
	_ = target

	path := fmt.Sprintf("/%s/%s", contsants.BASEPATH, contsants.STAT)
	var count int

	file, err := util.GetStream(path)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, target) {

			fields := strings.Fields(line)
			result, err := strconv.Atoi(fields[1])
			if nil != err {

			}

			count = result
			break
		}
	}

	return &count, nil
}

// static string getOsName();
func GetOsName() (*string, error) {
	path := "/etc/os-release"
	file, err := util.GetStream(path)
	target := "PRETTY_NAME="

	if nil != err {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var osName string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, target) {
			start := strings.Index(line, "=") + 2 // to remove ="
			end := len(line) - 1
			osName = line[start:end]
			break
		}
	}

	return &osName, nil

}

// static std::string printCpuStats(std::vector<std::string> values1, std::vector<std::string>values2);
// TODO: make this work later
func PrintCpuStats() {

	// float activeTime = getSysActiveCpuTime(values2) - getSysActiveCpuTime(values1);
	// float idleTime = getSysIdleCpuTime(values2) - getSysIdleCpuTime(values1);
	// float totalTime = activeTime + idleTime;
	// float result = 100.0*(activeTime / totalTime);
	// return to_string(result);

}

// };
