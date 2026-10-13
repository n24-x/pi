package lifecycle

// StartHookRecords returns the info of OnStart hooks that successfully ran
// till the end, including their caller and execution time. Used to report
// timeout errors on Start.
func (l *Lifecycle) StartHookRecords() HookRecords {
	l.mu.Lock()
	defer l.mu.Unlock()
	r := make(HookRecords, len(l.startRecords))
	copy(r, l.startRecords)
	return r
}

// StopHookRecords returns the info of OnStop hooks that successfully ran
// till the end, including their caller and execution time. Used to report
// timeout errors on Stop.
func (l *Lifecycle) StopHookRecords() HookRecords {
	l.mu.Lock()
	defer l.mu.Unlock()
	r := make(HookRecords, len(l.stopRecords))
	copy(r, l.stopRecords)
	return r
}

// RunningHookCaller returns the name of the hook that was running when a Start/Stop
// hook timed out.
func (l *Lifecycle) RunningHookCaller() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.runningHook.callerFrame.Function
}
