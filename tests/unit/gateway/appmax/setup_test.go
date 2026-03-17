package appmax_test

import (
	"os"
	"testing"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/tests/unit/support"
)

func TestMain(m *testing.M) {
	os.Exit(support.RunWithFrameworkBootstrap(m))
}
