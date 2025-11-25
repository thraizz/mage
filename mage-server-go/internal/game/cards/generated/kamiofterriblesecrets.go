package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kami Of Terrible Secrets", NewKamiOfTerribleSecrets)
}

// NewKamiOfTerribleSecrets creates a Kami Of Terrible Secrets
// {3}{B} - CREATURE
func NewKamiOfTerribleSecrets(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kami Of Terrible Secrets")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
