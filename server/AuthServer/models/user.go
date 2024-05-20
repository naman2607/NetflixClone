//Data models representing entities such as users

package user

import "time"

type UserBasicDetails struct {
	Id        int       `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
