package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func (p *GamePerson) putValue(value int, length int, offset int) {
	for i := offset; i < offset+length; i++ {
		bitNumber := i % 8
		byteNumber := i / 8

		j := length - (i - offset) - 1

		currentValue := (value >> j) & 1
		if currentValue > 0 {
			p.mask[byteNumber] |= (1 << (7 - bitNumber))
		}
	}
}

func (p *GamePerson) getValue(offset int, length int) int {
	value := 0
	for i := offset; i < offset+length; i++ {
		byteNumber := i / 8
		bitNumber := i % 8

		j := length - (i - offset) - 1
		currentValue := (p.mask[byteNumber] & (1 << (7 - bitNumber))) >> (7 - bitNumber)
		if currentValue > 0 {
			value += 1 << j
		}
	}
	return value
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(mana, 10, 0)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(health, 10, 10)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(respect, 4, 20)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(strength, 4, 24)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(experience, 4, 28)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(level, 4, 32)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(1, 1, 36)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(1, 1, 37)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(1, 1, 38)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.putValue(personType, 2, 39)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z int32
	gold    uint32
	name    [42]byte
	mask    [6]byte
	/*
		mp  0 - 1000 =>  10 bits
		hp  0 - 1000 =>  10 bits

		respect    0 - 10 => 4 bit
		strength   0 - 10 => 4 bit
		experience 0 - 10 => 4 bit
		level      0 - 10 => 4 bit

		house  0 - 1 => 1 bit
		family 0 - 1 => 1 bit
		gun    0 - 1 => 1 bit
		typ    0 - 2 => 2 bit
	*/
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	length := 0
	for _, b := range p.name {
		if b == 0 {
			break
		}
		length++
	}
	return string(p.name[:length])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return p.getValue(0, 10)
}

func (p *GamePerson) Health() int {
	return p.getValue(10, 10)
}

func (p *GamePerson) Respect() int {
	return p.getValue(20, 4)
}

func (p *GamePerson) Strength() int {
	return p.getValue(24, 4)
}

func (p *GamePerson) Experience() int {
	return p.getValue(28, 4)
}

func (p *GamePerson) Level() int {
	return p.getValue(32, 4)
}

func (p *GamePerson) HasHouse() bool {
	return p.getValue(36, 1) > 0
}

func (p *GamePerson) HasGun() bool {
	return p.getValue(37, 1) > 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.getValue(38, 1) > 0
}

func (p *GamePerson) Type() int {
	return p.getValue(39, 2)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
