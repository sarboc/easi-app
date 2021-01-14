package cedareasi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"

	apiclient "github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/client"
	apioperations "github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/client/operations"
	"github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/models"
	apimodels "github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/models"
	"github.com/cmsgov/easi-app/pkg/testhelpers"
)

func buildClient(t *testing.T) (*apiclient.EASiCoreAPI, runtime.ClientAuthInfoWriter) {
	t.Helper()

	host := os.Getenv("CEDAR_API_URL")
	key := os.Getenv("SDK_CEDAR_KEY")

	if host == "" && key == "" {
		t.Skip("CEDAR SDK test not configured")
	}

	// create the transport
	transport := httptransport.New(host, apiclient.DefaultBasePath, []string{"https"})

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)

	// Set auth header
	auth := httptransport.APIKeyAuth("x-Gateway-APIKey", "header", key)

	return client, auth
}

func TestSDKCedar(t *testing.T) {
	client, auth := buildClient(t)
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	hc := &http.Client{Transport: customTransport}
	ctx := context.Background()

	i1 := testhelpers.NewSystemIntake()
	now := time.Now().UTC()
	i1.SubmittedAt = &now
	i1.ArchivedAt = &now
	i1.DecidedAt = &now
	p1 := apioperations.NewIntakegovernancePOST5ParamsWithContext(ctx).WithHTTPClient(hc)
	p1.Body = &apimodels.Intake{
		Governance: systemIntakeToGovernanceIntake(&i1),
	}
	r1, err := client.Operations.IntakegovernancePOST5(p1, auth)
	if err != nil {
		t.Errorf("e1: %v\n", err)
		return
	}
	t.Logf("Intake1 %s: %s\n\t\t%v\n", i1.ID.String(), *r1.Payload.Response.Result, r1.Payload.Response.Message)

	t.Run("intake identity collisions", func(t *testing.T) {
		_, err = client.Operations.IntakegovernancePOST5(p1, auth)
		assert.Error(t, err)
	})

	t.Run("intake round trip", func(t *testing.T) {
		p2 := apioperations.NewIntakegovernanceidGET6ParamsWithContext(ctx).WithHTTPClient(hc).WithID(i1.ID.String())
		r2, err := client.Operations.IntakegovernanceidGET6(p2, auth)
		assert.NoError(t, err)
		out, err := json.MarshalIndent(*r2.Payload.Intake, "", "\t")
		assert.NoError(t, err)
		t.Logf("Retrieved: %s\n", out)

		p3 := apioperations.NewIntakegovernanceidPUT6ParamsWithContext(ctx).WithHTTPClient(hc).WithID(i1.ID.String())
		p3.Body = &models.IntakeUpdate{
			Governance: r2.Payload.Intake.Governance,
		}
		// p3.Body.Governance = r2.Payload.Intake.Governance
		p3.Body.Governance.SystemName = fmt.Sprintf("%s - %s", time.Now().UTC(), r2.Payload.Intake.Governance.SystemName)
		_, err = client.Operations.IntakegovernanceidPUT6(p3, auth)
		assert.NoError(t, err)

		p4 := apioperations.NewIntakegovernanceidGET6ParamsWithContext(ctx).WithHTTPClient(hc).WithID(i1.ID.String())
		r4, err := client.Operations.IntakegovernanceidGET6(p4, auth)
		assert.NoError(t, err)
		out, err = json.MarshalIndent(*r4.Payload.Intake, "", "\t")
		assert.NoError(t, err)
		t.Logf("Updated: %s\n", out)
	})
}
