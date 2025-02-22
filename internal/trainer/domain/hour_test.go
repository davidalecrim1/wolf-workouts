package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewAvailableHour(t *testing.T) {
	t.Run("schedule training in valid hour", func(t *testing.T) {
		t.Parallel()

		timestamp := time.Now().Truncate(time.Hour).AddDate(0, 0, 2)
		h, err := NewAvailableHour(timestamp)

		require.NoError(t, err)
		require.NotNil(t, h)
	})

	t.Run("do not allow to schedule training in the past", func(t *testing.T) {
		t.Parallel()

		timestamp := time.Now().Truncate(time.Hour).AddDate(0, 0, -1)
		h, err := NewAvailableHour(timestamp)

		require.ErrorIs(t, err, ErrPastHour)
		require.Nil(t, h)
	})

	t.Run("do not allow to schedule too into the future", func(t *testing.T) {
		t.Parallel()

		daysIntoTheFuture := 8
		timestamp := time.Now().Truncate(time.Hour).AddDate(0, 0, daysIntoTheFuture)
		h, err := NewAvailableHour(timestamp)

		require.ErrorIs(t, err, ErrTooDistantDate)
		require.Nil(t, h)
	})

	t.Run("make sure the training schedule hour is full hour", func(t *testing.T) {
		t.Parallel()

		timestamp := time.Now()
		h, err := NewAvailableHour(timestamp)

		require.ErrorIs(t, err, ErrNotFullHour)
		require.Nil(t, h)
	})
}
