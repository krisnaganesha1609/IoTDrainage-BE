package requests

type RegisterTokenRequest struct {
	DeviceID      string `json:"device_id" validate:"required"`
	AppInstanceID string `json:"app_instance_id" validate:"required"`
	FCMToken      string `json:"fcm_token" validate:"required"`
}
