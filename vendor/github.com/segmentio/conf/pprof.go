package conf

import "runtime"

// PPROF is a configuration struct which can be used to configure the runtime
// profilers of programs.
//
//	config := struct{
//		PPROF `conf:"pprof"`
//	}{
//		PPROF: conf.DefaultPPROF(),
//	}
//	conf.Load(&config)
//	conf.SetPPROF(config.PPROF)
//
type PPROF struct {
	BlockProfileRate     int `conf:"block-profile-rate"     help:"Sets the block profile rate to enable runtime profiling of blocking operations, zero disables block profiling." validate:"min=0"`
	MutexProfileFraction int `conf:"mutex-profile-fraction" help:"Sets the mutex profile fraction to enable runtime profiling of lock contention, zero disables mutex profiling." validate:"min=0"`
}

// DefaultPPROF returns the default value of a PPROF struct. Note that the
// zero-value is valid, DefaultPPROF differs because it captures the current
// configuration of the program's runtime.
func DefaultPPROF() PPROF {
	return PPROF{
		MutexProfileFraction: runtime.SetMutexProfileFraction(-1),
	}
}

// SetPPROF configures the runtime profilers based on the given PPROF config.
func SetPPROF(config PPROF) {
	runtime.SetBlockProfileRate(config.BlockProfileRate)
	runtime.SetMutexProfileFraction(config.MutexProfileFraction)
}
