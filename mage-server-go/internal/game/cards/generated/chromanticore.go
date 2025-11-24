package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chromanticore", NewChromanticore)
}

// NewChromanticore creates a Chromanticore
// {W}{U}{B}{R}{G} - ENCHANTMENT CREATURE
// Flying, FirstStrike, Vigilance, Trample, Lifelink
func NewChromanticore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chromanticore")
	card.ManaCost = "{W}{U}{B}{R}{G}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"MANTICORE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability2)
	ability3 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability3)
	ability4 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability4)
	ability5, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FirstStrikeAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("VigilanceAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("TrampleAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("LifelinkAbility", abilities.AttachmentTypeAura)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability5)
	ability6, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEnchantedEffect(4, 4)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability6)
	return card, nil
}
