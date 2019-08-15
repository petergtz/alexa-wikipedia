package paragraph_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/bodychoppers/dumb"
	"github.com/petergtz/alexa-wikipedia/bodychoppers/paragraph"

	"testing"
)

func TestParagraph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Paragraph Suite")
}

const body = `Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse.
Die äußerste Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht.
Zwischen der Bastschicht und dem Holz befindet sich bei Gymnospermen.
Hinsichtlich des inneren Baus des Baumstamms weichen die zu den Einkeimblättrigen ab.
`

var _ = Describe("Paragraph", func() {
	Context("MaxBodyPart 100", func() {
		It("chops into 4 body parts", func() {
			c := paragraph.BodyChopper{
				MaxBodyPartLen: 100,
				Fallback: &dumb.BodyChopper{
					MaxBodyPartLen: 100,
				},
			}
			position := 0
			positionWithinBodyPart := 0

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse.` + "\n"))

			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(67))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Die äußerste Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht.` + "\n"))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(145))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Zwischen der Bastschicht und dem Holz befindet sich bei Gymnospermen.` + "\n"))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(215))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Hinsichtlich des inneren Baus des Baumstamms weichen die zu den Einkeimblättrigen ab.` + "\n"))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(1))
			Expect(positionWithinBodyPart).To(Equal(0))
		})
	})

	Context("MaxBodyPart 80", func() {
		It("chops into 5 body parts using the dumb BodyChopper as fallback", func() {
			c := paragraph.BodyChopper{
				MaxBodyPartLen: 80,
				Fallback: &dumb.BodyChopper{
					MaxBodyPartLen: 80,
				},
			}
			position := 0
			positionWithinBodyPart := 0

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse.` + "\n"))

			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(67))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Die äußerste Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht.` + "\n"))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(145))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Zwischen der Bastschicht und dem Holz befindet sich bei Gymnospermen.` + "\n"))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(215))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Hinsichtlich des inneren Baus des Baumstamms weichen die zu den Einkeimblättrig`))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(0))
			Expect(positionWithinBodyPart).To(Equal(295))

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`en ab.` + "\n"))
			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(1))
			Expect(positionWithinBodyPart).To(Equal(0))
		})
	})

	Context("MaxBodyPart 500", func() {
		It("chops into 1 body part", func() {

			body := `Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse.
Die äußerste Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht.
Zwischen der Bastschicht und dem Holz befindet sich bei Gymnospermen.
Hinsichtlich des inneren Baus des Baumstamms weichen die zu den Einkeimblättrigen ab.`

			c := paragraph.BodyChopper{
				MaxBodyPartLen: 500,
				Fallback: &dumb.BodyChopper{
					MaxBodyPartLen: 500,
				},
			}
			position := 0
			positionWithinBodyPart := 0

			Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(body))

			position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
			Expect(position).To(Equal(1))
			Expect(positionWithinBodyPart).To(Equal(0))
		})
	})

})
