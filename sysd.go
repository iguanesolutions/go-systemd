package systemd

import "os"

// GetInvocationID returns the systemd invocation ID.
// If exists is false, we have not been launched by systemd.
func GetInvocationID() (ID string, exists bool) {
	return os.LookupEnv("INVOCATION_ID")
}
