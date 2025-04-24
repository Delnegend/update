package utils

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetExeVersion(filePath string) (string, error) {
	size, err := windows.GetFileVersionInfoSize(filePath, nil)
	if err != nil || size == 0 {
		return "", errors.New("can't get version info size")
	}
	buf := make([]byte, size)
	err = windows.GetFileVersionInfo(filePath, 0, uint32(len(buf)), unsafe.Pointer(&buf[0]))
	if err != nil {
		return "", errors.New("can't get version info")
	}
	var block *uint16
	var blockLen uint32
	err = windows.VerQueryValue(unsafe.Pointer(&buf[0]), `\`, unsafe.Pointer(&block), &blockLen)
	if err != nil {
		return "", errors.New("can't query version value")
	}
	// VS_FIXEDFILEINFO structure
	type VS_FIXEDFILEINFO struct {
		Signature        uint32
		StrucVersion     uint32
		FileVersionMS    uint32
		FileVersionLS    uint32
		ProductVersionMS uint32
		ProductVersionLS uint32
		FileFlagsMask    uint32
		FileFlags        uint32
		FileOS           uint32
		FileType         uint32
		FileSubtype      uint32
		FileDateMS       uint32
		FileDateLS       uint32
	}
	info := (*VS_FIXEDFILEINFO)(unsafe.Pointer(block))
	ms := info.FileVersionMS
	ls := info.FileVersionLS
	version := fmt.Sprintf("%d.%d.%d.%d",
		ms>>16, ms&0xFFFF, ls>>16, ls&0xFFFF)
	return version, nil
}
