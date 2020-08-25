package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"go.uber.org/zap"

	"github.com/cmsgov/easi-app/pkg/appcontext"
	"github.com/cmsgov/easi-app/pkg/apperrors"
	"github.com/cmsgov/easi-app/pkg/models"
)

// NewAuthorizeFetchSystemIntakesByEuaID is a service to authorize FetchSystemIntakesByEuaID
func NewAuthorizeFetchSystemIntakesByEuaID() func(ctx context.Context, euaID string) (bool, error) {
	return func(ctx context.Context, euaID string) (bool, error) {
		return true, nil
	}
}

// NewFetchSystemIntakesByEuaID is a service to fetch system intakes by EUA id
func NewFetchSystemIntakesByEuaID(
	config Config,
	fetch func(c context.Context, euaID string) (models.SystemIntakes, error),
	authorize func(c context.Context, euaID string) (bool, error),
) func(c context.Context, e string) (models.SystemIntakes, error) {
	return func(ctx context.Context, euaID string) (models.SystemIntakes, error) {
		logger := appcontext.ZLogger(ctx)
		ok, err := authorize(ctx, euaID)
		if err != nil {
			logger.Error("failed to authorize fetch system intakes")
			return models.SystemIntakes{}, err
		}
		if !ok {
			return models.SystemIntakes{}, &apperrors.UnauthorizedError{Err: err}
		}
		intakes, err := fetch(ctx, euaID)
		if err != nil {
			logger.Error("failed to fetch system intakes")
			return models.SystemIntakes{}, &apperrors.QueryError{
				Err:       err,
				Model:     intakes,
				Operation: apperrors.QueryFetch,
			}
		}
		return intakes, nil
	}
}

// NewCreateSystemIntake is a service to create a business case
func NewCreateSystemIntake(
	config Config,
	create func(c context.Context, intake *models.SystemIntake) (*models.SystemIntake, error),
) func(c context.Context, i *models.SystemIntake) (*models.SystemIntake, error) {
	return func(ctx context.Context, intake *models.SystemIntake) (*models.SystemIntake, error) {
		logger := appcontext.ZLogger(ctx)
		principal := appcontext.Principal(ctx)
		if !principal.AllowEASi() {
			// Default to failure to authorize and create a quick audit log
			logger.With(zap.Bool("Authorized", false)).
				With(zap.String("Operation", "CreateSystemIntake")).
				Info("something went wrong fetching the eua id from the context")
			return &models.SystemIntake{}, &apperrors.UnauthorizedError{}
		}
		intake.EUAUserID = principal.ID()
		// app validation belongs here
		createdIntake, err := create(ctx, intake)
		if err != nil {
			logger.Error("failed to create a system intake")
			return &models.SystemIntake{}, &apperrors.QueryError{
				Err:       err,
				Model:     intake,
				Operation: apperrors.QueryPost,
			}
		}
		return createdIntake, nil
	}
}

func determineState(existingStatus, newStatus models.SystemIntakeStatus) int {
	if existingStatus == models.SystemIntakeStatusDRAFT {
		if newStatus == models.SystemIntakeStatusDRAFT {
			return models.EditDraftIntake
		} else if newStatus == models.SystemIntakeStatusSUBMITTED {
			return models.SubmitIntake
		}
	} else if existingStatus == models.SystemIntakeStatusSUBMITTED {
		if newStatus == models.SystemIntakeStatusACCEPTED {
			return models.DecideIntakeAccepted
		} else if newStatus == models.SystemIntakeStatusAPPROVED {
			return models.DecideIntakeApproved
		} else if newStatus == models.SystemIntakeStatusCLOSED {
			return models.DecideIntakeClosed
		}
	}
	return models.BadTransition
}

