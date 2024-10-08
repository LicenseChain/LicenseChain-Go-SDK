
package licensechain

import "fmt"

type LicenseError struct {
    Message string
}

func (e *LicenseError) Error() string {
    return fmt.Sprintf("License Error: %s", e.Message)
}
