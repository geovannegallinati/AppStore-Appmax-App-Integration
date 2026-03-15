package services_test

import (
	"os"
	"testing"

	"github.com/geovanne-gallinati/AppStoreAppDemo/tests/unit/support"
)

func TestMain(m *testing.M) {
	os.Exit(support.RunWithFrameworkBootstrap(m))
}