// NewUpdateSystemIntake is a service to update a system intake
func NewUpdateSystemIntake(
	fetch func(c context.Context, id uuid.UUID) (*models.SystemIntake, error),
	canDecideIntake bool,
	updateDRAFTIntake func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error),
	submitIntake func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error),
	decideIntakeACCEPTED func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error),
	decideIntakeAPPROVED func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error),
	decideIntakeCLOSED func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error),
) func(c context.Context, i *models.SystemIntake) (*models.SystemIntake, error) {
	return func(ctx context.Context, intake *models.SystemIntake) (*models.SystemIntake, error) {
		existingIntake, fetchErr := fetch(ctx, intake.ID)
		if fetchErr != nil {
			return &models.SystemIntake{}, &apperrors.QueryError{
				Err:       fetchErr,
				Operation: apperrors.QueryFetch,
				Model:     existingIntake,
			}
		}

		transition := determineState(existingIntake.Status, intake.Status)
		switch transition {
		case models.EditDraftIntake:
			return updateDRAFTIntake(ctx, existingIntake, intake)
		case models.SubmitIntake:
			return submitIntake(ctx, existingIntake, intake)
		case models.DecideIntakeAccepted:
			return decideIntakeACCEPTED(ctx, existingIntake, intake)
		case models.DecideIntakeApproved:
			return decideIntakeAPPROVED(ctx, existingIntake, intake)
		case models.DecideIntakeClosed:
			return decideIntakeCLOSED(ctx, existingIntake, intake)
		default:
			return &models.SystemIntake{}, &apperrors.ResourceConflictError{
				Err:        errors.New("invalid intake status change"),
				Resource:   intake,
				ResourceID: intake.ID.String(),
			}
		}
	}
}

// NewUpdateDRAFTSystemIntake serives DRAFT-DRAFT transition
func NewUpdateDRAFTSystemIntake(
	config Config,
	authorize func(context.Context, *models.SystemIntake) (bool, error),
	update func(c context.Context, intake *models.SystemIntake) (*models.SystemIntake, error),
) func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error) {
	return func(ctx context.Context, existingIntake *models.SystemIntake, updatingIntake *models.SystemIntake) (*models.SystemIntake, error) {
		ok, err := authorize(ctx, existingIntake)
		if err != nil {
			return &models.SystemIntake{}, err
		}
		if !ok {
			return &models.SystemIntake{}, &apperrors.UnauthorizedError{Err: err}
		}

		updatedTime := config.clock.Now()
		updatingIntake.UpdatedAt = &updatedTime

		updatingIntake, err = update(ctx, updatingIntake)
		if err != nil {
			return &models.SystemIntake{}, &apperrors.QueryError{
				Err:       err,
				Model:     updatingIntake,
				Operation: apperrors.QuerySave,
			}
		}
		return updatingIntake, nil
	}
}

// NewSubmitSystemIntake services DRAFT-SUBMIT transition
func NewSubmitSystemIntake(
	config Config,
	authorize func(context.Context, *models.SystemIntake) (bool, error),
	update func(c context.Context, intake *models.SystemIntake) (*models.SystemIntake, error),
	validateAndSubmit func(intake *models.SystemIntake, logger *zap.Logger) (string, error),
	sendSubmitEmail func(requester string, intakeID uuid.UUID) error,
) func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error) {
	return func(ctx context.Context, existingIntake *models.SystemIntake, updatingIntake *models.SystemIntake) (*models.SystemIntake, error) {
		ok, err := authorize(ctx, existingIntake)
		if err != nil {
			return &models.SystemIntake{}, err
		}
		if !ok {
			return &models.SystemIntake{}, &apperrors.UnauthorizedError{Err: err}
		}

		updatedTime := config.clock.Now()
		updatingIntake.UpdatedAt = &updatedTime

		if updatingIntake.AlfabetID.Valid {
			err := &apperrors.ResourceConflictError{
				Err:        errors.New("intake has already been submitted to CEDAR"),
				ResourceID: updatingIntake.ID.String(),
				Resource:   updatingIntake,
			}
			return &models.SystemIntake{}, err
		}

		updatingIntake.SubmittedAt = &updatedTime
		alfabetID, validateAndSubmitErr := validateAndSubmit(updatingIntake, appcontext.ZLogger(ctx))
		if validateAndSubmitErr != nil {
			return &models.SystemIntake{}, validateAndSubmitErr
		}
		if alfabetID == "" {
			return &models.SystemIntake{}, &apperrors.ExternalAPIError{
				Err:       errors.New("submission was not successful"),
				Model:     updatingIntake,
				ModelID:   updatingIntake.ID.String(),
				Operation: apperrors.Submit,
				Source:    "CEDAR EASi",
			}
		}
		updatingIntake.AlfabetID = null.StringFrom(alfabetID)
		updatingIntake, err = update(ctx, updatingIntake)
		if err != nil {
			return &models.SystemIntake{}, &apperrors.QueryError{
				Err:       err,
				Model:     updatingIntake,
				Operation: apperrors.QuerySave,
			}
		}
		// only send an email when everything went ok
		err = sendSubmitEmail(updatingIntake.Requester, updatingIntake.ID)
		if err != nil {
			return &models.SystemIntake{}, err
		}

		return updatingIntake, nil
	}
}

