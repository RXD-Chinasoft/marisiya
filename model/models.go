package model
import(
	"fmt"
)
type Friend struct {
	Id int64 `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Friends []int64 `json:"friends" db:"friends"`
}

func (f Friend) String() string {
	return fmt.Sprintf("friend is => id : %d, email : %s, Friend : %v \n", f.Id, f.Email, f.Friends)
}