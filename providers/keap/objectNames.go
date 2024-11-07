package keap

import (
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/internal/datautils"
	"github.com/amp-labs/connectors/providers/keap/metadata"
)

const (
	objectNameAffiliates           = "affiliates"
	objectNameAppointments         = "appointments"
	objectNameAutomationCategories = "automationCategories"
	objectNameCompanies            = "companies"
	objectNameContacts             = "contacts"
	objectNameOrders               = "orders"
	objectNameSubscriptions        = "subscriptions"
	objectNameEmails               = "emails"
	objectNameFiles                = "files"
	objectNameNotes                = "notes"
	objectNameOpportunities        = "opportunities"
	objectNamePaymentMethodConfigs = "paymentMethodConfigs"
	objectNameProducts             = "products"
	objectNameHooks                = "hooks"
	objectNameTags                 = "tags"
	objectNameTagCategories        = "tag_categories"
	objectNameTasks                = "tasks"
	objectNameUsers                = "users"
)

// Supported object names can be found under schemas.json.
var supportedObjectsByRead = metadata.Schemas.ObjectNames() //nolint:gochecknoglobals

var supportedObjectsByCreate = map[common.ModuleID]datautils.StringSet{ //nolint:gochecknoglobals
	ModuleV1: datautils.NewSet(
		objectNameAffiliates,
		objectNameAppointments,
		objectNameCompanies,
		objectNameContacts,
		objectNameOrders,
		objectNameSubscriptions,
		objectNameEmails,
		objectNameFiles,
		objectNameNotes,
		objectNameOpportunities,
		objectNameProducts,
		objectNameHooks,
		objectNameTags,
		objectNameTagCategories,
		objectNameTasks,
		objectNameUsers,
	),
	ModuleV2: datautils.NewSet(
		objectNameAffiliates,
		objectNameAutomationCategories,
		objectNameCompanies,
		objectNameContacts,
		objectNameEmails,
		objectNamePaymentMethodConfigs,
		objectNameSubscriptions,
		objectNameTags,
		objectNameTagCategories,
	),
}

var supportedObjectsByUpdatePATCH = map[common.ModuleID]datautils.StringSet{ //nolint:gochecknoglobals
	ModuleV1: datautils.NewSet(
		objectNameAppointments,
		objectNameCompanies,
		objectNameContacts,
		objectNameNotes,
		objectNameOpportunities,
		objectNameProducts,
		objectNameTasks,
	),
	ModuleV2: datautils.NewSet(
		objectNameCompanies,
		objectNameContacts,
		objectNameTags,
		objectNameTagCategories,
	),
}

var supportedObjectsByUpdatePUT = map[common.ModuleID]datautils.StringSet{ //nolint:gochecknoglobals
	ModuleV1: datautils.NewSet(
		objectNameFiles,
		objectNameHooks,
	),
	ModuleV2: datautils.NewSet(
		objectNameAutomationCategories,
	),
}

var supportedObjectsByDelete = map[common.ModuleID]datautils.StringSet{ //nolint:gochecknoglobals
	ModuleV1: datautils.NewSet(
		objectNameAppointments,
		objectNameContacts,
		objectNameOrders,
		objectNameEmails,
		objectNameFiles,
		objectNameNotes,
		objectNameProducts,
		objectNameHooks,
		objectNameTasks,
	),
	ModuleV2: datautils.NewSet(
		objectNameAutomationCategories,
		objectNameCompanies,
		objectNameContacts,
		objectNameEmails,
		objectNameTagCategories,
	),
}

// ObjectNameToWritePath maps ObjectName to URL path used for Write operation.
//
// Some of the ignored endpoints:
// "/v1/account/profile" -- update single profile
// "/v1/emails/queue" -- send email to list of contacts
// "/v1/emails/unsync" -- un-sync a batch of email records
// "/v1/emails/sync" -- create a set of email records
// "/v2/businessProfile" -- update single profile.
var ObjectNameToWritePath = datautils.NewDefaultMap(map[string]string{ //nolint:gochecknoglobals
	objectNameTagCategories:        "tags/categories",
	objectNameAutomationCategories: "automationCategory", // API uses singular form. Others are consistently plural.
},
	func(objectName string) (jsonPath string) {
		return objectName
	},
)
