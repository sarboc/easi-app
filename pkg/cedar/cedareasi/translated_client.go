package cedareasi

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/cmsgov/easi-app/pkg/apperrors"
	"github.com/cmsgov/easi-app/pkg/appvalidation"
	apiclient "github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/client"
	apioperations "github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/client/operations"
	apimodels "github.com/cmsgov/easi-app/pkg/cedar/cedareasi/gen/models"
	"github.com/cmsgov/easi-app/pkg/models"
	"github.com/cmsgov/easi-app/pkg/validate"
)

// TranslatedClient is an API client for CEDAR EASi using EASi language
type TranslatedClient struct {
	client        *apiclient.EASiCore
	apiAuthHeader runtime.ClientAuthInfoWriter
}

// NewTranslatedClient returns an API client for CEDAR EASi using EASi language
func NewTranslatedClient(cedarHost string, cedarAPIKey string) TranslatedClient {
	// create the transport
	transport := httptransport.New(cedarHost, apiclient.DefaultBasePath, []string{"https"})

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)

	// Set auth header
	apiKeyHeaderAuth := httptransport.APIKeyAuth("x-Gateway-APIKey", "header", cedarAPIKey)

	return TranslatedClient{client, apiKeyHeaderAuth}
}

// FetchSystems fetches a system list from CEDAR
func (c TranslatedClient) FetchSystems(logger *zap.Logger) (models.SystemShorts, error) {
	resp, err := c.client.Operations.SystemsGET1(nil, c.apiAuthHeader)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch system from CEDAR with error: %v", err))
		return models.SystemShorts{}, err
	}

	systems := make([]models.SystemShort, len(resp.Payload.Systems))
	for index, system := range resp.Payload.Systems {
		systems[index] = models.SystemShort{
			ID:      *system.ID,
			Name:    *system.SystemName,
			Acronym: system.SystemAcronym,
		}
	}
	return systems, nil
}

// ValidateSystemIntakeForCedar validates all required fields to ensure we won't get errors for contents of the request
func ValidateSystemIntakeForCedar(intake *models.SystemIntake, logger *zap.Logger) error {
	expectedError := apperrors.ValidationError{
		Err:         errors.New("validation failed"),
		Validations: apperrors.Validations{},
		ModelID:     intake.ID.String(),
		Model:       intake,
	}
	const validationMessage = "is required"
	if validate.RequireUUID(intake.ID) {
		expectedError.WithValidation("ID", validationMessage)
	}
	if validate.RequireString(intake.EUAUserID) {
		expectedError.WithValidation("EUAUserID", validationMessage)
	}
	if validate.RequireString(string(intake.Status)) {
		expectedError.WithValidation("Status", validationMessage)
	}
	if validate.RequireString(intake.Requester) {
		expectedError.WithValidation("Requester", validationMessage)
	}
	if validate.RequireNullString(intake.Component) {
		expectedError.WithValidation("Component", validationMessage)
	}
	if validate.RequireNullString(intake.BusinessOwner) {
		expectedError.WithValidation("BusinessOwner", validationMessage)
	}
	if validate.RequireNullString(intake.BusinessOwnerComponent) {
		expectedError.WithValidation("BusinessOwnerComponent", validationMessage)
	}
	if validate.RequireNullString(intake.ProductManager) {
		expectedError.WithValidation("ProductManager", validationMessage)
	}
	if validate.RequireNullString(intake.ProductManagerComponent) {
		expectedError.WithValidation("ProductManagerComponent", validationMessage)
	}
	if validate.RequireNullString(intake.ProjectName) {
		expectedError.WithValidation("ProjectName", validationMessage)
	}
	if validate.RequireNullBool(intake.ExistingFunding) {
		expectedError.WithValidation("ExistingFunding", validationMessage)
	}
	if intake.ExistingFunding.Bool {
		if validate.RequireNullString(intake.FundingSource) {
			expectedError.WithValidation("FundingSource", validationMessage)
		}
		if intake.FundingSource.Valid && validate.FundingNumberInvalid(intake.FundingSource.String) {
			expectedError.WithValidation("FundingSource", "must be a 6 digit string")
		}
	}
	if validate.RequireNullString(intake.BusinessNeed) {
		expectedError.WithValidation("BusinessNeed", validationMessage)
	}
	if validate.RequireNullString(intake.Solution) {
		expectedError.WithValidation("Solution", validationMessage)
	}
	if validate.RequireNullBool(intake.EASupportRequest) {
		expectedError.WithValidation("EASupportRequest", validationMessage)
	}
	if validate.RequireNullString(intake.ProcessStatus) {
		expectedError.WithValidation("ProcessStatus", validationMessage)
	}
	if validate.RequireNullString(intake.ExistingContract) {
		expectedError.WithValidation("ExistingContract", validationMessage)
	}
	if validate.RequireTime(*intake.SubmittedAt) {
		expectedError.WithValidation("SubmittedAt", validationMessage)
	}
	if len(expectedError.Validations) > 0 {
		return &expectedError
	}
	return nil
}

