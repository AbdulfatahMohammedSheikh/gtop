package contsants

// contains the const dirs that will be used to read data from
const (
	BASEPATH   = "/proc/"
	CMDLINE    = "/cmdline"
	STATUSPATH = "/status"
	STAT       = "stat"
	UPTIME     = "/uptime"
	MEMINFO    = "/meminfo"
	VERION     = "version"

	// System clock ticks per second
	SYSTEMClkTck = 100

	// Index of the start time field in /proc/pid/STAT
	STARTTIMEIndex = 21 // this may requied to be deleted cuase i use 13-24 range
)

// TODO: add enum for info about the cpu status

/**
enum CPUStates {
    S_USER = 1,
    S_NICE,
    S_SYSTEM,
    S_IDLE,
    S_IOWAIT,
    S_IRQ,
    S_SOFTIRQ,
    S_STEAL,
    S_GUEST,
    S_GUEST_NICE
};
*/
