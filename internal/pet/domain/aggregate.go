package domain

import "fmt"

type (
	PetId  string
	UserId string
)

type Pet struct {
	Id     PetId
	UserId UserId
}

func (p *Pet) String() string {
	return fmt.Sprintf("Pet(id=%s,userId=%s)", p.Id, p.UserId)
}
