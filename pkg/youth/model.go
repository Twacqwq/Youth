package pkg

type Member struct {
	MemberId int
	Status   bool `json:"-"`
}
