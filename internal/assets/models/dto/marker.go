package dto

type Marker interface {
	isDTO()
}

func (Bond) isDTO()     {}
func (Share) isDTO()    {}
func (Etf) isDTO()      {}
func (Currency) isDTO() {}
