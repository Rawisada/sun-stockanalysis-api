package realtime

import "sun-stockanalysis-api/internal/models"

type CompositeAlertNotifier struct {
	notifiers []AlertEventNotifier
}

func NewCompositeAlertNotifier(notifiers ...AlertEventNotifier) *CompositeAlertNotifier {
	filtered := make([]AlertEventNotifier, 0, len(notifiers))
	for _, notifier := range notifiers {
		if notifier != nil {
			filtered = append(filtered, notifier)
		}
	}
	return &CompositeAlertNotifier{notifiers: filtered}
}

func (n *CompositeAlertNotifier) Notify(event *models.AlertEvent, message string) {
	if n == nil {
		return
	}
	for _, notifier := range n.notifiers {
		notifier.Notify(event, message)
	}
}
