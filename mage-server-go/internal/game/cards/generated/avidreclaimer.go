package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Avid Reclaimer", NewAvidReclaimer)
}

// NewAvidReclaimer creates a Avid Reclaimer
// {2}{G} - CREATURE
func NewAvidReclaimer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Avid Reclaimer")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DRUID"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	return card, nil
}
