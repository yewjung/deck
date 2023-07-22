package main

import (
	"errors"
	"math/rand"
	"time"
)

func numberOf1Bits(deck uint64) uint64 {
	deck = (deck & 0x5555555555555555) + ((deck >> 1) & 0x5555555555555555)
	deck = (deck & 0x3333333333333333) + ((deck >> 2) & 0x3333333333333333)
	deck = (deck & 0x0f0f0f0f0f0f0f0f) + ((deck >> 4) & 0x0f0f0f0f0f0f0f0f)
	deck = (deck & 0x00ff00ff00ff00ff) + ((deck >> 8) & 0x00ff00ff00ff00ff)
	deck = (deck & 0x0000ffff0000ffff) + ((deck >> 16) & 0x0000ffff0000ffff)
	return (deck & 0x00000000ffffffff) + ((deck >> 32) & 0x00000000ffffffff)
}

func pickRandom(noOfBits uint64) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(int(noOfBits))
}

// return new deck after flipping bit
func setBitToZero(deck uint64, bitRank uint64) (uint64, uint64) {
	bitPos := posOfBitRank(deck, bitRank)
	posFromLeft := 64 - bitPos
	deck ^= 1 << posFromLeft
	return deck, posFromLeft
}

func posOfBitRank(v uint64, r uint64) uint64 {
	var a, b, c, d, t, s uint64
	a = (v & 0x5555555555555555) + ((v >> 1) & 0x5555555555555555)
	b = (a & 0x3333333333333333) + ((a >> 2) & 0x3333333333333333)
	c = (b & 0x0f0f0f0f0f0f0f0f) + ((b >> 4) & 0x0f0f0f0f0f0f0f0f)
	d = (c & 0x00ff00ff00ff00ff) + ((c >> 8) & 0x00ff00ff00ff00ff)
	t = d>>32 + d>>48

	// Now do branchless select!
	s = 64

	s -= ((t - r) & 256) >> 3
	r -= (t & ((t - r) >> 8))
	t = (d >> (s - 16)) & 0xff

	s -= ((t - r) & 256) >> 4
	r -= (t & ((t - r) >> 8))
	t = (c >> (s - 8)) & 0xf

	s -= ((t - r) & 256) >> 5
	r -= (t & ((t - r) >> 8))
	t = (b >> (s - 4)) & 0x7

	s -= ((t - r) & 256) >> 6
	r -= (t & ((t - r) >> 8))
	t = (a >> (s - 2)) & 0x3

	s -= ((t - r) & 256) >> 7
	r -= (t & ((t - r) >> 8))
	t = (v >> (s - 1)) & 0x1

	s -= ((t - r) & 256) >> 8
	return 65 - s
}

func DrawCard(deck uint64) (uint64, Card, error) {
	if deck == 0 {
		Logger.Error(DECK_IS_EMPTY)
		return 0, Card{}, errors.New(DECK_IS_EMPTY)
	}
	noOfBits := numberOf1Bits(deck)
	// 1 <= noOfBits <= 52
	bitToFlip := pickRandom(noOfBits)
	// 0 <= bitToFlip <= 51
	bitToFlip64 := uint64(bitToFlip)
	deck, bitPosFromLeft := setBitToZero(deck, bitToFlip64+1)
	// 0 <= butPosFromLeft <= 51
	return deck, convertBitToCard(bitPosFromLeft), nil
}

var suits = [4]string{"H", "S", "C", "D"}
var ranks = [13]string{
	"A",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"10",
	"J",
	"Q",
	"K",
}

func convertBitToCard(bitPos uint64) Card {
	rank := bitPos % 13
	suit := bitPos / 13
	return Card{Rank: ranks[rank], Suit: suits[suit]}
}

type Card struct {
	Rank string `json:"rank"`
	Suit string `json:"suit"`
}
