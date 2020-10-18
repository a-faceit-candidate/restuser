package restuser

// User describes the main item of the user service: a user.
type User struct {
	// ID is generated by the service when the user is created. It is a valid UUID.
	ID string `json:"id" example:"c3e11b46-109c-11eb-adc1-0242ac120002" format:"uuid"`
	// CreatedAt is set by the service when the user is created.
	// It's formatted as an RFC3339 timestamp.
	CreatedAt string `json:"created_at" example:"2006-01-02T15:04:05Z" format:"date-time"`
	// UpdatedAt is set by the service when the user is updated.
	// It's formatted as an RFC3339 timestamp. For a recently created user, it equals the CreatedAt field.
	UpdatedAt string `json:"updated_at" example:"2006-01-02T15:04:05Z" format:"date-time"`
	// Name is the name of the user.
	Name string `json:"name" example:"John Doe"`
	// Email is the email of the user.
	Email string `json:"email" example:"john@colega.eu" format:"email"`
	// Country is the country code of the user, in ISO 3166-1 alpha-2 formatted as lowercase two character string.
	// Country is not validated to exist when a user is created.
	Country string `json:"country" example:"es"`
}

// ErrorResponse is used to provide further details on non-successful responses.
type ErrorResponse struct {
	Message string `json:"message" example:"Something terrible happened."`
}
