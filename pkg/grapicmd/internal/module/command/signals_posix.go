// +build linux darwin

package command

import (
	"os"
	"syscall"
)

var (
	propagatableSignals = []os.Signal{
		os.Interrupt,
		// os.Kill,
		syscall.SIGABRT,
		syscall.SIGALRM,
		// syscall.SIGBUS,
		// syscall.SIGCHLD,
		syscall.SIGCONT,
		// syscall.SIGFPE,
		syscall.SIGHUP,
		syscall.SIGILL,
		syscall.SIGINT,
		syscall.SIGIO,
		syscall.SIGIOT,
		// syscall.SIGKILL,
		syscall.SIGPIPE,
		syscall.SIGPROF,
		syscall.SIGQUIT,
		// syscall.SIGSEGV,
		syscall.SIGSTOP,
		// syscall.SIGSYS,
		syscall.SIGTERM,
		syscall.SIGTRAP,
		syscall.SIGTSTP,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGURG,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGVTALRM,
		syscall.SIGWINCH,
		syscall.SIGXCPU,
		syscall.SIGXFSZ,
	}
)
