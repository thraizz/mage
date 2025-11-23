package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cavalier Of Dawn", NewCavalierOfDawn)
}

// NewCavalierOfDawn creates a Cavalier Of Dawn
// {2}{W}{W}{W} - CREATURE
// Vigilance
func NewCavalierOfDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cavalier Of Dawn")
	card.ManaCost = "{2}{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "KNIGHT"}
	card.Power = "4"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
