//go:build !remote

package generate

import (
	"fmt"
	"strings"
)

// ParseDevice parses device mapping string to a src, dest & permissions string
func ParseDevice(device string) (string, string, string, error) {
	var src string
	var dst string
	permissions := "rwm"
	arr := strings.Split(device, ":")
	switch len(arr) {
	case 3:
		if !IsValidDeviceMode(arr[2]) {
			if arr[2] == "" {
				return "", "", "", fmt.Errorf("in device mapping %q: device mode cannot be empty. Either omit the second : or specify a valid device mode (combination of r, w, m)", device)
			}
			return "", "", "", fmt.Errorf("in device mapping %q: invalid device mode %q. Valid device modes are combinations of r (read), w (write), and m (mknod)", device, arr[2])
		}
		permissions = arr[2]
		fallthrough
	case 2:
		if IsValidDeviceMode(arr[1]) {
			permissions = arr[1]
		} else {
			if len(arr[1]) > 0 && arr[1][0] != '/' {
				if arr[1] == "" {
					return "", "", "", fmt.Errorf("in device mapping %q: device mode cannot be empty. Either omit the : or specify a valid device mode (combination of r, w, m)", device)
				}
				return "", "", "", fmt.Errorf("in device mapping %q: invalid device mode %q. Valid device modes are combinations of r (read), w (write), and m (mknod)", device, arr[1])
			}
			dst = arr[1]
		}
		fallthrough
	case 1:
		src = arr[0]
	default:
		return "", "", "", fmt.Errorf("invalid device specification: %q. Expected format: /host/path[:/container/path[:mode]]", device)
	}

	if dst == "" {
		dst = src
	}
	return src, dst, permissions, nil
}

// IsValidDeviceMode checks if the mode for device is valid or not.
// IsValid mode is a composition of r (read), w (write), and m (mknod).
func IsValidDeviceMode(mode string) bool {
	var legalDeviceMode = map[rune]bool{
		'r': true,
		'w': true,
		'm': true,
	}
	if mode == "" {
		return false
	}
	for _, c := range mode {
		if !legalDeviceMode[c] {
			return false
		}
		legalDeviceMode[c] = false
	}
	return true
}
