//go:build !remote

package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDevice(t *testing.T) {
	tests := []struct {
		device string
		src    string
		dst    string
		perm   string
	}{
		{"/dev/foo", "/dev/foo", "/dev/foo", "rwm"},
		{"/dev/foo:/dev/bar", "/dev/foo", "/dev/bar", "rwm"},
		{"/dev/foo:/dev/bar:rw", "/dev/foo", "/dev/bar", "rw"},
		{"/dev/foo:rw", "/dev/foo", "/dev/foo", "rw"},
		{"/dev/foo::rw", "/dev/foo", "/dev/foo", "rw"},
	}
	for _, test := range tests {
		src, dst, perm, err := ParseDevice(test.device)
		assert.NoError(t, err)
		assert.Equal(t, src, test.src)
		assert.Equal(t, dst, test.dst)
		assert.Equal(t, perm, test.perm)
	}
}

func TestParseDeviceErrors(t *testing.T) {
	tests := []struct {
		device      string
		expectedErr string
	}{
		// Empty device mode scenarios (the main issue being fixed)
		{"/dev/fuse::", "in device mapping \"/dev/fuse::\": device mode cannot be empty. Either omit the second : or specify a valid device mode (combination of r, w, m)"},
		{"/dev/null::", "in device mapping \"/dev/null::\": device mode cannot be empty. Either omit the second : or specify a valid device mode (combination of r, w, m)"},
		{"/dev/null:foo:", "in device mapping \"/dev/null:foo:\": device mode cannot be empty. Either omit the second : or specify a valid device mode (combination of r, w, m)"},
		
		// Invalid device mode characters
		{"/dev/null:rd", "in device mapping \"/dev/null:rd\": invalid device mode \"rd\". Valid device modes are combinations of r (read), w (write), and m (mknod)"},
		{"/dev/null:xyz", "in device mapping \"/dev/null:xyz\": invalid device mode \"xyz\". Valid device modes are combinations of r (read), w (write), and m (mknod)"},
		{"/dev/null:x", "in device mapping \"/dev/null:x\": invalid device mode \"x\". Valid device modes are combinations of r (read), w (write), and m (mknod)"},
		
		// Too many colons
		{"/dev/null::::", "invalid device specification: \"/dev/null::::\". Expected format: /host/path[:/container/path[:mode]]"},
		{"/dev/null:::::", "invalid device specification: \"/dev/null:::::\". Expected format: /host/path[:/container/path[:mode]]"},
	}
	
	for _, test := range tests {
		t.Run(test.device, func(t *testing.T) {
			_, _, _, err := ParseDevice(test.device)
			assert.Error(t, err)
			assert.Equal(t, test.expectedErr, err.Error())
		})
	}
}
