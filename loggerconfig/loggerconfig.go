// Logging configuration used by cli/ binaries.
package loggerconfig

import "os"
import "github.com/op/go-logging"

func Use() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{level:4.4s}%{id:4.4x}%{color:reset} %{message}\n%{color}%{level:4.4s}%{id:4.4x}%{color:reset} %{module} / %{shortfile} / %{longfunc}()\n\n",
		)))
}
