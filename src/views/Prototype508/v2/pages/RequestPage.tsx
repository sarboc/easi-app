import React, { useState } from 'react';
import { useHistory, useLocation, useParams } from 'react-router-dom';
import { TabPanel } from '@cmsgov/design-system/dist/esnext/Tabs/TabPanel';
import { Tabs } from '@cmsgov/design-system/dist/esnext/Tabs/Tabs';
import { Button, Table } from '@trussworks/react-uswds';
import { DateTime } from 'luxon';

import Modal from 'components/Modal';

import DateField from '../components/DateField';
import {
  ProgressIndicator,
  ProgressStatus,
  ProgressStep
} from '../components/Progress';
import RequestStatusField from '../components/RequestStatusField';
import SecondaryNavigation from '../components/SecondaryNavigation';
import useDocumentTitle from '../hooks/DocumentTitle';
import { useGlobalState } from '../state';
import {
  Document,
  DocumentType,
  Note,
  Project,
  RequestStep,
  RequestStepStatus
} from '../types';

import '@cmsgov/design-system/dist/css/index.css';
import './index.scss';

const AddTestDateModal = ({
  document,
  updateDocument
}: {
  document: Document;
  updateDocument: (document: Document) => void;
}) => {
  const [dateModalIsOpen, setDateModalIsOpen] = useState(false);

  return (
    <>
      <button
        type="button"
        className="usa-button usa-button--unstyled"
        onClick={() => {
          setDateModalIsOpen(true);
        }}
      >
        Add date
        <span className="usa-sr-only">to {documentName(document)}</span>
      </button>

      <Modal
        title="Add a test date"
        isOpen={dateModalIsOpen}
        closeModal={() => {
          setDateModalIsOpen(false);
        }}
        className="add-date-modal"
      >
        <DateField
          setDate={d => {
            // eslint-disable-next-line no-param-reassign
            document.testDate = d;
          }}
        />
        <button
          type="submit"
          className="usa-button"
          onClick={() => {
            updateDocument(document);
            setDateModalIsOpen(false);
          }}
        >
          Add date
        </button>
        <button
          type="button"
          className="usa-button usa-button--unstyled"
          onClick={() => {
            setDateModalIsOpen(false);
          }}
        >
          Don&rsquo;t add a test date
        </button>
      </Modal>
    </>
  );
};

const UpdateStatusModal = ({
  project,
  updateProject
}: {
  project: Project;
  updateProject: (project: Project) => void;
}) => {
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [projectStatus, setProjectStatus] = useState(project.status);
  const [date, setDate] = useState<DateTime>();

  return (
    <>
      <button
        type="button"
        className="usa-button usa-button--unstyled"
        onClick={() => setModalIsOpen(true)}
      >
        Change Status
      </button>
      <Modal
        title="Change Project Status"
        isOpen={modalIsOpen}
        closeModal={() => {
          setModalIsOpen(false);
        }}
        className="status-modal"
      >
        <div className="status-modal__content">
          <RequestStatusField
            projectStatus={projectStatus}
            projectName={project.name}
            setProjectStatus={setProjectStatus}
            setDate={setDate}
          />

          <div className="status-modal__footer">
            <p className="usa-prose">
              Changing the project status will send an email to all members of
              the 508 team letting them know about the new status.
            </p>

            <button
              type="submit"
              className="usa-button"
              onClick={() => {
                // eslint-disable-next-line no-param-reassign
                project.status = projectStatus;
                // eslint-disable-next-line no-param-reassign
                project.lastUpdatedAt = DateTime.local();
                Object.entries(project.stepStatuses).forEach(
                  ([requestStep, stepStatus]) => {
                    if (!stepStatus) {
                      return;
                    }
                    if (requestStep === projectStatus) {
                      // eslint-disable-next-line no-param-reassign
                      stepStatus.date = date;
                      // eslint-disable-next-line no-param-reassign
                      stepStatus.status = ProgressStatus.Current;
                    } else {
                      // eslint-disable-next-line no-param-reassign
                      stepStatus.status = ProgressStatus.Completed;
                    }
                  }
                );

                if (!project.stepStatuses[projectStatus]) {
                  // eslint-disable-next-line no-param-reassign
                  project.stepStatuses[projectStatus] = {
                    status: ProgressStatus.Current,
                    date
                  };
                }
                updateProject(project);
                setModalIsOpen(false);
              }}
            >
              Change status and send email
            </button>
            <button
              type="button"
              className="usa-button usa-button--unstyled"
              onClick={() => {
                setModalIsOpen(false);
              }}
            >
              Don&rsquo;t change projects status
            </button>
          </div>
        </div>
      </Modal>
    </>
  );
};

