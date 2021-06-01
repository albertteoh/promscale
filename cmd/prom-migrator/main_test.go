// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license.

package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timescale/promscale/pkg/migration-tool/utils"
)

func TestParseFlags(t *testing.T) {
	cases := []struct {
		name            string
		input           []string
		expectedConf    *config
		failsValidation bool
		errMessage      string
	}{
		{
			name:  "pass_normal",
			input: []string{"-start=1970-01-01T00:16:40+00:00", "-end=1970-01-01T00:16:41+00:00", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false"},
			expectedConf: &config{
				name:             "prom-migrator",
				start:            "1970-01-01T00:16:40+00:00",
				end:              "1970-01-01T00:16:41+00:00",
				mint:             1000000,
				mintSec:          1000,
				maxt:             1001000,
				maxtSec:          1001,
				humanReadable:    true,
				maxSlabSizeBytes: 524288000,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
				concurrentPull:     1,
				progressEnabled:    false,
			},
			failsValidation: false,
		}, {
			name: "pass_rt_info",
			input: []string{"-start=1970-01-01T00:16:40+00:00", "-end=1970-01-01T00:16:41+00:00", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false",
				"-reader-timeout=50m", "-reader-retry-delay=1m", "-reader-max-retries=10", "-reader-on-timeout=skip", "-reader-on-error=retry",
				"-writer-timeout=50m", "-writer-retry-delay=1m", "-writer-max-retries=15", "-writer-on-timeout=abort", "-writer-on-error=retry"},
			expectedConf: &config{
				name:             "prom-migrator",
				start:            "1970-01-01T00:16:40+00:00",
				end:              "1970-01-01T00:16:41+00:00",
				mint:             1000000,
				mintSec:          1000,
				maxt:             1001000,
				maxtSec:          1001,
				humanReadable:    true,
				maxSlabSizeBytes: 524288000,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout * 10,
					Delay:        time.Minute,
					OnTimeoutStr: "skip",
					OnErrStr:     "retry",
					MaxRetry:     10,
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout * 10,
					Delay:        time.Minute,
					OnTimeoutStr: "abort",
					OnErrStr:     "retry",
					MaxRetry:     15,
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
				concurrentPull:     1,
				progressEnabled:    false,
			},
			failsValidation: false,
		},
		{
			name:  "pass_normal with inverted commas",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false"},
			expectedConf: &config{
				name:             "prom-migrator",
				start:            "'1970-01-01T00:16:40+00:00'",
				end:              "'1970-01-01T00:16:41+00:00'",
				mint:             1000000,
				mintSec:          1000,
				maxt:             1001000,
				maxtSec:          1001,
				humanReadable:    true,
				maxSlabSizeBytes: 524288000,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
				concurrentPull:     1,
				progressEnabled:    false,
			},
			failsValidation: false,
		},
		{
			name:  "pass_normal_size",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-max-read-size=100MB"},
			expectedConf: &config{
				name:    "prom-migrator",
				start:   "'1970-01-01T00:16:40+00:00'",
				end:     "'1970-01-01T00:16:41+00:00'",
				mint:    1000000,
				mintSec: 1000,
				maxt:    1001000,
				maxtSec: 1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				concurrentPull:     1,
				maxSlabSizeBytes:   104857600,
				humanReadable:      true,
				maxSlabSize:        "100MB",
				concurrentPush:     1,
				progressEnabled:    false,
			},
			failsValidation: false,
		},
		{
			name:  "pass_normal_size_with_space",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-max-read-size=100 MB"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				concurrentPull:     1,
				maxSlabSizeBytes:   104857600,
				maxSlabSize:        "100 MB",
				concurrentPush:     1,
				progressEnabled:    false,
			},
			failsValidation: false,
		},
		{
			name:  "pass_normal_size_with_concurrent_implements",
			input: []string{"-start=1000", "-end=1001", "-human-readable-time=false", "-progress-enabled=false", "-concurrent-pull=16", "-concurrent-push=8", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "1000",
				end:           "1001",
				humanReadable: false,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				concurrentPull:     16,
				maxSlabSizeBytes:   524288000,
				maxSlabSize:        "500MB",
				concurrentPush:     8,
				progressEnabled:    false,
			},
			failsValidation: false,
		},
		{
			name:  "fail_normal_size",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-max-read-size=100MBB"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				maxSlabSize:        "100MBB",
				concurrentPush:     1,
				concurrentPull:     1,
				progressEnabled:    false,
			},
			failsValidation: true,
			errMessage:      `parsing byte-size: Unrecognized size suffix MBB`,
		},
		{
			name:  "fail_invalid_suffix",
			input: []string{"-start=1000", "-end=1001", "-human-readable-time=false", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-max-read-size=100PP"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "1000",
				end:           "1001",
				humanReadable: false,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				maxSlabSize:        "100PP",
				concurrentPush:     1,
				progressEnabled:    false,
				concurrentPull:     1,
			},
			failsValidation: true,
			errMessage:      `parsing byte-size: Unrecognized size suffix PP`,
		},
		{
			name:  "pass_normal_regex",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-progress-metric-name=progress_migration_up"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "progress_migration_up",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    false,
				concurrentPull:     1,
				maxSlabSizeBytes:   524288000,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: false,
		},
		{
			name:  "fail_invalid_regex",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-progress-metric-name=_progress_migration-_up"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "_progress_migration-_up",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    false,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `invalid metric-name regex match: prom metric must match ^[a-zA-Z_:][a-zA-Z0-9_:]*$: recieved: _progress_migration-_up`,
		},
		{
			name:  "fail_invalid_regex_2",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-progress-metric-name=0_progress_migration_up"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "0_progress_migration_up",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    false,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `invalid metric-name regex match: prom metric must match ^[a-zA-Z_:][a-zA-Z0-9_:]*$: recieved: 0_progress_migration_up`,
		},
		{
			name:  "fail_no_mint",
			input: []string{},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         defaultStartTime,
				end:           time.Now().Format(time.RFC3339),
				humanReadable: true,
				maxt:          time.Now().Unix() * 1000,
				maxtSec:       time.Now().Unix(),
				readerClient: utils.ClientRuntime{
					URL:          "",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `mint should be provided for the migration to begin`,
		},
		{
			name:  "fail_all_default",
			input: []string{"-start=1", "-human-readable-time=false"},
			expectedConf: &config{
				name:    "prom-migrator",
				start:   "1",
				end:     fmt.Sprintf("%d", time.Now().Unix()),
				mint:    1000,
				mintSec: 1,
				maxt:    time.Now().Unix() * 1000,
				maxtSec: time.Now().Unix(),
				readerClient: utils.ClientRuntime{
					URL:          "",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `remote read storage url and remote write storage url must be specified. Without these, data migration cannot begin`,
		},
		{
			name:  "fail_all_default_space",
			input: []string{"-start='1970-01-01T00:00:01+00:00'", "-reader-url=  ", "-writer-url= "},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:00:01+00:00'",
				end:           time.Now().Format(time.RFC3339),
				humanReadable: true,
				mint:          1000,
				mintSec:       1,
				maxt:          time.Now().Unix() * 1000,
				maxtSec:       time.Now().Unix(),
				readerClient: utils.ClientRuntime{
					URL:          "  ",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          " ",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `remote read storage url and remote write storage url must be specified. Without these, data migration cannot begin`,
		},
		{
			name:  "fail_empty_read_url",
			input: []string{"-start='1970-01-01T00:00:01+00:00'", "-writer-url=http://localhost:9201/write"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:00:01+00:00'",
				end:           time.Now().Format(time.RFC3339),
				humanReadable: true,
				mint:          1000,
				mintSec:       1,
				maxt:          time.Now().Unix() * 1000,
				maxtSec:       time.Now().Unix(),
				readerClient: utils.ClientRuntime{
					URL:          "",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `remote read storage url needs to be specified. Without read storage url, data migration cannot begin`,
		},
		{
			name:  "fail_empty_write_url",
			input: []string{"-start='1970-01-01T00:00:01+00:00'", "-reader-url=http://localhost:9090/api/v1/read"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:00:01+00:00'",
				end:           time.Now().Format(time.RFC3339),
				humanReadable: true,
				mint:          1000,
				mintSec:       1,
				maxt:          time.Now().Unix() * 1000,
				maxtSec:       time.Now().Unix(),
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `remote write storage url needs to be specified. Without write storage url, data migration cannot begin`,
		},
		{
			name:  "fail_mint_greater_than_maxt",
			input: []string{"-start='2001-09-09T01:46:41+00:00'", "-end='2001-09-09T01:46:40+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'2001-09-09T01:46:41+00:00'",
				end:           "'2001-09-09T01:46:40+00:00'",
				humanReadable: true,
				mint:          1000000001000,
				mintSec:       1000000001,
				maxt:          1000000000000,
				maxtSec:       1000000000,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `invalid input: minimum timestamp value (start) cannot be greater than the maximum timestamp value (end)`,
		},
		{
			name:  "fail_progress_enabled_but_no_read_write_storage_url_provided",
			input: []string{"-start='2001-09-09T01:46:41+00:00'", "-end='2001-09-09T01:46:40+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'2001-09-09T01:46:41+00:00'",
				end:           "'2001-09-09T01:46:40+00:00'",
				humanReadable: true,
				mint:          1000000001000,
				mintSec:       1000000001,
				maxt:          1000000000000,
				maxtSec:       1000000000,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				progressMetricURL:  "",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: true,
			errMessage:      `invalid input: minimum timestamp value (start) cannot be greater than the maximum timestamp value (end)`,
		},
		{
			name:  "pass_progress_enabled_and_read_write_storage_url_provided",
			input: []string{"-start='2001-09-09T01:46:41+00:00'", "-end='2001-09-09T01:46:42+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-metric-url=http://localhost:9201/read"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'2001-09-09T01:46:41+00:00'",
				end:           "'2001-09-09T01:46:42+00:00'",
				humanReadable: true,
				mint:          1000000001000,
				mintSec:       1000000001,
				maxt:          1000000002000,
				maxtSec:       1000000002,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressMetricURL:  "http://localhost:9201/read",
				progressEnabled:    true,
				concurrentPull:     1,
				maxSlabSizeBytes:   524288000,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
			},
			failsValidation: false,
			errMessage:      `invalid input: minimum timestamp value (mint) cannot be greater than the maximum timestamp value (maxt)`,
		},
		// Mutual exclusive tests.
		{
			name: "pass_normal_exclusive_password",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-read-auth-password=password",
				"-la-increment=7m", "-max-read-duration=7h"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				maxReadDuration:    time.Hour * 7,
				laIncrement:        time.Minute * 7,
				progressMetricURL:  "",
				maxSlabSizeBytes:   524288000,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
				progressEnabled:    false,
				concurrentPull:     1,
				readerAuth:         utils.Auth{Password: "password"},
			},
			failsValidation: false,
		},
		{
			name:  "pass_normal_exclusive_bearer_token",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-read-auth-bearer-token=token"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressMetricURL:  "",
				maxSlabSizeBytes:   524288000,
				maxSlabSize:        "500MB",
				concurrentPush:     1,
				progressEnabled:    false,
				concurrentPull:     1,
				readerAuth:         utils.Auth{BearerToken: "token"},
			},
			failsValidation: false,
		},
		{
			name:  "fail_non_exclusive_bearer_token_and_password",
			input: []string{"-start='1970-01-01T00:16:40+00:00'", "-end='1970-01-01T00:16:41+00:00'", "-reader-url=http://localhost:9090/api/v1/read", "-writer-url=http://localhost:9201/write", "-progress-enabled=false", "-read-auth-password=password", "-read-auth-bearer-token=token"},
			expectedConf: &config{
				name:          "prom-migrator",
				start:         "'1970-01-01T00:16:40+00:00'",
				end:           "'1970-01-01T00:16:41+00:00'",
				humanReadable: true,
				mint:          1000000,
				mintSec:       1000,
				maxt:          1001000,
				maxtSec:       1001,
				readerClient: utils.ClientRuntime{
					URL:          "http://localhost:9090/api/v1/read",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				writerClient: utils.ClientRuntime{
					URL:          "http://localhost:9201/write",
					Timeout:      defaultTimeout,
					Delay:        defaultRetryDelay,
					OnTimeoutStr: "retry",
					OnErrStr:     "abort",
				},
				progressMetricName: "prom_migrator_progress",
				maxReadDuration:    defaultMaxReadDuration,
				laIncrement:        defaultLaIncrement,
				progressMetricURL:  "",
				maxSlabSize:        "500MB",
				concurrentPush:     1,
				progressEnabled:    false,
				concurrentPull:     1,
				readerAuth:         utils.Auth{Password: "password", BearerToken: "token"},
			},
			failsValidation: true,
			errMessage:      `reader auth validation: at most one of basic_auth, oauth2, bearer_token & bearer_token_file must be configured`,
		},
	}

	for _, c := range cases {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		config := new(config)
		parseFlags(config, c.input)

		err := validateConf(config)
		if c.failsValidation {
			if err == nil {
				t.Fatalf(fmt.Sprintf("%s should have failed", c.name))
			}
			assert.Equal(t, c.errMessage, err.Error(), fmt.Sprintf("validation: %s", c.name))
		} else {
			assert.NoError(t, err, fmt.Sprintf("parsing input into config: %s", c.name))
		}

		err = parseClientInfo(c.expectedConf)
		assert.NoError(t, err, fmt.Sprintf("parsing expected flags: %s", c.name))

		assert.Equal(t, c.expectedConf, config, fmt.Sprintf("parse-flags: %s", c.name))
		if err != nil && !c.failsValidation {
			assert.NoError(t, err, c.name)
		}
	}
}