// NewDecideSystemIntake services SUBMITTED-(ACCEPTED/APPROVED/CLOSED) transition
func NewDecideSystemIntake(
	config Config,
	authorize func(context.Context, *models.SystemIntake) (bool, error),
	fetchRequesterEmail func(logger *zap.Logger, euaID string) (string, error),
	update func(c context.Context, intake *models.SystemIntake) (*models.SystemIntake, error),
	sendReviewEmail func(emailText string, recipientAddress string) error,
) func(context.Context, *models.SystemIntake, *models.SystemIntake) (*models.SystemIntake, error) {
	return func(ctx context.Context, existingIntake *models.SystemIntake, updatingIntake *models.SystemIntake) (*models.SystemIntake, error) {
		ok, err := authorize(ctx, existingIntake)
		if err != nil {
			return &models.SystemIntake{}, err
		}
		if !ok {
			return &models.SystemIntake{}, &apperrors.UnauthorizedError{Err: err}
		}

		updatedTime := config.clock.Now()
		updatingIntake.UpdatedAt = &updatedTime

		recipientAddress, err := fetchRequesterEmail(appcontext.ZLogger(ctx), existingIntake.EUAUserID)
		if err != nil {
			return &models.SystemIntake{}, err
		}
		if recipientAddress == "" {
			return &models.SystemIntake{}, &apperrors.ExternalAPIError{
				Err:       errors.New("email address fetch was not successful"),
				Model:     existingIntake,
				ModelID:   updatingIntake.ID.String(),
				Operation: apperrors.Fetch,
				Source:    "CEDAR LDAP",
			}
		}

		existingIntake.Status = updatingIntake.Status
		existingIntake.GrtReviewEmailBody = updatingIntake.GrtReviewEmailBody
		existingIntake.RequesterEmailAddress = null.StringFrom(recipientAddress)
		existingIntake.DecidedAt = &updatedTime
		existingIntake.UpdatedAt = &updatedTime
		// This ensures only certain fields can be modified.
		updatingIntake, err = update(ctx, existingIntake)
		if err != nil {
			return &models.SystemIntake{}, &apperrors.QueryError{
				Err:       err,
				Model:     updatingIntake,
				Operation: apperrors.QuerySave,
			}
		}

		err = sendReviewEmail(updatingIntake.GrtReviewEmailBody.String, recipientAddress)
		if err != nil {
			return &models.SystemIntake{}, err
		}

		return updatingIntake, nil
	}
}

