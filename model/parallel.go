package model

// Parallel describes parallel dependencies.
type Parallel struct {
	Seq     Sequence    `json:",omitempty"`
	Service *ServiceDep `json:",omitempty"`
	MaxPar  int         `json:",omitempty"`
}

// Validate performs validation and sets defaults.
func (p Parallel) Validate(r *Registry) error {
	return nil
}
