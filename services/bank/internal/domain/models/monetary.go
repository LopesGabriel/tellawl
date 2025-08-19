package models

type Monetary struct {
	Value  int
	Offset int
}

func (m Monetary) Sum(amount Monetary) Monetary {
	normalizedValue1 := m.Value / m.Offset
	normalizedValue2 := amount.Value / amount.Offset

	return Monetary{
		Value:  m.Offset * (normalizedValue1 + normalizedValue2),
		Offset: m.Offset,
	}
}

func (m Monetary) Sub(amount Monetary) Monetary {
	normalizedValue1 := m.Value / m.Offset
	normalizedValue2 := amount.Value / amount.Offset

	return Monetary{
		Value:  m.Offset * (normalizedValue1 - normalizedValue2),
		Offset: m.Offset,
	}
}
