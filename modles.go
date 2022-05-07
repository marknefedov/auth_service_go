package main

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionName string    `json:"session_name" bson:"session_name"`
	OpaqueToken string    `json:"opaque_token" bson:"opaque_token"`
	LastUsed    time.Time `json:"last_used" bson:"last_used"`
}

type User struct {
	Uuid     uuid.UUID `json:"uuid" bson:"uuid"`
	FullName string    `json:"full_name" bson:"full_name"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
	Sessions []Session `json:"sessions" bson:"sessions"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" bson:"full_name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type SessionResponse struct {
	SessionUUID string `json:"session_uuid" bson:"session_uuid"`
}

type AccsessTokenRequest struct {
	SessionUUID string `json:"session_uuid" bson:"session_uuid"`
}

type AccsessTokenResponse struct {
	AccsessToken string `json:"accsess_token"`
}