// NewAuthorizeUserIsIntakeRequester returns a function
// that authorizes a user for archiving a system intake
func NewAuthorizeUserIsIntakeRequester(logger *zap.Logger) func(
	c context.Context,
	i *models.SystemIntake,
) (bool, error) {
	return func(ctx context.Context, intake *models.SystemIntake) (bool, error) {
		principal := appcontext.Principal(ctx)
		if !principal.AllowEASi() {
			logger.Error("unable to get EUA ID from context")
			return false, &apperrors.ContextError{
				Operation: apperrors.ContextGet,
				Object:    "EUA ID",
			}
		}

		// If intake doesn't exist or owned by user, authorize
		if intake == nil || principal.ID() == intake.EUAUserID {
			logger.With(zap.Bool("Authorized", true)).
				With(zap.String("Operation", "UpdateSystemIntake")).
				Info("user authorized to save system intake")
			return true, nil
		}
		// Default to failure to authorize and create a quick audit log
		logger.With(zap.Bool("Authorized", false)).
			With(zap.String("Operation", "UpdateSystemIntake")).
			Info("unauthorized attempt to save system intake")
		return false, nil
	}
}

// NewAuthorizeUserIsGRT is a service to authorize FetchSystemIntakesByEuaID
func NewAuthorizeUserIsGRT() func(context.Context, *models.SystemIntake) (bool, error) {
	return func(ctx context.Context, intake *models.SystemIntake) (bool, error) {
		return true, nil
	}
}

// NewArchiveSystemIntake is a service to archive a system intake
func NewArchiveSystemIntake(
	config Config,
	fetch func(c context.Context, id uuid.UUID) (*models.SystemIntake, error),
	update func(c context.Context, intake *models.SystemIntake) (*models.SystemIntake, error),
	archiveBusinessCase func(context.Context, uuid.UUID) error,
	authorize func(context context.Context, intake *models.SystemIntake) (bool, error),
) func(context.Context, uuid.UUID) error {
	return func(ctx context.Context, id uuid.UUID) error {
		intake, fetchErr := fetch(ctx, id)
		if fetchErr != nil {
			return &apperrors.QueryError{
				Err:       fetchErr,
				Operation: apperrors.QueryFetch,
				Model:     intake,
			}
		}
		ok, err := authorize(ctx, intake)
		if err != nil {
			return err
		}
		if !ok {
			return &apperrors.UnauthorizedError{Err: err}
		}

		// We need to archive any associated business case
		if intake.BusinessCaseID != nil {
			err = archiveBusinessCase(ctx, *intake.BusinessCaseID)
			if err != nil {
				return err
			}
		}

		updatedTime := config.clock.Now()
		intake.UpdatedAt = &updatedTime
		intake.Status = models.SystemIntakeStatusARCHIVED
		intake.ArchivedAt = &updatedTime

		intake, err = update(ctx, intake)
		if err != nil {
			return &apperrors.QueryError{
				Err:       err,
				Model:     intake,
				Operation: apperrors.QuerySave,
			}
		}

		return nil
	}
}

// NewAuthorizeFetchSystemIntakeByID is a service to authorize FetchSystemIntakeByID
func NewAuthorizeFetchSystemIntakeByID() func(ctx context.Context, intake *models.SystemIntake) (bool, error) {
	return func(ctx context.Context, intake *models.SystemIntake) (bool, error) {
		return true, nil
	}
}

// NewFetchSystemIntakeByID is a service to fetch the system intake by intake id
func NewFetchSystemIntakeByID(
	config Config,
	fetch func(c context.Context, id uuid.UUID) (*models.SystemIntake, error),
	authorize func(c context.Context, i *models.SystemIntake) (bool, error),
) func(c context.Context, u uuid.UUID) (*models.SystemIntake, error) {
	return func(ctx context.Context, id uuid.UUID) (*models.SystemIntake, error) {
		logger := appcontext.ZLogger(ctx)
		intake, err := fetch(ctx, id)
		if err != nil {
			logger.Error("failed to fetch system intake")
			return &models.SystemIntake{}, &apperrors.QueryError{
				Err:       err,
				Model:     intake,
				Operation: apperrors.QueryFetch,
			}
		}
		ok, err := authorize(ctx, intake)
		if err != nil {
			logger.Error("failed to authorize fetch system intake")
			return &models.SystemIntake{}, err
		}
		if !ok {
			return &models.SystemIntake{}, &apperrors.UnauthorizedError{Err: err}
		}
		return intake, nil
	}
}
