package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jugan The Rising Star", NewJuganTheRisingStar)
}

// NewJuganTheRisingStar creates a Jugan The Rising Star
// {3}{G}{G}{G} - CREATURE
// Flying
func NewJuganTheRisingStar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jugan The Rising Star")
	card.ManaCost = "{3}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
