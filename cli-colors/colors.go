package clicolors

import (
	"fmt"
	"math/rand"
	"time"
)

type ColorMode string

const (
	HTML     ColorMode = "html"
	Terminal ColorMode = "terminal"
)

type ColorName string

const (
	LightGreen ColorName = "LightGreen"
	LightBlue  ColorName = "LightBlue"
	LightRed   ColorName = "LightRed"
	Red        ColorName = "Red"
	White      ColorName = "White"
	Reset      ColorName = "Reset"
	Random     ColorName = "Random"
)

var Colors = map[ColorMode]map[ColorName]string{
	HTML: {
		LightGreen: "<span color='#A9DFBF'>",
		LightBlue:  "<span color='#A9CCE3'>",
		LightRed:   "<span color='#F5B7B1'>",
		Red:        "<span color='#E74C3C'>",
		White:      "<span color='#FFFFFF'>",
		Reset:      "</span>",
	},
	Terminal: {
		LightGreen: "\033[38;2;169;223;191m",
		LightBlue:  "\033[38;2;169;204;227m",
		LightRed:   "\033[38;2;245;183;177m",
		Red:        "\033[38;2;231;76;60m",
		White:      "\033[38;2;255;255;255m",
		Reset:      "\033[0m",
	},
}

func GetMinutesDifference(t time.Time) float64 {
	now := time.Now()
	diff := t.Sub(now).Abs()
	return diff.Minutes()
}

func GetUnreadColor(colorMode ColorMode, c int) string {
	if c > 15 {
		return GetColor(colorMode, Red)
	} else if c > 5 {
		return GetColor(colorMode, LightBlue)
	} else if c > 0 {
		return GetColor(colorMode, LightGreen)
	} else {
		return GetColor(colorMode, White)
	}
}

func GetEventColor(colorMode ColorMode, t time.Time) string {
	minutes := GetMinutesDifference(t)

	if minutes <= 5 {
		return GetColor(colorMode, Red)
	} else if minutes < 15 {
		return GetColor(colorMode, LightRed)
	} else if minutes < 30 {
		return GetColor(colorMode, LightBlue)
	} else if minutes < 60 {
		return GetColor(colorMode, LightGreen)
	} else {
		return GetColor(colorMode, White)
	}
}

func GetColor(mode ColorMode, colorName ColorName) string {
	if colorMode, exists := Colors[mode]; exists {
		if color, exists := colorMode[colorName]; exists {
			return color
		} else {
			return randomColor(mode)
		}
	}
	return ""
}

func randomColor(colorMode ColorMode) string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)

	if colorMode == HTML {
		return fmt.Sprintf("<span style='color: #%02X%02X%02X'>", r, g, b)
	} else {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
}
