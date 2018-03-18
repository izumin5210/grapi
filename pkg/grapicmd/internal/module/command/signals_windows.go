// +build windows

package command

import (
	"os"
	"syscall"
)

var (
	propagatableSignals = []os.Signal{
		os.Interrupt,
		syscall.SIGABRT,
		syscall.SIGALRM,
		// syscall.SIGBUS,
		// syscall.SIGCHLD,
		// syscall.SIGFPE,
		syscall.SIGHUP,
		syscall.SIGILL,
		syscall.SIGINT,
		// syscall.SIGKILL,
		syscall.SIGPIPE,
		syscall.SIGQUIT,
		// syscall.SIGSEGV,
		// syscall.SIGSYS,
		syscall.SIGTERM,
		syscall.SIGTRAP,
	}
)
