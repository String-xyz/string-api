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
	Type  string `json:"type"`
	Value string `json:"value"`
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

type PhoneURL struct {
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

type InquiryField struct {
	AddressStreet1 StringValue `json:"addressStreet1"`
	AddressStreet2 StringValue `json:"addressStreet2"`
}

type InquiryAttributes struct {
	Attribute
	ReferenceId      *string       `json:"referenceId"`
	Behaviors        Behavior      `json:"behaviors"`
	Notes            *string       `json:"notes"`
	Tags             []interface{} `json:"tags"`
	PreviousStepName string        `json:"previousStepName"`
	NextStepName     string        `json:"nextStepName"`
	Fields           InquiryField  `json:"fields"`
}

type VerificationAttributes struct {
	Attribute
	CountryCode    *string    `json:"countryCode"`
	LeftPhotoUrl   *string    `json:"leftPhotoUrl"`
	RightPhotoUrl  *string    `json:"rightPhotoUrl"`
	CenterPhotoUrl *string    `json:"centerPhotoUrl"`
	PhotoUrls      []PhoneURL `json:"photoUrls"`
	Checks         []Check    `json:"checks"`
	CaptureMethod  string     `json:"captureMethod"`
}

type IncludeAttributes struct {
	VerificationAttributes
	SelfiePhoto             *string   `json:"selfiePhoto"`
	SelfiePhotoUrl          *string   `json:"selfiePhotoUrl"`
	FrontPhotoUrl           *PhotoURL `json:"frontPhotoUrl"`
	BackPhotoUrl            *PhoneURL `json:"backPhotoUrl"`
	VideoUrl                *string   `json:"videoUrl"`
	IdClass                 string    `json:"idClass"`
	CaptureMethod           string    `json:"captureMethod"`
	EntityConfidenceScore   int       `json:"entityConfidenceScore"`
	EntityConfidenceReasons []string  `json:"entityConfidenceReasons"`
	NameFirst               string    `json:"nameFirst"`
	NameMiddle              *string   `json:"nameMiddle"`
	NameLast                string    `json:"nameLast"`
	NameSuffix              *string   `json:"nameSuffix"`
	Birthdate               string    `json:"birthdate"`
	AddressStreet1          string    `json:"addressStreet1"`
	AddressStreet2          *string   `json:"addressStreet2"`
	AddressCity             string    `json:"addressCity"`
	AddressSubdivision      string    `json:"addressSubdivision"`
	AddressPostalCode       string    `json:"addressPostalCode"`
	IssuingAuthority        string    `json:"issuingAuthority"`
	IssuingSubdivision      string    `json:"issuingSubdivision"`
	Nationality             *string   `json:"nationality"`
	DocumentNumber          *string   `json:"documentNumber"`
	VisaStatus              *string   `json:"visaStatus"`
	IssueDate               string    `json:"issueDate"`
	ExpirationDate          string    `json:"expirationDate"`
	Designations            *string   `json:"designations"`
	Birthplace              *string   `json:"birthplace"`
	Endorsements            *string   `json:"endorsements"`
	Height                  *string   `json:"height"`
	Sex                     string    `json:"sex"`
	Restrictions            *string   `json:"restrictions"`
	VehicleClass            *string   `json:"vehicleClass"`
	IdentificationNumber    string    `json:"identificationNumber"`
}

type AccountField struct {
	Name                  HashValue  `json:"name"`
	Address               HashValue  `json:"address"`
	IdentificationNumbers ArrayValue `json:"identificationNumbers"`
	Birthdate             Value      `json:"birthdate"`
	PhoneNumber           Value      `json:"phoneNumber"`
	EmailAddress          Value      `json:"emailAddress"`
	SelfiePhoto           Value      `json:"selfiePhoto"`
}

type Include struct {
	Id            string            `json:"id"`
	Type          string            `json:"type"`
	Atrributes    IncludeAttributes `json:"attributes"`
	Relationships Relationships     `json:"relationships"`
}

type Inquiry struct {
	Id             string         `json:"id"`
	Status         string         `json:"status"`
	CreatedAt      string         `json:"createdAt"`
	LastUpdatedAt  string         `json:"lastUpdatedAt"`
	CompletedSteps CompletedSteps `json:"completedSteps"`
}

type InquiryPayload struct {
	AccountId string `json:"accountId"`
	Template  string `json:"template"`
}

type CompletedSteps struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type IdentityVerification struct {
	Id             string         `json:"id"`
	Status         string         `json:"status"`
	CreatedAt      string         `json:"createdAt"`
	LastUpdatedAt  string         `json:"lastUpdatedAt"`
	CompletedSteps CompletedSteps `json:"completedSteps"`
}

type IdentityVerificationPayload struct {
	AccountId string `json:"accountId"`
	Template  string `json:"template"`
}

type Verification struct {
	Id            string                 `json:"id"`
	Attributes    VerificationAttributes `json:"attributes"`
	Relationships Relationship           `json:"relationships"`
}

type Account struct {
	Type       string    `json:"type"`
	Id         string    `json:"id"`
	Attributes Attribute `json:"attributes"`
}

type AccountResponse struct {
	Data Account `json:"data"`
}

type ListAccountResponse struct {
	Data  []Account `json:"data"`
	Links Link      `json:"links"`
}

type ListInquiryResponse struct {
	Data  []Inquiry `json:"data"`
	Links Link      `json:"links"`
}
