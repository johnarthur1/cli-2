package v7_test

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/command/commandfakes"
	. "code.cloudfoundry.org/cli/command/v7"
	"code.cloudfoundry.org/cli/command/v7/v7fakes"
	"code.cloudfoundry.org/cli/resources"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("service command", func() {
	const (
		serviceInstanceName = "fake-service-instance-name"
		serviceInstanceGUID = "fake-service-instance-guid"
		spaceName           = "fake-space-name"
		spaceGUID           = "fake-space-guid"
		orgName             = "fake-org-name"
		username            = "fake-user-name"
	)
	var (
		cmd             ServiceCommand
		testUI          *ui.UI
		fakeConfig      *commandfakes.FakeConfig
		fakeSharedActor *commandfakes.FakeSharedActor
		fakeActor       *v7fakes.FakeActor
		executeErr      error
	)

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeSharedActor = new(commandfakes.FakeSharedActor)
		fakeActor = new(v7fakes.FakeActor)

		cmd = ServiceCommand{
			BaseCommand: BaseCommand{
				UI:          testUI,
				Config:      fakeConfig,
				SharedActor: fakeSharedActor,
				Actor:       fakeActor,
			},
		}

		fakeConfig.CurrentUserReturns(configv3.User{Name: username}, nil)

		fakeConfig.TargetedSpaceReturns(configv3.Space{
			GUID: spaceGUID,
			Name: spaceName,
		})

		fakeConfig.TargetedOrganizationReturns(configv3.Organization{
			Name: orgName,
		})

		fakeActor.GetServiceInstanceDetailsReturns(
			v7action.ServiceInstanceWithRelationships{
				ServiceInstance: resources.ServiceInstance{
					GUID: serviceInstanceGUID,
					Name: serviceInstanceName,
				},
				SharedStatus: v7action.SharedStatus{
					IsShared: false,
				},
			},
			v7action.Warnings{"warning one", "warning two"},
			nil,
		)

		setPositionalFlags(&cmd, serviceInstanceName)
	})

	It("checks the user is logged in, and targeting an org and space", func() {
		Expect(fakeSharedActor.CheckTargetCallCount()).To(Equal(1))
		orgChecked, spaceChecked := fakeSharedActor.CheckTargetArgsForCall(0)
		Expect(orgChecked).To(BeTrue())
		Expect(spaceChecked).To(BeTrue())
	})

	When("the --guid flag is specified", func() {
		BeforeEach(func() {
			setFlag(&cmd, "--guid")
		})

		It("looks up the service instance and prints the GUID and warnings", func() {
			Expect(executeErr).NotTo(HaveOccurred())

			Expect(fakeActor.GetServiceInstanceDetailsCallCount()).To(Equal(1))
			actualName, actualSpaceGUID := fakeActor.GetServiceInstanceDetailsArgsForCall(0)
			Expect(actualName).To(Equal(serviceInstanceName))
			Expect(actualSpaceGUID).To(Equal(spaceGUID))

			Expect(testUI.Out).To(Say(`^%s\n$`, serviceInstanceGUID))
			Expect(testUI.Err).To(SatisfyAll(
				Say("warning one"),
				Say("warning two"),
			))
		})
	})

	When("it is a user-provided service instance", func() {
		const (
			routeServiceURL = "https://route.com"
			syslogURL       = "https://syslog.com"
			tags            = "foo, bar"
		)

		BeforeEach(func() {
			fakeActor.GetServiceInstanceDetailsReturns(
				v7action.ServiceInstanceWithRelationships{
					ServiceInstance: resources.ServiceInstance{
						GUID:            serviceInstanceGUID,
						Name:            serviceInstanceName,
						Type:            resources.UserProvidedServiceInstance,
						SyslogDrainURL:  types.NewOptionalString(syslogURL),
						RouteServiceURL: types.NewOptionalString(routeServiceURL),
						Tags:            types.NewOptionalStringSlice(strings.Split(tags, ", ")...),
					},
				},
				v7action.Warnings{"warning one", "warning two"},
				nil,
			)
		})

		It("looks up the service instance and prints the details and warnings", func() {
			Expect(executeErr).NotTo(HaveOccurred())

			Expect(fakeActor.GetServiceInstanceDetailsCallCount()).To(Equal(1))
			actualName, actualSpaceGUID := fakeActor.GetServiceInstanceDetailsArgsForCall(0)
			Expect(actualName).To(Equal(serviceInstanceName))
			Expect(actualSpaceGUID).To(Equal(spaceGUID))

			Expect(testUI.Out).To(SatisfyAll(
				Say(`Showing info of service %s in org %s / space %s as %s...\n`, serviceInstanceName, orgName, spaceName, username),
				Say(`\n`),
				Say(`name:\s+%s\n`, serviceInstanceName),
				Say(`guid:\s+\S+\n`),
				Say(`type:\s+user-provided`),
				Say(`tags:\s+%s\n`, tags),
				Say(`route service url:\s+%s\n`, routeServiceURL),
				Say(`syslog drain url:\s+%s\n`, syslogURL),
			))

			Expect(testUI.Err).To(SatisfyAll(
				Say("warning one"),
				Say("warning two"),
			))
		})
	})

	When("it is a managed service instance", func() {
		const (
			dashboardURL               = "https://dashboard.com"
			tags                       = "foo, bar"
			servicePlanName            = "fake-service-plan-name"
			serviceOfferingName        = "fake-service-offering-name"
			serviceOfferingDescription = "an amazing service"
			serviceOfferingDocs        = "https://service.docs.com"
			serviceBrokerName          = "fake-service-broker-name"
			lastOperationType          = "create"
			lastOperationState         = "in progress"
			lastOperationDescription   = "doing amazing work"
			lastOperationStartTime     = "a second ago"
			lastOperationUpdatedTime   = "just now"
		)

		BeforeEach(func() {
			fakeActor.GetServiceInstanceDetailsReturns(
				v7action.ServiceInstanceWithRelationships{
					ServiceInstance: resources.ServiceInstance{
						GUID:         serviceInstanceGUID,
						Name:         serviceInstanceName,
						Type:         resources.ManagedServiceInstance,
						DashboardURL: types.NewOptionalString(dashboardURL),
						Tags:         types.NewOptionalStringSlice(strings.Split(tags, ", ")...),
						LastOperation: resources.LastOperation{
							Type:        lastOperationType,
							State:       lastOperationState,
							Description: lastOperationDescription,
							CreatedAt:   lastOperationStartTime,
							UpdatedAt:   lastOperationUpdatedTime,
						},
					},
					ServiceOffering: resources.ServiceOffering{
						Name:             serviceOfferingName,
						Description:      serviceOfferingDescription,
						DocumentationURL: serviceOfferingDocs,
					},
					ServicePlanName:   servicePlanName,
					ServiceBrokerName: serviceBrokerName,
					SharedStatus: v7action.SharedStatus{
						IsShared: true,
					},
				},
				v7action.Warnings{"warning one", "warning two"},
				nil,
			)
		})

		It("looks up the service instance and prints the details and warnings", func() {
			Expect(executeErr).NotTo(HaveOccurred())

			Expect(fakeActor.GetServiceInstanceDetailsCallCount()).To(Equal(1))
			actualName, actualSpaceGUID := fakeActor.GetServiceInstanceDetailsArgsForCall(0)
			Expect(actualName).To(Equal(serviceInstanceName))
			Expect(actualSpaceGUID).To(Equal(spaceGUID))

			Expect(testUI.Out).To(SatisfyAll(
				Say(`Showing info of service %s in org %s / space %s as %s...\n`, serviceInstanceName, orgName, spaceName, username),
				Say(`\n`),
				Say(`name:\s+%s\n`, serviceInstanceName),
				Say(`guid:\s+\S+\n`),
				Say(`type:\s+managed`),
				Say(`broker:\s+%s`, serviceBrokerName),
				Say(`offering:\s+%s`, serviceOfferingName),
				Say(`plan:\s+%s`, servicePlanName),
				Say(`tags:\s+%s\n`, tags),
				Say(`description:\s+%s\n`, serviceOfferingDescription),
				Say(`documentation:\s+%s\n`, serviceOfferingDocs),
				Say(`dashboard url:\s+%s\n`, dashboardURL),
				Say(`\n`),
				Say(`Showing status of last operation from service instance %s...\n`, serviceInstanceName),
				Say(`\n`),
				Say(`status:\s+%s %s\n`, lastOperationType, lastOperationState),
				Say(`message:\s+%s\n`, lastOperationDescription),
				Say(`started:\s+%s\n`, lastOperationStartTime),
				Say(`updated:\s+%s\n`, lastOperationUpdatedTime),
			))

			Expect(testUI.Err).To(SatisfyAll(
				Say("warning one"),
				Say("warning two"),
			))
		})

		Context("service instances sharing", func() {
			When("service instance is shared", func() {
				It("shows shared information", func() {
					Expect(testUI.Out).To(SatisfyAll(
						Say(`Sharing:`),
						Say(`This service instance is currently shared.`),
					))
				})
			})

			When("service is not shared", func() {
				BeforeEach(func() {
					fakeActor.GetServiceInstanceDetailsReturns(
						v7action.ServiceInstanceWithRelationships{
							ServiceInstance: resources.ServiceInstance{},
							SharedStatus: v7action.SharedStatus{
								IsShared: false,
							},
						},
						v7action.Warnings{},
						nil,
					)
				})

				It("displays that the service is not shared", func() {
					Expect(testUI.Out).To(SatisfyAll(
						Say(`Sharing:`),
						Say(`This service instance is not currently being shared.`),
					))
				})
			})

			When("the service instance sharing feature is disabled", func() {
				BeforeEach(func() {
					fakeActor.GetServiceInstanceDetailsReturns(
						v7action.ServiceInstanceWithRelationships{
							ServiceInstance: resources.ServiceInstance{},
							SharedStatus: v7action.SharedStatus{
								FeatureFlagIsDisabled: true,
							},
						},
						v7action.Warnings{},
						nil,
					)
				})

				It("displays that the sharing feature is disabled", func() {
					Expect(testUI.Out).To(SatisfyAll(
						Say(`Sharing:\n`),
						Say(`\n`),
						Say(`The "service_instance_sharing" feature flag is disabled for this Cloud Foundry platform.\n`),
						Say(`\n`),
					))
				})
			})

			When("the service instance sharing feature is enabled", func() {
				BeforeEach(func() {
					fakeActor.GetServiceInstanceDetailsReturns(
						v7action.ServiceInstanceWithRelationships{
							ServiceInstance: resources.ServiceInstance{},
							SharedStatus: v7action.SharedStatus{
								FeatureFlagIsDisabled: false,
							},
						},
						v7action.Warnings{},
						nil,
					)
				})

				It("does not display a warning", func() {
					Expect(testUI.Out).NotTo(
						Say(`The "service_instance_sharing" feature flag is disabled for this Cloud Foundry platform.`),
					)
				})
			})

			When("the offering does not allow service instance sharing", func() {
				BeforeEach(func() {
					fakeActor.GetServiceInstanceDetailsReturns(
						v7action.ServiceInstanceWithRelationships{
							ServiceInstance: resources.ServiceInstance{},
							SharedStatus: v7action.SharedStatus{
								OfferingDisablesSharing: true,
							},
						},
						v7action.Warnings{},
						nil,
					)
				})

				It("displays that the sharing feature is disabled", func() {
					Expect(testUI.Out).To(SatisfyAll(
						Say(`Sharing:\n`),
						Say(`\n`),
						Say(`Service instance sharing is disabled for this service offering.\n`),
						Say(`\n`),
					))
				})
			})

			When("the offering does allow service instance sharing", func() {
				BeforeEach(func() {
					fakeActor.GetServiceInstanceDetailsReturns(
						v7action.ServiceInstanceWithRelationships{
							ServiceInstance: resources.ServiceInstance{},
							SharedStatus: v7action.SharedStatus{
								OfferingDisablesSharing: false,
							},
						},
						v7action.Warnings{},
						nil,
					)
				})

				It("does not display a warning", func() {
					Expect(testUI.Out).NotTo(
						Say(`Service instance sharing is disabled for this service offering.`),
					)
				})
			})
		})

	})

	When("there is a problem looking up the service instance", func() {
		BeforeEach(func() {
			fakeActor.GetServiceInstanceDetailsReturns(
				v7action.ServiceInstanceWithRelationships{},
				v7action.Warnings{"warning one", "warning two"},
				errors.New("boom"),
			)
		})

		It("prints warnings and returns an error", func() {
			Expect(executeErr).To(MatchError("boom"))

			Expect(testUI.Out).NotTo(Say(`.`), "output not empty!")
			Expect(testUI.Err).To(SatisfyAll(
				Say("warning one"),
				Say("warning two"),
			))
		})
	})

	When("checking the target returns an error", func() {
		BeforeEach(func() {
			fakeSharedActor.CheckTargetReturns(errors.New("explode"))
		})

		It("returns the error", func() {
			Expect(executeErr).To(MatchError("explode"))
		})
	})
})