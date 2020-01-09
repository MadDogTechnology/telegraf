package nats

import (
	"testing"

	"github.com/MadDogTechnology/telegraf/plugins/serializers"
	"github.com/MadDogTechnology/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestConnectAndWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	server := []string{"nats://" + testutil.GetLocalHost() + ":4222"}
	s, _ := serializers.NewInfluxSerializer()
	n := &NATS{
		Servers:    server,
		Subject:    "telegraf",
		serializer: s,
	}

	// Verify that we can connect to the NATS daemon
	err := n.Connect()
	require.NoError(t, err)

	// Verify that we can successfully write data to the NATS daemon
	err = n.Write(testutil.MockMetrics())
	require.NoError(t, err)
}
