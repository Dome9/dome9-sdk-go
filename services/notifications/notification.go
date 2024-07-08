﻿package notifications

import (
	"fmt"
	"net/http"
	"time"
)

const (
	RESTfulServicePathNotification = "notification"
)

// Constants
type AssessmentFindingOrigin int

const (
	ComplianceEngine AssessmentFindingOrigin = iota
	Magellan
	MagellanAwsGuardDuty        = 191
	Serverless                  = 2
	AwsInspector                = 50
	ServerlessSecurityAnalyzer  = 51
	ExternalFindingSource       = 100
	Qualys                      = 101
	Tenable                     = 102
	AwsGuardDuty                = 103
	KubernetesImageScanning     = 104
	KubernetesRuntimeAssurance  = 105
	ContainersRuntimeProtection = 106
	WorkloadChangeMonitoring    = 107
	ImageAssurance              = 7
	SourceCodeAssurance         = 8
	InfrastructureAsCode        = 9
	CIEM                        = 10
	Incident                    = 11
)

type NotificationTriggerType int

const (
	Report NotificationTriggerType = iota
	Single
	Scheduled
)

type NotificationOutputType int

const (
	Default NotificationOutputType = iota
	Detailed
	Summary
	FullCsv
	FullCsvZip
	ExecutivePlatform
	JsonFullEntity
	JsonSimpleEntity
	PlainText
	TemplateBased
	CustomOutputFormat
)

// Models
type BaseNotificationViewModel struct {
	Name                 string                               `json:"name" validate:"required"`
	Description          string                               `json:"description"`
	AlertsConsole        bool                                 `json:"alertsConsole" default:"true"`
	SendOnEachOccurrence bool                                 `json:"sendOnEachOccurrence"`
	Origin               string                               `json:"origin" validate:"required"`
	IntegrationSettings  NotificationIntegrationSettingsModel `json:"integrationSettingsModel" validate:"required"`
}

type NotificationIntegrationSettingsModel struct {
	ReportsIntegrationSettings            []ReportNotificationIntegrationSettings    `json:"reportsIntegrationSettings"`
	SingleNotificationIntegrationSettings []SingleNotificationIntegrationSettings    `json:"singleNotificationIntegrationSettings"`
	ScheduledIntegrationSettings          []ScheduledNotificationIntegrationSettings `json:"scheduledIntegrationSettings"`
}

type BaseNotificationIntegrationSettings struct {
	IntegrationId string                       `json:"integrationId" validate:"required"`
	OutputType    string                       `json:"outputType"`
	Filter        ComplianceNotificationFilter `json:"filter"`
}

type SingleNotificationIntegrationSettings struct {
	BaseNotificationIntegrationSettings
	Payload string `json:"payload"`
}

type ReportNotificationIntegrationSettings struct {
	BaseNotificationIntegrationSettings
}

type ScheduledNotificationIntegrationSettings struct {
	BaseNotificationIntegrationSettings
	CronExpression string `json:"cronExpression" validate:"required,cron"`
}

type ComplianceNotificationFilter struct{}

type PutNotificationViewModel struct {
	BaseNotificationViewModel
	Id string `json:"id" validate:"required"`
}

type PostNotificationViewModel struct {
	BaseNotificationViewModel
}

type ResponseNotificationViewModel struct {
	BaseNotificationViewModel
	Id        string    `json:"id" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

// APIs

func (service *Service) Create(body PostNotificationViewModel) (*ResponseNotificationViewModel, *http.Response, error) {
	v := new(ResponseNotificationViewModel)
	resp, err := service.Client.NewRequestDo("POST", RESTfulServicePathNotification, nil, body, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetAll() (*[]ResponseNotificationViewModel, *http.Response, error) {
	v := new([]ResponseNotificationViewModel)
	resp, err := service.Client.NewRequestDo("GET", RESTfulServicePathNotification, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetById(id string) (*ResponseNotificationViewModel, *http.Response, error) {
	if id == "" {
		return nil, nil, fmt.Errorf("id parameter must be passed")
	}

	v := new(ResponseNotificationViewModel)
	relativeURL := fmt.Sprintf("%s/%s", RESTfulServicePathNotification, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(name string) (*ResponseNotificationViewModel, *http.Response, error) {
	if name == "" {
		return nil, nil, fmt.Errorf("name parameter must be passed")
	}

	v := new(ResponseNotificationViewModel)
	relativeURL := fmt.Sprintf("%s?name=%s", RESTfulServicePathNotification, name)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(body PutNotificationViewModel) (*ResponseNotificationViewModel, *http.Response, error) {
	if body.Id == "" {
		return nil, nil, fmt.Errorf("id parameter must be passed")
	}

	v := new(ResponseNotificationViewModel)
	resp, err := service.Client.NewRequestDo("PUT", RESTfulServicePathNotification, nil, body, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Delete(id string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", RESTfulServicePathNotification, id)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
