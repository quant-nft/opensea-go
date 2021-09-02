package opensea_go

type Trait struct {
	TraitType string  `json:"traitType" bson:"trait_type"`
	Value     string  `json:"value" bson:"value"`
	Percent   float64 `json:"percent" bson:"percent"`
	Score     float64 `json:"score" bson:"score"`
}

type Meta []Trait

// RarityScore calculate the rarity score like rarity.tools
// See https://raritytools.medium.com/ranking-rarity-understanding-rarity-calculation-methods-86ceaeb9b98c
//   [Rarity Score for a Trait Value] = 1 / ([Number of Items with that Trait Value] / [Total Number of Items in Collection])
func RarityScore(metas []Meta) []Meta {
	sumOfType := make(map[string]int)
	sumOfTraits := make(map[string]map[string]int)
	for _, m := range metas {
		for _, t := range m {
			sumOfType[t.TraitType] += 1
			if mt, ok := sumOfTraits[t.TraitType]; ok {
				mt[t.Value] += 1
			} else {
				sumOfTraits[t.TraitType] = make(map[string]int)
				sumOfTraits[t.TraitType][t.Value] += 1
			}
		}
	}
	//total := len(metas)
	for i := 0; i < len(metas); i++ {
		for j := 0; j < len(metas[i]); j++ {
			traitType := metas[i][j].TraitType
			traitValue := metas[i][j].Value
			metas[i][j].Percent = float64(sumOfTraits[traitType][traitValue]) / float64(sumOfType[traitType])
			metas[i][j].Score = float64(sumOfType[traitType]) / float64(sumOfTraits[traitType][traitValue])
		}
	}
	return metas
}