func submitSystemIntake(validatedIntake *models.SystemIntake, c TranslatedClient, logger *zap.Logger) (string, error) {
	id := validatedIntake.ID.String()
	submissionTime := validatedIntake.SubmittedAt.String()
	statusAsString := string(validatedIntake.Status)
	params := apioperations.NewIntakegovernancePOST5Params()
	governanceIntake := apimodels.GovernanceIntake{
		BusinessNeeds:           &validatedIntake.BusinessNeed.String,
		BusinessOwner:           &validatedIntake.BusinessOwner.String,
		BusinessOwnerComponent:  &validatedIntake.BusinessOwnerComponent.String,
		EaCollaborator:          validatedIntake.EACollaborator.String,
		EaSupportRequest:        &validatedIntake.EASupportRequest.Bool,
		EuaUserID:               &validatedIntake.EUAUserID,
		ExistingContract:        &validatedIntake.ExistingContract.String,
		ExistingFunding:         &validatedIntake.ExistingFunding.Bool,
		FundingSource:           validatedIntake.FundingSource.String,
		ID:                      &id,
		Isso:                    validatedIntake.ISSO.String,
		OitSecurityCollaborator: validatedIntake.OITSecurityCollaborator.String,
		ProcessStatus:           &validatedIntake.ProcessStatus.String,
		ProductManager:          &validatedIntake.ProductManager.String,
		ProductManagerComponent: &validatedIntake.ProductManagerComponent.String,
		Requester:               &validatedIntake.Requester,
		RequesterComponent:      &validatedIntake.Component.String,
		Solution:                &validatedIntake.Solution.String,
		Status:                  &statusAsString,
		SubmittedAt:             &submissionTime,
		SubmittedTime:           &submissionTime,
		SystemName:              &validatedIntake.ProjectName.String,
		TrbCollaborator:         validatedIntake.TRBCollaborator.String,
	}
	governanceConversion := []*apimodels.GovernanceIntake{
		&governanceIntake,
	}
	params.Body = &apimodels.Intake{
		Governance: governanceConversion,
	}
	resp, err := c.client.Operations.IntakegovernancePOST5(params, c.apiAuthHeader)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to submit intake for CEDAR with error: %v", err))
		return "", &apperrors.ExternalAPIError{
			Err:       err,
			Model:     validatedIntake,
			ModelID:   id,
			Operation: apperrors.Submit,
			Source:    "CEDAR",
		}
	}
	alfabetID := ""
	if *resp.Payload.Response.Result != "success" {
		return "", &apperrors.ExternalAPIError{
			Err:       errors.New("CEDAR return result: " + *resp.Payload.Response.Result),
			ModelID:   validatedIntake.ID.String(),
			Model:     validatedIntake,
			Operation: apperrors.Submit,
			Source:    "CEDAR",
		}
	}
	alfabetID = resp.Payload.Response.Message[0]
	return alfabetID, nil
}

// ValidateAndSubmitSystemIntake submits a system intake to CEDAR
func (c TranslatedClient) ValidateAndSubmitSystemIntake(intake *models.SystemIntake, logger *zap.Logger) (string, error) {
	err := ValidateSystemIntakeForCedar(intake, logger)
	if err != nil {
		return "", err
	}
	return submitSystemIntake(intake, c, logger)
}

