package model

type UserType string

const (
	UserTypeClient   UserType = "client"
	UserTypeOperator UserType = "operator"
)

func (UserType) Values() []string {
	return []string{
		string(UserTypeClient),
		string(UserTypeOperator),
	}
}
