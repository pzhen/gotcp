package Utils

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"syscall"
)

//由于linux/unix系统默认打开文件句柄数量有限
//此函数将修改程序运行时系统的打开句柄个数
//Mac下golang 版本 1.9.2可以,在高版本修改失败
func SetSysLimit()  {
	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", errors.WithStack(err))
		os.Exit(1)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", errors.WithStack(err))
		os.Exit(1)
	}
}