// ValidateBusinessCaseForCedar validates all required fields to ensure we won't get errors for contents of the request
func (c TranslatedClient) ValidateBusinessCaseForCedar(businessCase *models.BusinessCase) error {
	expectedError := apperrors.ValidationError{
		Err:         errors.New("validation failed"),
		Validations: apperrors.Validations{},
		ModelID:     businessCase.ID.String(),
		Model:       businessCase,
	}
	const validationMessage = "is required"
	if validate.RequireNullString(businessCase.BusinessNeed) {
		expectedError.WithValidation("BusinessNeed", validationMessage)
	}
	if validate.RequireNullString(businessCase.BusinessOwner) {
		expectedError.WithValidation("BusinessOwner", validationMessage)
	}
	if validate.RequireNullString(businessCase.CMSBenefit) {
		expectedError.WithValidation("CMSBenefit", validationMessage)
	}
	// ToDo: decided_at is required in the swagger. Should it be?
	if validate.RequireString(businessCase.EUAUserID) {
		expectedError.WithValidation("EUAUserID", validationMessage)
	}
	// ToDo: hosting_needs is required in swagger. What field is that?
	if validate.RequireUUID(businessCase.ID) {
		expectedError.WithValidation("ID", validationMessage)
	}
	// ToDo: lifecycle_id is required in swagger. Should it be?
	if validate.RequireNullString(businessCase.ProjectName) {
		expectedError.WithValidation("ProjectName", validationMessage)
	}
	if validate.RequireNullString(businessCase.Requester) {
		expectedError.WithValidation("Requester", validationMessage)
	}
	if validate.RequireNullString(businessCase.ProjectName) {
		expectedError.WithValidation("ProjectName", validationMessage)
	}
	if validate.RequireNullString(businessCase.RequesterPhoneNumber) {
		expectedError.WithValidation("RequesterPhoneNumber", validationMessage)
	}
	if validate.RequireString(string(businessCase.Status)) {
		expectedError.WithValidation("Status", validationMessage)
	}
	if validate.RequireTime(*businessCase.InitialSubmittedAt) {
		expectedError.WithValidation("InitialSubmittedAt", validationMessage)
	}
	if validate.RequireNullString(businessCase.SuccessIndicators) {
		expectedError.WithValidation("SuccessIndicators", validationMessage)
	}
	// ToDo: user_interface is required in swagger. What is that?
	// ToDo: withdrawn_at is required in swagger. Should it be?
	if validate.RequireNullString(businessCase.AsIsTitle) {
		expectedError.WithValidation("AsIsTitle", validationMessage)
	}
	if validate.RequireNullString(businessCase.AsIsSummary) {
		expectedError.WithValidation("AsIsSummary", validationMessage)
	}
	if validate.RequireNullString(businessCase.AsIsPros) {
		expectedError.WithValidation("AsIsPros", validationMessage)
	}
	if validate.RequireNullString(businessCase.AsIsCons) {
		expectedError.WithValidation("AsIsCons", validationMessage)
	}
	if validate.RequireNullString(businessCase.PreferredTitle) {
		expectedError.WithValidation("PreferredTitle", validationMessage)
	}
	if validate.RequireNullString(businessCase.PreferredSummary) {
		expectedError.WithValidation("PreferredSummary", validationMessage)
	}
	if validate.RequireNullString(businessCase.PreferredPros) {
		expectedError.WithValidation("PreferredPros", validationMessage)
	}
	if validate.RequireNullString(businessCase.PreferredCons) {
		expectedError.WithValidation("PreferredCons", validationMessage)
	}
	if validate.RequireNullString(businessCase.AlternativeATitle) {
		expectedError.WithValidation("AlternativeATitle", validationMessage)
	}
	if validate.RequireNullString(businessCase.AlternativeASummary) {
		expectedError.WithValidation("AlternativeASummary", validationMessage)
	}
	if validate.RequireNullString(businessCase.AlternativeAPros) {
		expectedError.WithValidation("AlternativeAPros", validationMessage)
	}
	if validate.RequireNullString(businessCase.AlternativeACons) {
		expectedError.WithValidation("AlternativeACons", validationMessage)
	}
	if appvalidation.IsAlternativeBPresent(businessCase) {
		if validate.RequireNullString(businessCase.AlternativeBTitle) {
			expectedError.WithValidation("AlternativeBTitle", validationMessage)
		}
		if validate.RequireNullString(businessCase.AlternativeBSummary) {
			expectedError.WithValidation("AlternativeBSummary", validationMessage)
		}
		if validate.RequireNullString(businessCase.AlternativeBPros) {
			expectedError.WithValidation("AlternativeBPros", validationMessage)
		}
		if validate.RequireNullString(businessCase.AlternativeBCons) {
			expectedError.WithValidation("AlternativeBCons", validationMessage)
		}
	}
	if len(expectedError.Validations) > 0 {
		return &expectedError
	}
	return nil
}
