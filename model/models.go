package model
import(
	"fmt"
)
type Friend struct {
	Id int64 `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Friends []int64 `json:"friends" db:"friends"`
	SubscribMgr  []string `json:"subscribMgr " db:"subscribMgr "`
}

func (f Friend) String() string {
	return fmt.Sprintf("friend is => id : %d, email : %s, Friends : %v \n", f.Id, f.Email, f.Friends)
}