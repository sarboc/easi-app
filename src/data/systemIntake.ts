import { DateTime } from 'luxon';

import cmsGovernanceTeams from 'constants/enums/cmsGovernanceTeams';
import {
  GovernanceCollaborationTeam,
  SystemIntakeForm
} from 'types/systemIntake';

// On the frontend, the field is now "requestName", but the backend API
// has it as "projectName". This was an update from design.
export const initialSystemIntakeForm: SystemIntakeForm = {
  id: '',
  euaUserID: '',
  requestName: '',
  status: 'INTAKE_DRAFT',
  requestType: 'NEW',
  requester: {
    name: '',
    component: '',
    email: ''
  },
  businessOwner: {
    name: '',
    component: ''
  },
  productManager: {
    name: '',
    component: ''
  },
  isso: {
    isPresent: null,
    name: ''
  },
  governanceTeams: {
    isPresent: null,
    teams: []
  },
  fundingSource: {
    isFunded: null,
    fundingNumber: '',
    source: ''
  },
  costs: {
    isExpectingIncrease: '',
    expectedIncreaseAmount: ''
  },
  contract: {
    hasContract: '',
    contractor: '',
    vehicle: '',
    startDate: {
      month: '',
      year: ''
    },
    endDate: {
      month: '',
      year: ''
    }
  },
  businessNeed: '',
  businessSolution: '',
  currentStage: '',
  needsEaSupport: null,
  grtReviewEmailBody: '',
  decidedAt: null,
  businessCaseId: null,
  submittedAt: null,
  updatedAt: null,
  createdAt: null,
  archivedAt: null,
  lcid: ''
};

export const prepareSystemIntakeForApi = (systemIntake: SystemIntakeForm) => {
  const getGovernanceCollaborator = (name: string) => {
    const selectedTeam = systemIntake.governanceTeams.teams.find(
      team => team.name === name
    );

    return selectedTeam ? selectedTeam.collaborator : '';
  };

  return {
    ...(systemIntake.id && {
      id: systemIntake.id
    }),
    status: systemIntake.status,
    requestType: systemIntake.requestType,
    requester: systemIntake.requester.name,
    component: systemIntake.requester.component,
    businessOwner: systemIntake.businessOwner.name,
    businessOwnerComponent: systemIntake.businessOwner.component,
    productManager: systemIntake.productManager.name,
    productManagerComponent: systemIntake.productManager.component,
    isso: systemIntake.isso.name,
    trbCollaborator: getGovernanceCollaborator('Technical Review Board'),
    oitSecurityCollaborator: getGovernanceCollaborator(
      "OIT's Security and Privacy Group"
    ),
    eaCollaborator: getGovernanceCollaborator('Enterprise Architecture'),
    projectName: systemIntake.requestName,
    existingFunding: systemIntake.fundingSource.isFunded,
    fundingNumber: systemIntake.fundingSource.fundingNumber,
    fundingSource: systemIntake.fundingSource.source,
    businessNeed: systemIntake.businessNeed,
    solution: systemIntake.businessSolution,
    processStatus: systemIntake.currentStage,
    eaSupportRequest: systemIntake.needsEaSupport,
    existingContract: systemIntake.contract.hasContract,
    grtReviewEmailBody: systemIntake.grtReviewEmailBody,
    costIncrease: systemIntake.costs.isExpectingIncrease,
    costIncreaseAmount: systemIntake.costs.expectedIncreaseAmount,
    contractor: systemIntake.contract.contractor,
    contractVehicle: systemIntake.contract.vehicle,
    contractStartMonth: systemIntake.contract.startDate.month,
    contractStartYear: systemIntake.contract.startDate.year,
    contractEndMonth: systemIntake.contract.endDate.month,
    contractEndYear: systemIntake.contract.endDate.year
  };
};

