package multipliers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/risk"
	"github.com/stackrox/rox/pkg/set"
)

const (
	riskyComponentCountFloor = 0
	riskyComponentCountCeil  = 10
	maxRiskyScore            = 1.5
	maxComponentsInMessage   = 5
)

var riskyComponents = set.NewStringSet(
	"apk",
	"apt",
	"bash",
	"curl",
	"dnf",
	"netcat",
	"nmap",
	"rpm",
	"sh",
	"tcsh",
	"telnet",
	"wget",
	"yum",
)

// riskyComponentCountMultiplier is a scorer for the components in an image that can be used by attackers.
type riskyComponentCountMultiplier struct{}

// NewRiskyComponents provides a multiplier that scores the data based on the the number of risky components in images.
func NewRiskyComponents() Multiplier {
	return &riskyComponentCountMultiplier{}
}

// Score takes a image and evaluates its risk based on image component counts.
func (c *riskyComponentCountMultiplier) Score(_ context.Context, msg proto.Message) *storage.Risk_Result {
	image, ok := msg.(*storage.Image)
	if !ok {
		return nil
	}

	// Get the largest number of risky components in an image
	// Create a name to version map of all the image components.
	presentComponents := set.NewStringSet()
	for _, component := range image.GetScan().GetComponents() {
		presentComponents.Add(component.GetName())
	}

	// Count how many known risky components match a labeled component.
	riskySet := riskyComponents.Intersect(presentComponents)

	if riskySet.Cardinality() == 0 {
		return nil
	}

	// Linear increase between 10 components and 20 components from weight of 1 to 1.5.
	score := float32(1.0) + float32(riskySet.Cardinality()-riskyComponentCountFloor)/float32(riskyComponentCountCeil-riskyComponentCountFloor)/float32(2)
	if score > maxRiskyScore {
		score = maxRiskyScore
	}

	return &storage.Risk_Result{
		Name: risk.RiskyImageComponent.DisplayTitle,
		Factors: []*storage.Risk_Result_Factor{
			{Message: generateMessage(image.GetName().GetFullName(), riskySet.AsSlice())},
		},
		Score: score,
	}
}

func generateMessage(imageName string, largestRiskySet []string) string {
	return fmt.Sprintf("%s %s", generatePrefix(imageName), generateSuffix(largestRiskySet))
}

func generatePrefix(imageName string) string {
	if imageName != "" {
		return fmt.Sprintf("Image %s", imageName)
	}
	return "An image"
}

func generateSuffix(largestRiskySet []string) string {
	if len(largestRiskySet) == 1 {
		return generateSuffixForOneComponent(largestRiskySet[0])
	}

	// Sort for message stability.
	sort.SliceStable(largestRiskySet, func(i, j int) bool {
		return largestRiskySet[i] < largestRiskySet[j]
	})

	if len(largestRiskySet) <= maxComponentsInMessage {
		return generateSuffixForMultipleButLessThanMax(largestRiskySet)
	}
	return generateSuffixForMoreThanMax(largestRiskySet)
}

func generateSuffixForOneComponent(riskyComponent string) string {
	return fmt.Sprintf("contains component %s", riskyComponent)
}

func generateSuffixForMultipleButLessThanMax(largestRiskySet []string) string {
	componentsString := strings.Join(largestRiskySet, ", ")
	return fmt.Sprintf("contains components useful for attackers: %s", componentsString)
}

func generateSuffixForMoreThanMax(largestRiskySet []string) string {
	componentsInMessage := largestRiskySet[:5]
	componentsString := strings.Join(componentsInMessage, ", ")
	diff := len(largestRiskySet) - maxComponentsInMessage
	return fmt.Sprintf("contains components: %s and %d other(s) that are useful for attackers", componentsString, diff)
}
