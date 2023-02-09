package scrcpy

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

var adbPath = "adb"

func SetAdbPath(path string) error {
	adbPath = path
	return CheckAdb()
}

func CheckAdb() error {
	cmd := exec.Command(adbPath, "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	if !strings.Contains(out.String(), "Android Debug Bridge version") {
		return errors.New("adb not found, please check if Adb environment exists on your system. To specify the Adb path, use the SetAdbPath(path) method")
	}
	return nil
}

type AdbDeviceInfo struct {
	DeviceID    string // 设备ID
	ProductName string // 产品名称
	Model       string // 型号
}

// AdbGetDevices 获取所有已连接的设备
func AdbGetDevices() ([]string, error) {
	cmd := exec.Command(adbPath, "devices")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(out.String(), "\n")
	var devices []string
	for _, line := range lines[1:] {
		if line == "List of devices attached" || !strings.Contains(line, "\tdevice") {
			continue
		}
		deviceID := strings.Split(line, "\t")[0]
		devices = append(devices, deviceID)
	}
	return devices, nil
}

// AdbGetDeviceInfo 获取设备信息
func AdbGetDeviceInfo(deviceID string) (AdbDeviceInfo, error) {
	cmd := exec.Command(adbPath, "-s", deviceID, "shell", "getprop")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return AdbDeviceInfo{}, err
	}
	lines := strings.Split(out.String(), "\n")
	var deviceInfo AdbDeviceInfo
	for _, line := range lines {
		if strings.Contains(line, "ro.product.name") {
			deviceInfo.ProductName = strings.Split(line, ": ")[1]
		} else if strings.Contains(line, "ro.product.model") {
			deviceInfo.Model = strings.Split(line, ": ")[1]
		}
	}
	deviceInfo.DeviceID = deviceID
	return deviceInfo, nil
}

// AdbPush 将本地文件推送到设备
func AdbPush(deviceID, local, remote string) error {
	cmd := exec.Command(adbPath, "-s", deviceID, "push", local, remote)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	if !strings.Contains(out.String(), "pushed") {
		return errors.New("failed to push file")
	}
	return nil
}

// AdbForward 将设备上的端口转发到本地
func AdbForward(deviceID, local, remote string) error {
	cmd := exec.Command(adbPath, "-s", deviceID, "forward", local, remote)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// AdbForwardRemove 移除设备上的端口转发
func AdbForwardRemove(deviceID, local string) error {
	cmd := exec.Command(adbPath, "-s", deviceID, "forward", "--remove", local)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// AdbShell 执行设备上的shell命令
func AdbShell(deviceID, command string) error {
	if cmd, err := AdbShellAsync(deviceID, command); err != nil {
		return err
	} else {
		return cmd.Wait()
	}
}

// AdbShellAsync 异步执行设备上的shell命令
func AdbShellAsync(deviceID, command string) (*exec.Cmd, error) {
	cmd := exec.Command(adbPath, "-s", deviceID, "shell", command)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
