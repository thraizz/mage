package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ezekiel Sims Spider Totem", NewEzekielSimsSpiderTotem)
}

// NewEzekielSimsSpiderTotem creates a Ezekiel Sims Spider Totem
// {4}{G} - CREATURE
// Reach
func NewEzekielSimsSpiderTotem(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ezekiel Sims Spider Totem")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER", "HUMAN", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
