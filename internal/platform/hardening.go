//go:build !openbsd

// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/platform/hardening.go

package platform

func PledgeCLI() error {
	return nil
}

func PledgeNetworkClient() error {
	return nil
}

func PledgeServer() error {
	return nil
}

func UnveilNone() error {
	return nil
}

func UnveilReadOnly(path string) error {
	return nil
}

func UnveilReadWrite(path string) error {
	return nil
}
