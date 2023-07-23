package persona

import "time"

type IdValue struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type HashValue struct {
	Type  string           `json:"type"`
	Value map[string]Value `json:"value"`
}

type StringValue struct {
	Type  string  `json:"type"`
	Value *string `json:"value"`
}

type ArrayValue struct {
	Type  string      `json:"type"`
	Value []HashValue `json:"value"`
}

type Value struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Link struct {
	Prev *string `json:"prev"`
	Next *string `json:"next"`
}

type PhotoURL struct {
	Page          *string  `json:"page"`
	Url           *string  `json:"url"`
	fileName      *string  `json:"fileName"`
	NormalizedUrl *string  `json:"normalizedUrl"`
	OriginalUrls  []string `json:"originalUrls"`
	ByteSize      int      `json:"byteSize"`
}

type Check struct {
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Reason      []interface{} `json:"reason"`
	Requirement string        `json:"requirement"`
	Metadata    interface{}   `json:"metadata"`
}

type RelationshipId struct {
	Data *IdValue `json:"data"`
}

type RelationshipIds struct {
	Data []IdValue `json:"data"`
}

type Relationships struct {
	Account                     *RelationshipId  `json:"account"`
	Inquity                     *RelationshipId  `json:"inquiry"`
	Template                    *RelationshipId  `json:"template"`
	InquityTemplate             *RelationshipId  `json:"inquiryTemplate"`
	InquityTemplateVersion      *RelationshipId  `json:"inquiryTemplateVersion"`
	VerificationTemplate        *RelationshipId  `json:"verificationTemplate"`
	VerificationTemplateVersion *RelationshipId  `json:"verificationTemplateVersion"`
	Verifications               *RelationshipIds `json:"verifications"`
	Sessions                    *RelationshipIds `json:"sessions"`
	Documents                   *RelationshipIds `json:"documents"`
	DocumentFiles               *RelationshipIds `json:"documentFiles"`
	Selfies                     *RelationshipIds `json:"selfies"`
}

type Behavior struct {
	RequestSpoofAttempts   int     `json:"requestSpoofAttempts"`
	UserAgentSpoofAttempts int     `json:"userAgentSpoofAttempts"`
	DistractionEvents      int     `json:"distractionEvents"`
	HesitationBaseline     int     `json:"hesitationBaseline"`
	HesitationCount        int     `json:"hesitationCount"`
	HesitationTime         int     `json:"hesitationTime"`
	ShortcutCopies         int     `json:"shortcutCopies"`
	ShortcutPastes         int     `json:"shortcutPastes"`
	AutofillCancels        int     `json:"autofillCancels"`
	AutofillStarts         int     `json:"autofillStarts"`
	DevtoolsOpen           bool    `json:"devtoolsOpen"`
	CompletionTime         float64 `json:"completionTime"`
	HesitationPercentage   float64 `json:"hesitationPercentage"`
	BehaviorThreatLevel    string  `json:"behaviorThreatLevel"`
}

type AccountFields struct {
	Name                  HashValue   `json:"name"`
	Address               HashValue   `json:"address"`
	IdentificationNumbers ArrayValue  `json:"identificationNumbers"`
	Birthdate             Value       `json:"birthdate"`
	PhoneNumber           StringValue `json:"phoneNumber"`
	EmailAddress          StringValue `json:"emailAddress"`
	SelfiePhoto           Value       `json:"selfiePhoto"`
}

type InquiryFields struct {
	AddressStreet1 StringValue `json:"addressStreet1"`
	AddressStreet2 StringValue `json:"addressStreet2"`
}

type Attribute struct {
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"createdAt"`
	StartedAt       *time.Time `json:"startedAt"`
	FailedAt        *time.Time `json:"failedAt"`
	DecesionedAt    *time.Time `json:"decesionedAt"`
	MarkForReviewAt *time.Time `json:"markForReviewAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	RedactedAt      *time.Time `json:"redactedAt"`
	SubmittedAt     *time.Time `json:"submittedAt"`
	CompletedAt     *time.Time `json:"completedAt"`
	ExpiredAt       *time.Time `json:"expiredAt"`
}

type AccountAttributes struct {
	Attribute
	Fields AccountFields `json:"fields"`
}

