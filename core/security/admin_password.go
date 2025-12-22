package security

type AdminPassword interface {
	Compare(string) bool
}

type adminPassword struct {
	value string
}

func (p *adminPassword) Compare(other string) bool {
	return p.value == other
}

func NewAdminPassword(value string) AdminPassword {
	return &adminPassword{value: value}
}
