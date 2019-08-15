package paragraph

import "github.com/petergtz/alexa-wikipedia/wiki"

type BodyChopper struct {
	MaxBodyPartLen int
	Fallback       wiki.BodyChopper
}

func (c BodyChopper) FetchBodyPart(body string, currentPositionWithinSectionBody int) string {
	result := ""
	fitsIn := true
	for index, runeValue := range body[currentPositionWithinSectionBody:] {
		if index > c.MaxBodyPartLen {
			if result == "" {
				return c.Fallback.FetchBodyPart(body, currentPositionWithinSectionBody)
			}
			fitsIn = false
			break
		}
		if runeValue == '\n' {
			result = body[currentPositionWithinSectionBody : currentPositionWithinSectionBody+index+1]
		}
	}
	if fitsIn {
		result = body[currentPositionWithinSectionBody:]
	}
	return result
}

func (c BodyChopper) MoveToNextBodyPart(body string, currentPosition int, currentPositionWithinSectionBody int) (newPosition int, newPositionWithinSectionBody int) {
	fitsIn := true
	for index, runeValue := range body[currentPositionWithinSectionBody:] {
		if index > c.MaxBodyPartLen {
			if newPositionWithinSectionBody == 0 {
				return c.Fallback.MoveToNextBodyPart(body, currentPosition, currentPositionWithinSectionBody)
			}
			fitsIn = false
			break
		}
		if runeValue == '\n' {
			newPositionWithinSectionBody = currentPositionWithinSectionBody + index + 1
		}
	}
	if fitsIn {
		return currentPosition + 1, 0
	}
	return currentPosition, newPositionWithinSectionBody
}
