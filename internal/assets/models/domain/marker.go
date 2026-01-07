package domain

type Marker interface {
	isDomain()
}

func (Bond) isDomain()     {}
func (Share) isDomain()    {}
func (Etf) isDomain()      {}
func (Currency) isDomain() {}
