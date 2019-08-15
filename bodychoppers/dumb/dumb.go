package dumb

import "github.com/pkg/math"

type BodyChopper struct {
	MaxBodyPartLen int
}

func (c BodyChopper) MoveToNextBodyPart(body string, currentPosition int, currentPositionWithinSectionBody int) (newPosition int, newPositionWithinSectionBody int) {
	if currentPositionWithinSectionBody+c.MaxBodyPartLen >= len(body) {
		return currentPosition + 1, 0
	}
	return currentPosition, currentPositionWithinSectionBody + c.MaxBodyPartLen
}

func (c BodyChopper) FetchBodyPart(body string, currentPositionWithinSectionBody int) string {
	return body[currentPositionWithinSectionBody:math.Min(len(body), currentPositionWithinSectionBody+c.MaxBodyPartLen)]
}
