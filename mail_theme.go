package main

import "github.com/matcornic/hermes"

type MailTheme interface {
	hermes.Theme
	GetStyle() string
}
