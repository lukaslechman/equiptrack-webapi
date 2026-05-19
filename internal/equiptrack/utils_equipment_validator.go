package equiptrack

import (
	"fmt"
	"time"
)

func validateEquipment(equipment Equipment) (string, bool) {
	if equipment.Name == "" {
		return "Názov vybavenia je povinný", false
	}
	if equipment.SerialNumber == "" {
		return "Sériové číslo je povinné", false
	}
	if equipment.Manufacturer == "" {
		return "Výrobca je povinný", false
	}
	if equipment.PurchasePrice < 0 {
		return "Obstarávacia cena nemôže byť záporná", false
	}
	if equipment.LifespanYears < 1 || equipment.LifespanYears > 50 {
		return "Životnosť musí byť 1-50 rokov", false
	}

	// validácia formátu dátumu nákupu
	var purchaseDate time.Time
	if equipment.PurchaseDate == "" {
		return "Dátum nákupu je povinný", false
	}
	purchaseDate, err := time.Parse("2006-01-02", equipment.PurchaseDate)
	if err != nil {
		return fmt.Sprintf("Neplatný formát dátumu nákupu '%s', očakávaný formát YYYY-MM-DD", equipment.PurchaseDate), false
	}

	// dátum nákupu nesmie byť v budúcnosti
	if purchaseDate.After(time.Now()) {
		return "Dátum nákupu nemôže byť v budúcnosti", false
	}

	// validácia záruky ak je zadaná
	if equipment.WarrantyUntil != "" {
		warrantyDate, err := time.Parse("2006-01-02", equipment.WarrantyUntil)
		if err != nil {
			return fmt.Sprintf("Neplatný formát dátumu záruky '%s', očakávaný formát YYYY-MM-DD", equipment.WarrantyUntil), false
		}
		// záruka musí byť po dátume nákupu
		if warrantyDate.Before(purchaseDate) {
			return "Dátum záruky nemôže byť pred dátumom nákupu", false
		}
	}

	return "", true
}
