package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Raffine Scheming Seer", NewRaffineSchemingSeer)
}

// NewRaffineSchemingSeer creates a Raffine Scheming Seer
// {W}{U}{B} - CREATURE
// Flying
func NewRaffineSchemingSeer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Raffine Scheming Seer")
	card.ManaCost = "{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPHINX", "DEMON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}