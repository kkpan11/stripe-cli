package fixtures

import (
	"context"
	"embed"
	"fmt"
	"sort"

	"github.com/spf13/afero"

	"github.com/stripe/stripe-cli/pkg/stripe"
)

//go:embed triggers/*
var triggers embed.FS

// Events is a mapping of pre-built trigger events and the corresponding json file
var Events = map[string]string{
	"account.application.deauthorized":          "triggers/account.application.deauthorized.json",
	"account.updated":                           "triggers/account.updated.json",
	"balance.available":                         "triggers/balance.available.json",
	"billing_portal.configuration.created":      "triggers/billing_portal.configuration.created.json",
	"billing_portal.configuration.updated":      "triggers/billing_portal.configuration.updated.json",
	"billing_portal.session.created":            "triggers/billing_portal.session.created.json",
	"cash_balance.funds_available":              "triggers/cash_balance.funds_available.json",
	"charge.captured":                           "triggers/charge.captured.json",
	"charge.dispute.closed":                     "triggers/charge.dispute.closed.json",
	"charge.dispute.created":                    "triggers/charge.dispute.created.json",
	"charge.dispute.updated":                    "triggers/charge.dispute.updated.json",
	"charge.failed":                             "triggers/charge.failed.json",
	"charge.refunded":                           "triggers/charge.refunded.json",
	"charge.refund.updated":                     "triggers/charge.refund.updated.json",
	"charge.succeeded":                          "triggers/charge.succeeded.json",
	"checkout.session.async_payment_failed":     "triggers/checkout.session.async_payment_failed.json",
	"checkout.session.async_payment_succeeded":  "triggers/checkout.session.async_payment_succeeded.json",
	"checkout.session.completed":                "triggers/checkout.session.completed.json",
	"checkout.session.expired":                  "triggers/checkout.session.expired.json",
	"coupon.created":                            "triggers/coupon.created.json",
	"coupon.deleted":                            "triggers/coupon.deleted.json",
	"coupon.updated":                            "triggers/coupon.updated.json",
	"credit_note.created":                       "triggers/credit_note.created.json",
	"credit_note.updated":                       "triggers/credit_note.updated.json",
	"credit_note.voided":                        "triggers/credit_note.voided.json",
	"customer_cash_balance_transaction.created": "triggers/customer_cash_balance_transaction.created.json",
	"customer.discount.created":                 "triggers/customer.discount.created.json",
	"customer.discount.deleted":                 "triggers/customer.discount.deleted.json",
	"customer.discount.updated":                 "triggers/customer.discount.updated.json",
	"customer.created":                          "triggers/customer.created.json",
	"customer.deleted":                          "triggers/customer.deleted.json",
	"customer.updated":                          "triggers/customer.updated.json",
	"customer.source.created":                   "triggers/customer.source.created.json",
	"customer.source.updated":                   "triggers/customer.source.updated.json",
	"customer.subscription.created":             "triggers/customer.subscription.created.json",
	"customer.subscription.deleted":             "triggers/customer.subscription.deleted.json",
	"customer.subscription.updated":             "triggers/customer.subscription.updated.json",
	"customer.subscription.paused":              "triggers/customer.subscription.paused.json",
	"customer.subscription.trial_will_end":      "triggers/customer.subscription.trial_will_end.json",
	"identity.verification_session.canceled":    "triggers/identity.verification_session.canceled.json",
	"identity.verification_session.created":     "triggers/identity.verification_session.created.json",
	"identity.verification_session.redacted":    "triggers/identity.verification_session.redacted.json",
	"invoice.created":                           "triggers/invoice.created.json",
	"invoice.finalized":                         "triggers/invoice.finalized.json",
	"invoice.paid":                              "triggers/invoice.paid.json",
	"invoice.payment_action_required":           "triggers/invoice.payment_action_required.json",
	"invoice.payment_failed":                    "triggers/invoice.payment_failed.json",
	"invoice.payment_succeeded":                 "triggers/invoice.payment_succeeded.json",
	"invoice.updated":                           "triggers/invoice.updated.json",
	"invoice.deleted":                           "triggers/invoice.deleted.json",
	"invoice.voided":                            "triggers/invoice.voided.json",
	"invoice.sent":                              "triggers/invoice.sent.json",
	"invoice.marked_uncollectible":              "triggers/invoice.marked_uncollectible.json",
	"invoiceitem.created":                       "triggers/invoiceitem.created.json",
	"invoiceitem.deleted":                       "triggers/invoiceitem.deleted.json",
	"issuing_authorization.request":             "triggers/issuing_authorization.request.json",
	"issuing_authorization.request.eu":          "triggers/issuing_authorization.request.eu.json",
	"issuing_authorization.request.gb":          "triggers/issuing_authorization.request.gb.json",
	"issuing_card.created":                      "triggers/issuing_card.created.json",
	"issuing_card.created.eu":                   "triggers/issuing_card.created.eu.json",
	"issuing_card.created.gb":                   "triggers/issuing_card.created.gb.json",
	"issuing_cardholder.created":                "triggers/issuing_cardholder.created.json",
	"issuing_cardholder.created.eu":             "triggers/issuing_cardholder.created.eu.json",
	"issuing_cardholder.created.gb":             "triggers/issuing_cardholder.created.gb.json",
	"payment_intent.amount_capturable_updated":  "triggers/payment_intent.amount_capturable_updated.json",
	"payment_intent.created":                    "triggers/payment_intent.created.json",
	"payment_intent.payment_failed":             "triggers/payment_intent.payment_failed.json",
	"payment_intent.processing":                 "triggers/payment_intent.processing.json",
	"payment_intent.succeeded":                  "triggers/payment_intent.succeeded.json",
	"payment_intent.canceled":                   "triggers/payment_intent.canceled.json",
	"payment_link.created":                      "triggers/payment_link.created.json",
	"payment_link.updated":                      "triggers/payment_link.updated.json",
	"payment_intent.partially_funded":           "triggers/payment_intent.partially_funded.json",
	"payment_intent.requires_action":            "triggers/payment_intent.requires_action.json",
	"payment_method.attached":                   "triggers/payment_method.attached.json",
	"payment_method.detached":                   "triggers/payment_method.detached.json",
	"payment_method.updated":                    "triggers/payment_method.updated.json",
	"payout.created":                            "triggers/payout.created.json",
	"payout.updated":                            "triggers/payout.updated.json",
	"plan.created":                              "triggers/plan.created.json",
	"plan.deleted":                              "triggers/plan.deleted.json",
	"plan.updated":                              "triggers/plan.updated.json",
	"price.created":                             "triggers/price.created.json",
	"price.updated":                             "triggers/price.updated.json",
	"product.created":                           "triggers/product.created.json",
	"product.deleted":                           "triggers/product.deleted.json",
	"product.updated":                           "triggers/product.updated.json",
	"promotion_code.created":                    "triggers/promotion_code.created.json",
	"promotion_code.updated":                    "triggers/promotion_code.updated.json",
	"reporting.report_run.succeeded":            "triggers/reporting.report_run.succeeded.json",
	"setup_intent.canceled":                     "triggers/setup_intent.canceled.json",
	"setup_intent.created":                      "triggers/setup_intent.created.json",
	"setup_intent.setup_failed":                 "triggers/setup_intent.setup_failed.json",
	"setup_intent.succeeded":                    "triggers/setup_intent.succeeded.json",
	"setup_intent.requires_action":              "triggers/setup_intent.requires_action.json",
	"subscription_schedule.canceled":            "triggers/subscription_schedule.canceled.json",
	"subscription_schedule.created":             "triggers/subscription_schedule.created.json",
	"subscription_schedule.released":            "triggers/subscription_schedule.released.json",
	"subscription_schedule.updated":             "triggers/subscription_schedule.updated.json",
	"subscription.payment_succeeded":            "triggers/subscription.payment_succeeded.json",
	"subscription.payment_failed":               "triggers/subscription.payment_failed.json",
	"quote.created":                             "triggers/quote.created.json",
	"quote.canceled":                            "triggers/quote.canceled.json",
	"quote.finalized":                           "triggers/quote.finalized.json",
	"quote.accepted":                            "triggers/quote.accepted.json",
	"tax_rate.created":                          "triggers/tax_rate.created.json",
	"tax_rate.updated":                          "triggers/tax_rate.updated.json",
	"transfer.created":                          "triggers/transfer.created.json",
	"transfer.reversed":                         "triggers/transfer.reversed.json",
	"transfer.updated":                          "triggers/transfer.updated.json",
	"v1.billing.meter.error_report_triggered":   "triggers/v1.billing.meter.error_report_triggered.json",
	"v1.billing.meter.no_meter_found":           "triggers/v1.billing.meter.no_meter_found.json",
}

