package dumb_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/bodychoppers/dumb"

	"testing"
)

func TestDumb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dumb Suite")
}

var _ = Describe("Dumb", func() {
	It("chops", func() {
		body := `Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse.
Die äußerste Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht.
Zwischen der Bastschicht und dem Holz befindet sich bei Gymnospermen.
Hinsichtlich des inneren Baus des Baumstamms weichen die zu den Einkeimblättrigen ab.
`
		c := dumb.BodyChopper{80}
		position := 0
		positionWithinBodyPart := 0

		Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse.
Die äußerst`))

		position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
		Expect(position).To(Equal(0))
		Expect(positionWithinBodyPart).To(Equal(80))

		Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`e Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht.
Zwischen der Ba`))
		position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
		Expect(position).To(Equal(0))
		Expect(positionWithinBodyPart).To(Equal(160))

		Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`stschicht und dem Holz befindet sich bei Gymnospermen.
Hinsichtlich des inneren `))
		position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
		Expect(position).To(Equal(0))
		Expect(positionWithinBodyPart).To(Equal(240))

		Expect(c.FetchBodyPart(body, positionWithinBodyPart)).To(Equal(`Baus des Baumstamms weichen die zu den Einkeimblättrigen ab.
`))
		position, positionWithinBodyPart = c.MoveToNextBodyPart(body, position, positionWithinBodyPart)
		Expect(position).To(Equal(1))
		Expect(positionWithinBodyPart).To(Equal(0))
	})
})
