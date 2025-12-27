package curve

import "github.com/mjbagaoisan/humanstocksbot/internal/domain"

func Price(basePrice, slope, supply int64) int64 {
	return basePrice + slope*supply
}

func BuyCost(basePrice, slope, supply, qty int64) (int64, error) {

	if qty <= 0 {
		return 0, domain.ErrInvalidQuantity
	}

	// 1. Base Cost
	baseCost, err := safeMul(basePrice, qty)
	if err != nil {
		return 0, err
	}

	// 2. Linear Term
	linearTerm, err := safeMul(qty, supply)
	if err != nil {
		return 0, err
	}

	// 3. Quadratic Term: (qty * (qty - 1)) / 2
	qtyMinusOne, err := safeSub(qty, 1)
	if err != nil {
		return 0, err
	}
	quadraticNumerator, err := safeMul(qty, qtyMinusOne)
	if err != nil {
		return 0, err
	}
	quadraticTerm := quadraticNumerator / 2

	// 4. Calculate Slope Impact
	slopeFactor, err := safeAdd(linearTerm, quadraticTerm)
	if err != nil {
		return 0, err
	}
	slopeCost, err := safeMul(slope, slopeFactor)
	if err != nil {
		return 0, err
	}

	totalCost, err := safeAdd(baseCost, slopeCost)
	if err != nil {
		return 0, err
	}

	return totalCost, nil
}

func SellPayout(basePrice, slope, supply, qty int64) (int64, error) {

	if qty == 0 {
		return 0, nil
	}

	if supply < qty {
		return 0, domain.ErrInsufficientSupply
	}

	// selling brings supply down, so payout equals what it cost to buy those shares
	newSupply := supply - qty

	return BuyCost(basePrice, slope, newSupply, qty)
}
