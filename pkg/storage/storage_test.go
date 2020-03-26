package storage_test

import (
	"testing"

	"github.com/convox/console/pkg/storage"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require.IsType(t, &storage.Dynamo{}, storage.New("dynamo"))
	require.IsType(t, &storage.Dynamo{}, storage.New(""))
}