const documentName = (doc: Document) => {
  if (doc.type === DocumentType.Other) {
    return doc.otherName || '';
  }
  return doc.type.toString();
};

const DocumentTable = ({
  project,
  updateProject
}: {
  project: Project;
  updateProject: (project: Project) => void;
}) => {
  return (
    <Table bordered={false} fullWidth>
      <caption className="usa-sr-only">
        List of documents uploaded for {project.name}
      </caption>
      <thead>
        <tr>
          <th scope="col" style={{ whiteSpace: 'nowrap' }}>
            Document
          </th>
          <th scope="col" style={{ whiteSpace: 'nowrap' }}>
            Date Uploaded
          </th>
          <th scope="col" style={{ whiteSpace: 'nowrap' }}>
            Test Date
          </th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {project.documents
          .sort((a, b) => a.createdAt.toSeconds() - b.createdAt.toSeconds())
          .map(doc => {
            return (
              <tr key={doc.id}>
                <th scope="row">
                  {documentName(doc)}{' '}
                  {doc.type === DocumentType.TestResults && (
                    <span>- {doc.score}%</span>
                  )}
                </th>
                <td>{doc.createdAt.toFormat('LLLL d y')}</td>
                <td>
                  {doc.testDate ? (
                    doc.testDate.toFormat('LLLL d y')
                  ) : (
                    <AddTestDateModal
                      document={doc}
                      updateDocument={() => updateProject(project)}
                    />
                  )}
                </td>
                <td>
                  <button
                    type="button"
                    className="usa-button usa-button--unstyled margin-right-2"
                  >
                    View{' '}
                    <span className="usa-sr-only">{documentName(doc)}</span>
                  </button>
                  <button
                    type="button"
                    className="usa-button usa-button--unstyled"
                    onClick={() => {
                      project.documents.splice(
                        project.documents.indexOf(doc),
                        1
                      );
                      // eslint-disable-next-line no-param-reassign
                      project.banner = `${documentName(
                        doc
                      )} was removed from the project.`;
                      updateProject(project);
                    }}
                  >
                    Remove{' '}
                    <span className="usa-sr-only">{documentName(doc)}</span>
                  </button>
                </td>
              </tr>
            );
          })}
      </tbody>
    </Table>
  );
};

const getSelectedTabId = (hash: string) => {
  const id = hash.replace('#', '');
  if (['documents', 'notes', 'details'].includes(id)) {
    return id;
  }
  return 'documents';
};

const Timeline = ({ project }: { project: Project }) => {
  let encounteredCurrent = false;
  return (
    <ProgressIndicator>
      {Object.values(RequestStep).map(value => {
        const step = value as RequestStep;

        const status = project.stepStatuses[step] || {
          status: encounteredCurrent
            ? ProgressStatus.NotCompleted
            : ProgressStatus.Skipped
        };

        if (status.status === ProgressStatus.Current) {
          encounteredCurrent = true;
        }

        return (
          <ProgressStep name={step} status={status.status}>
            {stepContent(step, status)}
          </ProgressStep>
        );
      })}
    </ProgressIndicator>
  );
};