export const prepareSystemIntakeForApp = (
  systemIntake: any
): SystemIntakeForm => {
  const governanceTeams = () => {
    const teams: GovernanceCollaborationTeam[] = [];
    cmsGovernanceTeams.forEach(team => {
      if (systemIntake[team.collaboratorKey]) {
        teams.push({
          collaborator: systemIntake[team.collaboratorKey],
          name: team.value
        });
      }
    });
    return teams;
  };

  return {
    id: systemIntake.id || '',
    euaUserID: systemIntake.euaUserID || '',
    requestName: systemIntake.projectName || '',
    status: systemIntake.status || 'INTAKE_DRAFT',
    requestType: systemIntake.requestType || 'NEW',
    requester: {
      name: systemIntake.requester || '',
      component: systemIntake.component || '',
      email: systemIntake.requesterEmailAddress || ''
    },
    businessOwner: {
      name: systemIntake.businessOwner || '',
      component: systemIntake.businessOwnerComponent || ''
    },
    productManager: {
      name: systemIntake.productManager || '',
      component: systemIntake.productManagerComponent || ''
    },
    isso: {
      isPresent: !!systemIntake.isso || null,
      name: systemIntake.isso || ''
    },
    governanceTeams: {
      isPresent: governanceTeams().length !== 0 || null,
      teams: governanceTeams() || []
    },
    fundingSource: {
      isFunded:
        systemIntake.existingFunding === null
          ? null
          : systemIntake.existingFunding,
      fundingNumber: systemIntake.fundingNumber || '',
      source: systemIntake.fundingSource || ''
    },
    costs: {
      isExpectingIncrease: systemIntake.costIncrease || '',
      expectedIncreaseAmount: systemIntake.costIncreaseAmount || ''
    },
    contract: {
      hasContract: systemIntake.existingContract || '',
      contractor: systemIntake.contractor || '',
      vehicle: systemIntake.contractVehicle || '',
      startDate: {
        month: systemIntake.contractStartMonth || '',
        year: systemIntake.contractStartYear || ''
      },
      endDate: {
        month: systemIntake.contractEndMonth || '',
        year: systemIntake.contractEndYear || ''
      }
    },
    businessNeed: systemIntake.businessNeed || '',
    businessSolution: systemIntake.solution || '',
    currentStage: systemIntake.processStatus || '',
    needsEaSupport:
      systemIntake.eaSupportRequest === null
        ? null
        : systemIntake.eaSupportRequest,
    grtReviewEmailBody: systemIntake.grtReviewEmailBody || '',
    decidedAt: systemIntake.decidedAt
      ? DateTime.fromISO(systemIntake.decidedAt)
      : null,
    businessCaseId: systemIntake.businessCase || null,
    submittedAt: systemIntake.submittedAt
      ? DateTime.fromISO(systemIntake.submittedAt)
      : null,
    updatedAt: systemIntake.updatedAt
      ? DateTime.fromISO(systemIntake.updatedAt)
      : null,
    createdAt: systemIntake.createdAt
      ? DateTime.fromISO(systemIntake.createdAt)
      : null,
    archivedAt: systemIntake.archivedAt
      ? DateTime.fromISO(systemIntake.archivedAt)
      : null,
    lcid: systemIntake.lcid || ''
  };
};

export const convertIntakeToCSV = (intake: SystemIntakeForm) => {
  const collaboratorTeams: any = {};
  if (intake.governanceTeams.isPresent) {
    intake.governanceTeams.teams.forEach(team => {
      switch (team.name) {
        case 'Technical Review Board':
          collaboratorTeams.trbCollaborator = team.collaborator;
          break;
        case "OIT's Security and Privacy Group":
          collaboratorTeams.oitCollaborator = team.collaborator;
          break;
        case 'Enterprise Architecture':
          collaboratorTeams.eaCollaborator = team.collaborator;
          break;
        default:
          break;
      }
    });
  }
  return {
    ...intake,
    ...collaboratorTeams,
    contractStartDate: ['HAVE_CONTRACT', 'IN_PROGRESS'].includes(
      intake.contract.hasContract
    )
      ? `${intake.contract.startDate.month}/${intake.contract.startDate.year}`
      : '',
    contractEndDate: ['HAVE_CONTRACT', 'IN_PROGRESS'].includes(
      intake.contract.hasContract
    )
      ? `${intake.contract.endDate.month}/${intake.contract.endDate.year}`
      : '',
    submittedAt: intake.submittedAt && intake.submittedAt.toISO(),
    updatedAt: intake.updatedAt && intake.updatedAt.toISO(),
    createdAt: intake.createdAt && intake.createdAt.toISO(),
    decidedAt: intake.decidedAt && intake.decidedAt.toISO(),
    archivedAt: intake.archivedAt && intake.archivedAt.toISO()
  };
};
