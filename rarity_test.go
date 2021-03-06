package opensea_go

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRarityScore(t *testing.T) {
	tests := [][]Meta{
		// one NFT series
		{
			{
				TokenId: 1,
				Attributes: []Trait{
					{TraitType: "A", Value: "a1"},
					{TraitType: "B", Value: "b1"},
				},
			},
			{
				TokenId: 2,
				Attributes: []Trait{
					{TraitType: "A", Value: "a2"},
					{TraitType: "B", Value: "b1"},
				},
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
