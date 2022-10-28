package unit21

// //////////////////////////////////////////////////////////////////
// Entity
type u21Entity struct {
	GeneralData       *entityGeneral       `json:"general_data"`
	UserData          *entityPersonal      `json:"user_data,omitempty"`
	CommunicationData *entityCommunication `json:"communication_data,omitempty"`
	DigitalData       *entityDigitalData   `json:"digital_data,omitempty"`
	CustomData        *entityCustomData    `json:"custom_data,omitempty"`
	WorkflowOptions   *options             `json:"options,omitempty"`
}

type entityGeneral struct {
	EntityId      string   `json:"entity_id"`
	EntityType    string   `json:"entity_type"` //employee or business - says user in the spreadsheet?
	EntitySubType string   `json:"entity_subtype"`
	Status        string   `json:"status,omitempty"`
	RegisteredAt  int      `json:"registered_at"`  //date in seconds since 1/1/1970
	Tags          []string `json:"tags,omitempty"` //list of format: keyString:valueString
}

type entityPersonal struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
	BirthDay   int    `json:"day_of_birth,omitempty"`
	BirthMonth int    `json:"month_of_birth,omitempty"`
	BirthYear  int    `json:"year_of_birth,omitempty"`
	SSN        string `json:"ssn,omitempty"`
}

type entityCommunication struct {
	Emails []string `json:"email_addresses,omitempty"`
	Phones []string `json:"phone_numbers,omitempty"` //E.164 format +12125551395 ( '[+][country code][area code][local phone number]' https://en.wikipedia.org/wiki/E.164
}

type entityDigitalData struct {
	IpAddresses        []string `json:"ip_addresses,omitempty"` //ipv4 or ipv6
	ClientFingerprints []string `json:"client_fingerprints,omitempty"`
}

type entityCustomData struct {
	//more can be added to this as needed
	Platforms []string `json:"platforms,omitempty"` //where does this come from? do you want it called partnerName instead?
}

type options struct {
	ResolveGeoIp     bool `json:"resolve_geoip"`      //default true
	MergeCustomData  bool `json:"merge_custom_data"`  //default false https://docs.unit21.ai/docs/how-data-merges-on-updates#custom-data-merge-strategy
	UpsertOnConflict bool `json:"upsert_on_conflict"` //default true BUT should we change to false? don't need update endpoint if allowed
}

type createEntityResponse struct {
	Ignored           bool   `json:"ignored,omitempty"`
	EntityId          string `json:"entity_id"`
	PreviouslyExisted bool   `json:"previously_existed"`
	Unit21Id          string `json:"unit21_id"`
}

type updateEntityResponse struct {
	EntityId string `json:"entity_id"`
	Unit21Id string `json:"unit21_id"`
}

// //////////////////////////////////////////////////////////////////
// Instrument
type u21Instrument struct {
	InstrumentId       string                  `json:"instrument_id"`
	InstrumentType     string                  `json:"instrument_type"`
	InstrumentSubtype  string                  `json:"instrument_subtype"`
	Source             string                  `json:"source"`
	Status             string                  `json:"status,omitempty"`
	RegisteredAt       int                     `json:"registered_at"`
	ParentInstrumentId string                  `json:"parent_instrument_id,omitempty"`
	Entities           []instrumentEntity      `json:"entities,omitempty"`
	CustomData         any                     `json:"custom_data,omitempty"`
	DigitalData        *instrumentDigitalData  `json:"digital_data,omitempty"`
	LocationData       *instrumentLocationData `json:"location_data,omitempty"`
	Tags               []string                `json:"tags,omitempty"`
	Options            *options                `json:"options,omitempty"`
}

type instrumentEntity struct {
	EntityId       string `json:"entity_id"`
	RelationshipId string `json:"relationship_id"`
	EntityType     string `json:"entity_type"`
}

type instrumentDigitalData struct {
	IpAddresses []string `json:"ip_addresses,omitempty"`
}

type instrumentLocationData struct {
	Type           string `json:"type"`
	BuildingNumber string `json:"building_number"`
	UnitNumber     string `json:"unit_number"`
	StreetName     string `json:"street_name"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postal_code"`
	Country        string `json:"country"`
	VerifiedOn     int    `json:"verified_on"`
}

type instrumentCustomData struct {
	//more can be added to this as needed
	None any `json:"none"`
}

type createInstrumentResponse struct {
	Ignored           bool   `json:"ignored,omitempty"`
	InstrumentId      string `json:"instrument_id"`
	PreviouslyExisted bool   `json:"previously_existed"`
	Unit21Id          string `json:"unit21_id"`
}

type updateInstrumentResponse struct {
	InstrumentId string `json:"instrument_id"`
	Unit21Id     string `json:"unit21_id"`
}

// //////////////////////////////////////////////////////////////////
// Event

type u21event struct {
	GeneralData     *eventGeneral    `json:"general_data"`
	TransactionData *transactionData `json:"transaction_data"`
}

type eventGeneral struct {
	EventId      string       `json:"event_id"`
	EventType    string       `json:"event_type"`
	EventTime    int          `json:"event_time"`
	EventSubtype string       `json:"event_subtype"`
	Status       string       `json:"status"`
	Parents      *eventParent `json:"parents"`
	Tags         []string     `json:"tags,omitempty"`
}

type transactionData struct {
	Amount         float32 `json:"amount"`
	SentAmount     float32 `json:"sent_amount"`
	SentCurrency   string  `json:"sent_currency"`
	SenderEntityId string  `json:"sender_entity_id"`
}

type eventParent struct {
	EventId   string `json:"event_id"`
	EventType string `json:"event_type"`
}
