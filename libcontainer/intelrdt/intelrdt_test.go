// +build linux

package intelrdt

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestIntelRdtSetL3CacheSchema(t *testing.T) {
	if !IsCatEnabled() {
		return
	}

	helper := NewIntelRdtTestUtil(t)
	defer helper.cleanup()

	const (
		l3CacheSchemaBefore = "L3:0=f;1=f0"
		l3CacheSchemeAfter  = "L3:0=f0;1=f"
	)

	helper.writeFileContents(map[string]string{
		"schemata": l3CacheSchemaBefore + "\n",
	})

	helper.IntelRdtData.config.IntelRdt.L3CacheSchema = l3CacheSchemeAfter
	intelrdt := &IntelRdtManager{
		Config: helper.IntelRdtData.config,
		Path:   helper.IntelRdtPath,
	}
	if err := intelrdt.Set(helper.IntelRdtData.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}
	values := strings.Split(tmpStrings, "\n")
	value := values[0]

	if value != l3CacheSchemeAfter {
		t.Fatal("Got the wrong value, set 'schemata' failed.")
	}
}

func TestIntelRdtSetMemBwSchema(t *testing.T) {
	if !IsMbaEnabled() {
		return
	}

	helper := NewIntelRdtTestUtil(t)
	defer helper.cleanup()

	const (
		memBwSchemaBefore = "MB:0=20;1=70"
		memBwSchemeAfter  = "MB:0=70;1=20"
	)

	helper.writeFileContents(map[string]string{
		"schemata": memBwSchemaBefore + "\n",
	})

	helper.IntelRdtData.config.IntelRdt.MemBwSchema = memBwSchemeAfter
	intelrdt := &IntelRdtManager{
		Config: helper.IntelRdtData.config,
		Path:   helper.IntelRdtPath,
	}
	if err := intelrdt.Set(helper.IntelRdtData.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}
	values := strings.Split(tmpStrings, "\n")
	value := values[0]

	if value != memBwSchemeAfter {
		t.Fatal("Got the wrong value, set 'schemata' failed.")
	}
}

func TestIntelRdtSetMemBwScSchema(t *testing.T) {
	if !IsMbaScEnabled() {
		return
	}

	helper := NewIntelRdtTestUtil(t)
	defer helper.cleanup()

	const (
		memBwScSchemaBefore = "MB:0=5000;1=7000"
		memBwScSchemeAfter  = "MB:0=9000;1=4000"
	)

	helper.writeFileContents(map[string]string{
		"schemata": memBwScSchemaBefore + "\n",
	})

	helper.IntelRdtData.config.IntelRdt.MemBwSchema = memBwScSchemeAfter
	intelrdt := &IntelRdtManager{
		Config: helper.IntelRdtData.config,
		Path:   helper.IntelRdtPath,
	}
	if err := intelrdt.Set(helper.IntelRdtData.config); err != nil {
		t.Fatal(err)
	}

	tmpStrings, err := getIntelRdtParamString(helper.IntelRdtPath, "schemata")
	if err != nil {
		t.Fatalf("Failed to parse file 'schemata' - %s", err)
	}
	values := strings.Split(tmpStrings, "\n")
	value := values[0]

	if value != memBwScSchemeAfter {
		t.Fatal("Got the wrong value, set 'schemata' failed.")
	}
}

func TestFindIntelRdtMountpointDir(t *testing.T) {
	//isMbaScEnabled = false
	testError := "Expected: %q but got: %q"

	t.Run("Valid mountinfo with MBA Software Controller disabled", func(t *testing.T) {
		f, err := os.Open("testing/mountinfo")
		if err != nil {
			t.Fatal(err)
		}

		mountpoint, err := findIntelRdtMountpointDir(f)
		if err != nil {
			t.Fatalf("Failed to find intel rdt mountpoint directory: %s", err)
		}

		if mountpoint != "/sys/fs/resctrl" {
			t.Fatal("Got the wrong value of Intel RDT mountpoint directory")
		}

		if isMbaScEnabled == true {
			t.Fatal("Intel RDT MBA Software Controller should be disabled")
		}
	})

	t.Run("Valid mountinfo with MBA Software Controller enabled", func(t *testing.T) {
		f, err := os.Open("testing/mountinfo_mba_sc")
		if err != nil {
			t.Fatal(err)
		}

		mountpoint, err := findIntelRdtMountpointDir(f)
		if err != nil {
			t.Fatalf("Failed to find intel rdt mountpoint directory: %s", err)
		}

		if mountpoint != "/sys/fs/resctrl" {
			t.Fatal("Got the wrong value of Intel RDT mountpoint directory")
		}

		if isMbaScEnabled == false {
			t.Fatal("Intel RDT MBA Software Controller should be enabled")
		}
	})

	t.Run("Empty mountinfo", func(t *testing.T) {
		expectedNotFoundError := errors.New("mountpoint for Intel RDT not found")

		_, err := findIntelRdtMountpointDir(strings.NewReader(""))
		if err != nil {
			if err.Error() != expectedNotFoundError.Error() {
				t.Fatalf(testError, expectedNotFoundError, err)
			}
		} else {
			t.Fatalf(testError, expectedNotFoundError, err)
		}
	})

	t.Run("Invalid mountinfo", func(t *testing.T) {
		expectedNotFoundError := errors.New("mountpoint for Intel RDT not found : Parsing '.' failed: not enough fields (1)")

		_, err := findIntelRdtMountpointDir(strings.NewReader("."))
		if err != nil {
			if err.Error() != expectedNotFoundError.Error() {
				t.Fatalf(testError, expectedNotFoundError, err)
			}
		} else {
			t.Fatalf(testError, expectedNotFoundError, err)
		}
	})
}
