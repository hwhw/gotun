package gotun

import (
	"syscall"
	"os"
	"unsafe"
)

type IfreqFlags struct {
    Name [16] byte
    Flags uint16
}

func getTunTap(flags uint16) (tunfile *os.File, tundevice string, err os.Error) {
	var ifr IfreqFlags
	ifr.Flags = flags

	tunfile, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0666)
	if err != nil {
		return
	}
	_, _, errnop := syscall.Syscall(syscall.SYS_IOCTL, uintptr(tunfile.Fd()), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
	errno := int(errnop)
	if errno != 0 {
		err = os.NewError(syscall.Errstr(errno))
		return
	}
	strlen := 16
	for ; strlen > 0; strlen-- {
		if ifr.Name[strlen-1] != 0 {
			break
		}
	}
	tundevice = string(ifr.Name[0:strlen])
	return
}

func NewTun() (tunfile *os.File, tundevice string, err os.Error) {
	tunfile, tundevice, err = getTunTap(syscall.IFF_TUN | syscall.IFF_NO_PI)
	return
}

func NewTap() (tunfile *os.File, tundevice string, err os.Error) {
	tunfile, tundevice, err = getTunTap(syscall.IFF_TAP | syscall.IFF_NO_PI)
	return
}
