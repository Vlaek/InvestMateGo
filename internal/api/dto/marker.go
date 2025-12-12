package dto

type DTOMarker interface {
	isDTO()
}

func (BondDTO) isDTO()     {}
func (ShareDTO) isDTO()    {}
func (EtfDTO) isDTO()      {}
func (CurrencyDTO) isDTO() {}
