package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizePrometheusMetricName(t *testing.T) {
	cases := map[string]string{
		"test-1":                           "test_1",
		"testTest2":                        "test_test2",
		"test3":                            "test3",
		"test4-":                           "test4",
		"???test5????":                     "test5",
		"???test6?_?":                      "test6",
		"epa_chamo":                        "epa_chamo",
		"__epa_____funka_":                 "epa_funka",
		"_1_2__3__4_5":                     "1_2_3_4_5",
		"_1_2__$_%_%__%3__4_5":             "1_2_3_4_5",
		"_1_2__$_%_piero_testea%__%3__4_5": "1_2_piero_testea_3_4_5",
		"-otro-test-test-":                 "otro_test_test",
		"∞_tests-":                         "tests",
		"-∞-_tests-----":                   "tests",
		"∞-_tests----t-":                   "tests_t",
		"--_tests----t-%$&//(()))=extra":   "tests_t_extra",
	}
	for test, expected := range cases {
		assert.Equal(t, expected, sanitizeMetricName(test))
	}
}