// BuildFromFixtureFile creates a new fixture struct for a file
func BuildFromFixtureFile(fs afero.Fs, apiKey, stripeAccount, apiBaseURL, jsonFile string, skip, override, add, remove []string, edit bool) (*Fixture, error) {
	fixture, err := NewFixtureFromFile(
		fs,
		apiKey,
		stripeAccount,
		apiBaseURL,
		jsonFile,
		skip,
		override,
		add,
		remove,
		edit,
	)
	if err != nil {
		return nil, err
	}

	return fixture, nil
}

// BuildFromFixtureString creates a new fixture from a string
func BuildFromFixtureString(fs afero.Fs, apiKey, stripeAccount, apiBaseURL, raw string) (*Fixture, error) {
	fixture, err := NewFixtureFromRawString(fs, apiKey, stripeAccount, apiBaseURL, raw)
	if err != nil {
		return nil, err
	}
	return fixture, nil
}

// EventList prints out a padded list of supported trigger events for printing the help file
func EventList() string {
	var eventList string
	for _, event := range EventNames() {
		eventList += fmt.Sprintf("  %s\n", event)
	}

	return eventList
}

// EventNames returns an array of all the event names
func EventNames() []string {
	names := []string{}
	for name := range Events {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Trigger triggers a Stripe event.
func Trigger(ctx context.Context, event string, stripeAccount string, baseURL string, apiKey string, skip, override, add, remove []string, raw string, apiVersion string, edit bool) ([]string, error) {
	var fixture *Fixture
	var err error
	fs := afero.NewOsFs()

	// send event triggered
	telemetryClient := stripe.GetTelemetryClient(ctx)
	if telemetryClient != nil {
		go telemetryClient.SendEvent(ctx, "Triggered Event", event)
	}

	if len(raw) == 0 {
		if file, ok := Events[event]; ok {
			fixture, err = BuildFromFixtureFile(fs, apiKey, stripeAccount, baseURL, file, skip, override, add, remove, edit)
			if err != nil {
				return nil, err
			}
		} else {
			exists, _ := afero.Exists(fs, event)
			if !exists {
				return nil, fmt.Errorf("%s", fmt.Sprintf("The event `%s` is not supported by Stripe CLI. To trigger unsupported events, use the Stripe API or Dashboard to perform actions that lead to the event you want to trigger (for example, create a Customer to generate a `customer.created` event). You can also create a custom fixture: https://docs.stripe.com/cli/fixtures", event))
			}

			fixture, err = BuildFromFixtureFile(fs, apiKey, stripeAccount, baseURL, event, skip, override, add, remove, edit)
			if err != nil {
				return nil, err
			}
		}
	} else {
		fixture, err = BuildFromFixtureString(fs, apiKey, stripeAccount, baseURL, raw)
		if err != nil {
			return nil, err
		}
	}

	requestNames, err := fixture.Execute(ctx, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("%s", fmt.Sprintf("Trigger failed: %s\n", err))
	}

	return requestNames, nil
}

func reverseMap() map[string]string {
	reversed := make(map[string]string)
	for name, file := range Events {
		reversed[file] = name
	}

	return reversed
}
