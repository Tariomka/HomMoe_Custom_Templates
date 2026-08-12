//go:build integration_test

package dialogs

// SelectType picks the bonus type by its dropdown label and reports whether
// such a label exists. ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SelectType(label string) bool {
	return this.typeDropdown.SelectByName(label)
}

// SelectReceiver ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SelectReceiver(name string) bool {
	return this.receiverDropdown.SelectByName(name)
}

// ClickAdd ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) ClickAdd() { this.addBtn.Click() }

// ClickCancel ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) ClickCancel() { this.cancelBtn.Click() }

// ClickPickSpells ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) ClickPickSpells() { this.pickSpellBtn.Click() }

// ClickPickItem ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) ClickPickItem() { this.pickItemBtn.Click() }

// ClickRemoveSpell queues a click on the remove button of the picked spell at
// index and reports whether that row exists. ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) ClickRemoveSpell(index int) bool {
	if index < 0 || index >= len(this.spellRemoveBtns) {
		return false
	}

	this.spellRemoveBtns[index].Click()
	return true
}

// SelectedSpells ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SelectedSpells() []string { return this.selectedSpells }

// ErrorText ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) ErrorText() string { return this.errorText }

// SetMakeFree ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SetMakeFree(value bool) { this.makeFree.Value = value }

// SetMultiplier ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SetMultiplier(text string) { this.multiplierEdit.SetText(text) }

// SetMovement ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SetMovement(text string) { this.movementEdit.SetText(text) }

// SetItem ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SetItem(text string) { this.itemEdit.SetText(text) }

// SetResourceAmount ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SetResourceAmount(text string) { this.resourceEdit.SetText(text) }

// SpellCountLabel exposes the picked-spell caption for the current selection.
// ONLY FOR INTEGRATION TEST USE
func (this *BonusPickerDialog) SpellCountLabel() string {
	return this.handler.GetSpellCountLabel(len(this.selectedSpells))
}
