package game

type Equipment struct {
	Slots map[string]*ItemInstance `yaml:"slots,omitempty"`
}

func NewEquipment() Equipment {
	return Equipment{
		Slots: make(map[string]*ItemInstance),
	}
}

func (e *Equipment) GetSlots() map[string]*ItemInstance {
	return e.Slots
}

func (e *Equipment) Equip(slot string, item *ItemInstance) {
	e.Slots[slot] = item
}

func (e *Equipment) Unequip(slot string) *ItemInstance {
	item := e.Slots[slot]
	delete(e.Slots, slot)

	return item
}

func (e *Equipment) GetItem(slot string) *ItemInstance {
	return e.Slots[slot]
}
