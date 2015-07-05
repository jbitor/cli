// Logging configuration used by cli/ binaries.
package loggerconfig

import "os"
import "github.com/op/go-logging"

func Use() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{level:4.4s}%{color:reset}%{id:4.4x}%{color} %{message}\n%{color}%{level:4.4s}%{color:reset}%{id:4.4x}%{color} %{module:24.24s} %{shortfile:-24.24s} %{longfunc:32.32s}()\n",
		)))
}
