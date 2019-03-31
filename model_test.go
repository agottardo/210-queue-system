package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBasicFunctionality(t *testing.T) {
	config = ReadConfig()
	require.Zero(t, len(queue.Entries))
	require.Zero(t, len(UnservedEntries()))
	JoinQueue("Joe Student", "r3a1b", "Totally lost.")
	JoinQueue("Diligent Student", "r3a2b", "Totally lost, again.")
	require.Equal(t, 2, len(queue.Entries))
	require.Equal(t, 2, len(UnservedEntries()))
	require.Equal(t, "r3a1b", queue.Entries[0].CSid)
	require.Equal(t, "Joe Student", queue.Entries[0].Name)
	require.Equal(t, "r3a2b", queue.Entries[1].CSid)
	require.Equal(t, "Diligent Student", queue.Entries[1].Name)
	require.True(t, HasJoinedQueue("r3a1b"))
	require.True(t, HasJoinedQueue("r3a2b"))
	require.False(t, HasJoinedQueue("r3a3b"))
	require.Zero(t, NumTimesHelped("r3a1b"))
	require.Zero(t, NumTimesHelped("r3a2b"))
	require.Zero(t, NumTimesHelped("r3a3b"))
	ServeStudent("r3a3b") // Does nothing
	require.Equal(t, 2, len(queue.Entries))
	require.Equal(t, 2, len(UnservedEntries()))
	ServeStudent("r3a2b")
	require.Equal(t, 2, len(queue.Entries)) // Data remains in memory
	require.Equal(t, 1, len(UnservedEntries()))
	require.Equal(t, uint(1), NumTimesHelped("r3a2b"))
	require.Zero(t, NumTimesHelped("r3a1b"))
	exists1, position1 := QueuePositionForCSID("r3a1b")
	require.True(t, exists1)
	require.Zero(t, position1)
	exists2, position2 := QueuePositionForCSID("r3a2b")
	require.False(t, exists2)
	require.Zero(t, position2)
	require.Equal(t, uint(2), TotalNumStudentsHelped())
}
