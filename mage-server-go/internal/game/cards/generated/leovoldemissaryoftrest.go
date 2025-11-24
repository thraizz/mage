package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leovold Emissary Of Trest", NewLeovoldEmissaryOfTrest)
}

// NewLeovoldEmissaryOfTrest creates a Leovold Emissary Of Trest
// {B}{G}{U} - CREATURE
func NewLeovoldEmissaryOfTrest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leovold Emissary Of Trest")
	card.ManaCost = "{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
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