const ProjectPage = () => {
  const { id } = useParams();
  const { hash, pathname } = useLocation();
  const history = useHistory();

  const { state, updateProject } = useGlobalState();
  const project = state.projects[id];
  const [noteContent, setNoteContent] = useState('');
  hash.replace('#', '');

  useDocumentTitle(`EASi: Project page for ${project && project.name}`);

  if (!project) {
    return <main>Project not found</main>;
  }

  return (
    <>
      <SecondaryNavigation />
      <main
        id="main-content"
        className="easi-main-content margin-bottom-5"
        aria-label={`Project page for ${project.name}`}
      >
        <div className="grid-container">
          {project.banner && (
            <div
              className="usa-alert usa-alert--success usa-alert--slim margin-bottom-2 margin-top-2"
              role="alert"
            >
              <div className="usa-alert__body">
                <p className="usa-alert__text">{project.banner}</p>
              </div>
            </div>
          )}

          <h1>
            {project.name} {project.release}
          </h1>
        </div>

        <div className="grid-container">
          <div className="grid-row grid-gap-lg easi-column-reverse">
            <div className="grid-col-4">
              <h3>Timeline</h3>
              <div
                className="easi-grt__status-info text-gray-90 padding-top-1 padding-bottom-1"
                aria-label={`Status for ${project.name}`}
                aria-describedby="timeline-description"
              >
                <div className="usa-sr-only" id="timeline-description">
                  The timeline indicates where this request is within the 508
                  process and let’s you change the status.
                </div>
                <UpdateStatusModal
                  project={project}
                  updateProject={updateProject}
                />
                <Timeline project={project} />
              </div>
            </div>

            <div className="grid-col-8">
              <Tabs defaultSelectedId={getSelectedTabId(hash)}>
                <TabPanel id="documents" tab="Documents">
                  <button
                    type="button"
                    className="usa-button"
                    onClick={() => {
                      history.push(`${pathname}/upload`);
                    }}
                  >
                    Upload a document
                  </button>

                  <DocumentTable
                    project={project}
                    updateProject={updateProject}
                  />
                </TabPanel>
                <TabPanel id="notes" tab="Notes">
                  <form>
                    <label className="usa-label" htmlFor="input-type-textarea">
                      Add note
                    </label>
                    <textarea
                      className="usa-textarea"
                      id="input-type-textarea"
                      name="input-type-textarea"
                      value={noteContent}
                      onChange={e => {
                        setNoteContent(e.target.value);
                      }}
                      style={{ height: '100px' }}
                    />

                    <Button
                      className="margin-top-2"
                      type="button"
                      onClick={() => {
                        if (noteContent.trim().length > 0) {
                          project.notes.push({
                            id: Math.round(Math.random() * 10000000),
                            content: noteContent,
                            createdAt: DateTime.local(),
                            authorName: 'Aaron Allen'
                          });
                          project.banner = `Note added to ${project.name} project page.`;
                          project.lastUpdatedAt = DateTime.local();
                          updateProject(project);
                          setNoteContent('');
                        }
                      }}
                    >
                      Add Note
                    </Button>
                  </form>

                  <ol
                    className="note-list"
                    aria-label={`This is a list of all notes on ${project.name}.`}
                  >
                    {project.notes
                      .sort(
                        (a, b) =>
                          a.createdAt.toSeconds() - b.createdAt.toSeconds()
                      )
                      .map((activity: Note) => {
                        return (
                          <li className="" key={activity.id}>
                            <p className="margin-top-0 margin-bottom-1 text-pre-wrap">
                              {activity.content}
                            </p>
                            <span className="text-base-dark font-body-2xs">
                              by {activity.authorName}
                              <span aria-hidden="true">{' | '}</span>
                              {activity.createdAt.toFormat('LLLL d y')}
                            </span>
                            <hr aria-hidden="true" />
                          </li>
                        );
                      })}
                  </ol>
                </TabPanel>

                <TabPanel id="details" tab="Details and Past Requests">
                  <div
                    className="grid-container grid-gap"
                    style={{ padding: 0 }}
                  >
                    <div className="grid-row">
                      <div className="grid-col-6">
                        <dl className="detail-list">
                          <dt>Submitted date</dt>
                          <dd>{project.submissionDate.toFormat('LLLL d y')}</dd>
                          <dt>Business owner</dt>
                          <dd>
                            {project.businessOwner.name},{' '}
                            {project.businessOwner.component}
                          </dd>
                          <dt>Lifecycle ID</dt>
                          <dd>{project.lifecycleID}</dd>

                          <dt>Point of contact</dt>
                          <dd>
                            {project.pointOfContact.name}
                            <br />
                            <button
                              type="button"
                              className="usa-button usa-button--unstyled"
                              onClick={() => {
                                history.push('/508/v2/not-implemented');
                              }}
                            >
                              Update
                            </button>
                          </dd>
                        </dl>
                      </div>
                      <div className="grid-col-6">
                        <strong>Past Requests</strong>
                        <ol className="past-requests">
                          {project.pastRequests.map(request => (
                            <li>
                              <a href="/508/v2/not-implemented">
                                {request.name} {request.release}
                              </a>
                              <br />
                              Last tested on{' '}
                              {request.lastTestDate.toFormat('LLLL d y')}
                            </li>
                          ))}
                        </ol>
                      </div>
                    </div>
                  </div>
                </TabPanel>
              </Tabs>
            </div>
          </div>
        </div>
      </main>
    </>
  );
};

function stepContent(
  step: RequestStep,
  status:
    | RequestStepStatus
    | { status: ProgressStatus.NotCompleted; date: null }
) {
  if (
    [
      RequestStep.TestScheduled,
      RequestStep.ValidationTestingScheduled
    ].includes(step) &&
    status.date
  ) {
    return <>Test date: {status.date.toFormat('LLL d y')}</>;
  }
  if (step === RequestStep.RemediationInProgress && status.date) {
    return <>Start date: {status.date.toFormat('LLL d y')}</>;
  }
  return null;
}

export default ProjectPage;
