package opensea_go

import (
	"fmt"
	"sort"
)

type Trait struct {
	Type    string  `json:"type,omitempty" bson:"type,omitempty"`
	Value   string  `json:"value" bson:"value"`
	Percent float64 `json:"percent,omitempty" bson:"percent,omitempty"`
	Score   float64 `json:"score,omitempty" bson:"score,omitempty"`
}

type Meta struct {
	TokenId    int     `json:"tokenId" bson:"tokenId"`
	Attributes []Trait `json:"attributes" bson:"attributes"`
}

type Rarity struct {
	TokenId    int     `json:"tokenId" bson:"tokenId"`
	Owner      string  `json:"owner,omitempty" bson:"owner,omitempty"`
	Rank       int     `json:"rank,omitempty" bson:"rank,omitempty"`
	Score      float64 `json:"score,omitempty" bson:"score,omitempty"`
	Plus       int     `json:"plus,omitempty" bson:"plus,omitempty"`
	Attributes []Trait `json:"attributes" bson:"attributes"`
}

const (
	//columnsDing     = 4
	columnsDiscord  = 4
	columnsTelegram = 4
)

func (r Rarity) FormatDing() string {
	content := fmt.Sprintf(`
稀有度排名: %d
稀有度得分: %2.f`,
		r.Rank, r.Score,
	)
	for _, trait := range r.Attributes {
		content += "\n"
		if trait.Type != "" {
			content += fmt.Sprintf("%s", trait.Type)
		}
		if trait.Type != "" && (trait.Value != "" || trait.Score != 0) {
			content += ": "
		}
		if trait.Value != "" {
			content += fmt.Sprintf("%s", trait.Value)
		}
		if trait.Score != 0 {
			content += fmt.Sprintf(", %.2f", trait.Score)
		}
	}
	return content
}

func (r Rarity) FormatDiscord() string {
	content := fmt.Sprintf(`
稀有度排名: **%d**
稀有度得分: %2.f`,
		r.Rank, r.Score,
	)
	for i, trait := range r.Attributes {
		if i%columnsDiscord == 0 {
			content += "\n"
		} else {
			content += "\t"
		}
		content += "\n"
		if trait.Type != "" {
			content += fmt.Sprintf("%s", trait.Type)
		}
		if trait.Type != "" && (trait.Value != "" || trait.Score != 0) {
			content += ": "
		}
		if trait.Value != "" {
			content += fmt.Sprintf("%s", trait.Value)
		}
		if trait.Score != 0 {
			content += fmt.Sprintf(", %.2f", trait.Score)
		}
	}
	return content
}

func (r Rarity) FormatTelegram() string {
	return r.FormatDing()
}

// RarityScore calculate the rarity score like rarity.tools
// See https://raritytools.medium.com/ranking-rarity-understanding-rarity-calculation-methods-86ceaeb9b98c
//   [Rarity Score for a Trait Value] = 1 / ([Number of Items with that Trait Value] / [Total Number of Items in Collection])
func RarityScore(metas []Meta) []Rarity {
	sumOfType := make(map[string]int)
	sumOfTraits := make(map[string]map[string]int)
	for _, m := range metas {
		for _, t := range m.Attributes {
			sumOfType[t.Type] += 1
			if mt, ok := sumOfTraits[t.Type]; ok {
				mt[t.Value] += 1
			} else {
				sumOfTraits[t.Type] = make(map[string]int)
				sumOfTraits[t.Type][t.Value] += 1
			}
		}
	}
	var rarities []Rarity
	for i := 0; i < len(metas); i++ {
		var score float64
		for j := 0; j < len(metas[i].Attributes); j++ {
			traitType := metas[i].Attributes[j].Type
			traitValue := metas[i].Attributes[j].Value
			metas[i].Attributes[j].Percent = float64(sumOfTraits[traitType][traitValue]) / float64(sumOfType[traitType])
			metas[i].Attributes[j].Score = float64(sumOfType[traitType]) / float64(sumOfTraits[traitType][traitValue])
			score += metas[i].Attributes[j].Score
		}
		rarities = append(rarities, Rarity{
			TokenId:    metas[i].TokenId,
			Score:      score,
			Attributes: metas[i].Attributes,
		})
	}
	sort.Sort(sort.Reverse(byScore(rarities)))
	for i := 0; i < len(rarities); i++ {
		rarities[i].Rank = i + 1
	}
	return rarities
}

type byScore []Rarity

func (s byScore) Len() int {
	return len(s)
}
func (s byScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byScore) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}
