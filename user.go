package main

import "mini-gambler/money"

type User struct {
	ID       int64        `json:"-"`
	Username string       `json:"username"`
	Balance  *money.Money `json:"balance"`
}

func (u *User) Debit(amount money.Money) User {
	u.Balance.Debit(amount)
	return *u
}

func (u *User) Credit(amount money.Money) User {
	u.Balance.Credit(amount)
	return *u
}
