package opensea_go

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRarityScore(t *testing.T) {
	tests := [][]Meta{
		{
			{
				{TraitType: "A", Value: "a1"},
				{TraitType: "B", Value: "b1"},
			},
			{
				{TraitType: "A", Value: "a2"},
				{TraitType: "B", Value: "b1"},
			},
		},
	}
	for _, tt := range tests {
		scores := RarityScore(tt)
		b, err := json.MarshalIndent(scores, "", "  ")
		require.NoError(t, err)
		t.Logf(string(b))
	}
}