type CommonFields struct {
	Attribute
	// City of residence address. Not all international addresses use this attribute.
	AddresCity string `json:"addressCity,omitempty"`
	// Street name of residence address.
	AddressStreet1 string `json:"addressStreet1,omitempty"`
	// Extension of residence address, usually apartment or suite number.
	AddressStreet2 string `json:"addressStreet2,omitempty"`
	// State or subdivision of residence address. In the US,
	// this should be the unabbreviated name. Not all international addresses use this attribute.
	AddressSubdivision string `json:"addressSubdivision,omitempty"`
	// Postal code of residence address. Not all international addresses use this attribute.
	AddressPostalCode string `json:"addressPostalCode,omitempty"`
	// Birthdate, must be in the format "YYYY-MM-DD".
	Birthdate string `json:"birthdate,omitempty"`
	// ISO 3166-1 alpha 2 country code of the government ID to be verified. This is generally their country of residence as well.
	CountryCode string `json:"countryCode,omitempty"`

	EmailAddress string `json:"emailAddress,omitempty"`
	// Given or first name.
	NameFirst string `json:"nameFirst,omitempty"`
	// Family or last name.
	NameLast string `json:"nameLast,omitempty"`

	NameMiddle string `json:"nameMiddle,omitempty"`

	PhoneNumber string `json:"phoneNumber,omitempty"`

	SocialSecurityNumber string `json:"socialSecurityNumber,omitempty"`
}

type CommonAttributes struct {
	SelfiePhoto             *string  `json:"selfiePhoto"`
	SelfiePhotoUrl          *string  `json:"selfiePhotoUrl"`
	FrontPhotoUrl           *string  `json:"frontPhotoUrl"`
	BackPhotoUrl            *string  `json:"backPhotoUrl"`
	VideoUrl                *string  `json:"videoUrl"`
	IdClass                 string   `json:"idClass"`
	CaptureMethod           string   `json:"captureMethod"`
	EntityConfidenceScore   float64  `json:"entityConfidenceScore"`
	EntityConfidenceReasons []string `json:"entityConfidenceReasons"`
	NameFirst               string   `json:"nameFirst"`
	NameMiddle              *string  `json:"nameMiddle"`
	NameLast                string   `json:"nameLast"`
	NameSuffix              *string  `json:"nameSuffix"`
	Birthdate               string   `json:"birthdate"`
	AddressStreet1          string   `json:"addressStreet1"`
	AddressStreet2          *string  `json:"addressStreet2"`
	AddressCity             string   `json:"addressCity"`
	AddressSubdivision      string   `json:"addressSubdivision"`
	AddressPostalCode       string   `json:"addressPostalCode"`
	IssuingAuthority        string   `json:"issuingAuthority"`
	IssuingSubdivision      string   `json:"issuingSubdivision"`
	Nationality             *string  `json:"nationality"`
	DocumentNumber          *string  `json:"documentNumber"`
	VisaStatus              *string  `json:"visaStatus"`
	IssueDate               string   `json:"issueDate"`
	ExpirationDate          string   `json:"expirationDate"`
	Designations            *string  `json:"designations"`
	Birthplace              *string  `json:"birthplace"`
	Endorsements            *string  `json:"endorsements"`
	Height                  *string  `json:"height"`
	Sex                     string   `json:"sex"`
	Restrictions            *string  `json:"restrictions"`
	VehicleClass            *string  `json:"vehicleClass"`
	IdentificationNumber    string   `json:"identificationNumber"`
}

type InquiryCreationAttributes struct {
	AccountId                string `json:"accountId"`
	CountryCode              string `json:"countryCode"`
	InquityTemplateId        string `json:"inquiryTemplateId"`
	InquityTemplateVersionId string `json:"inquiryTemplateVersionId"`
	// Template ID for flow requirements (use this field if your template ID starts with tmpl_).
	// You must pass in either template-id OR inquiry-template-id OR inquiry-template-version-id
	TemplateId        string `json:"templateId"`
	TemplateVersionId string `json:"templateVersionId"`
	// for styling
	ThemeId string `json:"themeId"`

	Fields CommonFields `json:"fields"`
}

type InquiryAttributes struct {
	Attribute
	ReferenceId      *string       `json:"referenceId"`
	Behaviors        Behavior      `json:"behaviors"`
	Notes            *string       `json:"notes"`
	Tags             []interface{} `json:"tags"`
	PreviousStepName string        `json:"previousStepName"`
	NextStepName     string        `json:"nextStepName"`
	Fields           InquiryFields `json:"fields"`
}

type CompletedSteps struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type VerificationAttributes struct {
	Attribute
	CommonAttributes
	CountryCode    *string    `json:"countryCode"`
	LeftPhotoUrl   *string    `json:"leftPhotoUrl"`
	RightPhotoUrl  *string    `json:"rightPhotoUrl"`
	CenterPhotoUrl *string    `json:"centerPhotoUrl"`
	PhotoUrls      []PhotoURL `json:"photoUrls"`
	Checks         []Check    `json:"checks"`
	CaptureMethod  string     `json:"captureMethod"`
}

type Included struct {
	Id            string                 `json:"id"`
	Type          string                 `json:"type"`
	Atrributes    VerificationAttributes `json:"attributes"`
	Relationships Relationships          `json:"relationships"`
}
