package helpers

import "math/rand"

type Game struct {
	Command string;
	Player float64;
	BallX float64;
	BallY float64;
}

type Init struct {
	Command string;
	IsLeft bool;
}

type Ready struct {
	Command string;
}

type Score struct {
	Command string;
	Scored bool;
	LeftP int;
	RightP int;
}

type Ball struct {
	X float64;
	Y float64;
	Dx float64;
	Dy float64;
	Velocity float64;
}

func (b *Ball)Move(delta float64) {
	b.X += b.Dx * delta * b.Velocity
	b.Y += b.Dy * delta * b.Velocity
}

func GetBall() Ball{
	var b Ball
	b.X = 396
	b.Y = 426
	b.Dx = 1
	b.Dy = -1
	b.Velocity = 260

	return b
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString (n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}