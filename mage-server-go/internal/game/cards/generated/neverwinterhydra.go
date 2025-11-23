package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Neverwinter Hydra", NewNeverwinterHydra)
}

// NewNeverwinterHydra creates a Neverwinter Hydra
// {X}{X}{G}{G} - CREATURE
// Trample
func NewNeverwinterHydra(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Neverwinter Hydra")
	card.ManaCost = "{X}{X}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HYDRA"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
