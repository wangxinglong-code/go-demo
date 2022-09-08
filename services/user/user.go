package user

import "go-demo/models"

func GetUser(uid int) (*models.User, error) {
	info, err := new(models.User).Get(uid)
	return info, err
}
