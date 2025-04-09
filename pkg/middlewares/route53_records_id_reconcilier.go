package middlewares

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/pkg/resource/aws"
)

// Since AWS returns the FQDN as the name of the remote record, we must change the Id of the
// state record to be equivalent (ZoneId_FQDN_Type_SetIdentifier)
// For a TXT record toto for zone example.com with Id 1234
// From AWS provider, we retrieve: 1234_toto.example.com_TXT
// From Terraform state, we retrieve: 1234_toto_TXT
type Route53RecordIDReconcilier struct{}

func NewRoute53RecordIDReconcilier() Route53RecordIDReconcilier {
	return Route53RecordIDReconcilier{}
}

func (m Route53RecordIDReconcilier) Execute(_, resourcesFromState *[]*resource.Resource) error {

	for _, stateResource := range *resourcesFromState {

		if stateResource.ResourceType() != aws.AwsRoute53RecordResourceType {
			continue
		}

		// Extracting attributes.
		zoneID := (*stateResource.Attrs)["zone_id"].(string)
		fqdn := (*stateResource.Attrs)["fqdn"].(string)
		recordType := (*stateResource.Attrs)["type"].(string)

		// Replacing the wildcard placeholder '\\052' with '*' in fqdn.
		normalizedFqdn := strings.ReplaceAll(fqdn, "\\052", "*")

		// Constructing the new ID.
		newId := strings.Join([]string{zoneID, normalizedFqdn, recordType}, "_")

		if newId != stateResource.Id {
			stateResource.Id = newId
			_ = stateResource.Attrs.SafeSet([]string{"id"}, newId)
			logrus.WithFields(logrus.Fields{
				"old_id": stateResource.ResourceId(),
				"new_id": newId,
			}).Debug("Normalized route53 record ID")
		}
	}

	return nil
}
