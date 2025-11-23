package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Briar Hydra", NewBriarHydra)
}

// NewBriarHydra creates a Briar Hydra
// {5}{G} - CREATURE
// Trample
func NewBriarHydra(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Briar Hydra")
	card.ManaCost = "{5}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PLANT", "HYDRA"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
