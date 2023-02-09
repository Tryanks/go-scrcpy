package scrcpy

import (
	"errors"
	"os/exec"
)

type ScrClient struct {
	DeviceId   string
	DeviceInfo AdbDeviceInfo
	AliveCmd   *exec.Cmd
	AVCodec    string
}

func NewScrcpy(serial ...string) (sc ScrClient, err error) {
	if err = CheckAdb(); err != nil {
		return
	}
	if len(serial) > 1 {
		err = errors.New("only one serial number is allowed")
		return
	}
	if len(serial) == 0 {
		devices, err := AdbGetDevices()
		if err != nil {
			return
		}
		if len(devices) == 0 {
			err = errors.New("no device connected")
			return
		}
		serial = append(serial, devices[0])
	}
	deviceId := serial[0]
	device, err := AdbGetDeviceInfo(deviceId)
	if err != nil {
		err = errors.New("can not connect to device: " + deviceId)
		return
	}
	sc = ScrClient{
		DeviceId:   deviceId,
		DeviceInfo: device,
	}
	return
}
