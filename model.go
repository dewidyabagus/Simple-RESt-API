package main

import "time"

type Account struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);not null" json:"email"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

/*
	Data Seeder :
	INSERT INTO `accounts`(`name`, `email`, `created_at`)
	VALUES ('John Wick', 'john.wick@email.com', NOW()), ('Dick Mathin', 'dick.marthin@email.com', NOW()), ('Herro Indo', 'herro@email.com', NOW());
*/
