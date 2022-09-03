package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompareTim(t *testing.T) {
	time1 := time.Now()
	time2 := time.Now()

	isEqual := CheckDate(time1, time2)
	require.True(t, isEqual)
}
