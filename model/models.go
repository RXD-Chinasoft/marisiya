package model
import(
	"fmt"
)
type Friend struct {
	Id int64 `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Friend int64 `json:"friend" db:"friend"`
}

func (f Friend) String() string {
	return fmt.Sprintf("friend is => id : %d, email : %s, Friend : %d \n", f.Id, f.Email, f.Friend)
}