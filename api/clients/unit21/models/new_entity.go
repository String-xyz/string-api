package unit21

type NewEntity struct {
	GeneralData       *Entity        `json:"general_data"`
	UserData          *User          `json:"user_data"`
	DocumentData      *Document      `json:"document_data,omitempty"` //no document data collected?
	InstrumentIds     []string       `json:"instrument_ids"`          //will these be created first?
	CommunicationData *Communication `json:"communication_data"`
	DigitalData       *DigitalInfo   `json:"digital_data"`
	LocationData      *Location      `json:"location_data,omitempty"` //no location data collected?
	CustomData        *Custom        `json:"custom_data"`
	WorkflowOptions   *Options       `json:"options,omitempty"`
}

type Entity struct {
	EntityId      string    `json:"entity_id"`
	EntityType    string    `json:"entity_type"` //employee or business - says user in the spreadsheet?
	EntitySubType string    `json:"entity_subtype"`
	Status        string    `json:"status"`
	RegisteredAt  int       `json:"registered_at"` //date in seconds since 1/1/1970
	Parents       []*Parent `json:"parents,omitempty"`
	Tags          []string  `json:"tags"` //list of format: keyString:valueString
}

type Parent struct {
	EntityId   string `json:"entity_id"`
	EntityType string `json:"entity_type"` //employee or business
}

type User struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
	BirthDay   int    `json:"day_of_birth,omitempty"`
	BirthMonth int    `json:"month_of_birth,omitempty"`
	BirthYear  int    `json:"year_of_birth,omitempty"`
	SSN        string `json:"ssn,omitempty"`
}

type Document struct {
	DocumentId   string `json:"document_id"`
	DocumentType string `json:"document_type"`
	State        string `json:"state"`       //CA or California
	Country      string `json:"country"`     //https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	IssuedAt     int    `json:"issued_at"`   //date in seconds since 1/1/1970
	ExpiresAt    int    `json:"expires_at"`  //date in seconds since 1/1/1970
	VerifiedOn   int    `json:"verified_on"` //date in seconds since 1/1/1970
}

type Communication struct {
	Emails []string `json:"email_addresses"`
	Phones []string `json:"phone_numbers"` //E.164 format +12125551395 ( '[+][country code][area code][local phone number]' https://en.wikipedia.org/wiki/E.164
}

type DigitalInfo struct {
	IpAddresses        []string `json:"ip_addresses"` //ipv4 or ipv6
	ClientFingerprints []string `json:"client_fingerprints"`
}

type Location struct {
	AddressType    string `json:"type"`            //max 24 characters
	BuildingNumber string `json:"building_number"` //max 24 characters
	UnitNumber     string `json:"unit_number"`     //max 24 characters
	Street         string `json:"street_name"`     //max 128 characters
	City           string `json:"city"`            //max 128 characters
	State          string `json:"state"`           //CA or California
	PostalCode     string `json:"postal_code"`
	Country        string `json:"country"`     //short code or long name
	VerifiedOn     int    `json:"verified_on"` //date in seconds since 1/1/1970
}

type Custom struct {
	//more can be added to this as needed
	Platform string `json:"platform"` //where does this come from? do you want it called partnerName instead?
	Tier     int    `json:"tier"`     //is this an int?
}

type Options struct {
	IdentityVerifications *Identity `json:"identity_verifications"`
	ResolveGeoIp          bool      `json:"resolve_geoip"`      //default true
	MergeCustomData       bool      `json:"merge_custom_data"`  //default false https://docs.unit21.ai/docs/how-data-merges-on-updates#custom-data-merge-strategy
	UpsertOnConflict      bool      `json:"upsert_on_conflict"` //default true BUT we should change to false
}

type Identity struct {
	WorkflowId       string `json:"workflow_id"`
	RunVerifications bool   `json:"run_verifications"`       //default false do we want this?
	Synchronous      bool   `json:"synchronous_response"`    //default false do we want this? can take 2 mins
	FullResponse     bool   `json:"include_full_response"`   //default false only true if Synchronous = true
	SocureAPIKey     string `json:"socure_override_api_key"` // what is this?
}

type EntityResponse struct {
	Ignored           bool   `json:"ignored,omitempty"`
	EntityId          string `json:"entity_id"`
	PreviouslyExisted bool   `json:"previusly_existed"`
	Unit21Id          string `json:"unit21_id"`
}
