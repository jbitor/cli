// Logging configuration used by cli/ binaries.
package loggerconfig

import "os"
import "github.com/op/go-logging"

func Use() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{module:10.10s} %{level:4.4s}%{color:reset}%{id:4.4x}%{color} %{message}\n%{color}%{module:10.10s} %{level:4.4s}%{color:reset}%{id:4.4x}%{color}     %{shortfile:-20.20s} %{longfunc:32.32s}()%{color:reset}",
		)))
}
