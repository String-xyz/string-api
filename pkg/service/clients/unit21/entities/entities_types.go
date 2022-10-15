package service

type NewEntity struct {
	GeneralData       *Entity        `json:"general_data"`
	UserData          *User          `json:"user_data,omitempty"`
	CommunicationData *Communication `json:"communication_data,omitempty"`
	DigitalData       *DigitalInfo   `json:"digital_data,omitempty"`
	CustomData        *Custom        `json:"custom_data,omitempty"`
	WorkflowOptions   *Options       `json:"options,omitempty"`
}

type Entity struct {
	EntityId      string   `json:"entity_id"`
	EntityType    string   `json:"entity_type"` //employee or business - says user in the spreadsheet?
	EntitySubType string   `json:"entity_subtype"`
	Status        string   `json:"status,omitempty"`
	RegisteredAt  int      `json:"registered_at"`  //date in seconds since 1/1/1970
	Tags          []string `json:"tags,omitempty"` //list of format: keyString:valueString
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

type Communication struct {
	Emails []string `json:"email_addresses,omitempty"`
	Phones []string `json:"phone_numbers,omitempty"` //E.164 format +12125551395 ( '[+][country code][area code][local phone number]' https://en.wikipedia.org/wiki/E.164
}

type DigitalInfo struct {
	IpAddresses        []string `json:"ip_addresses,omitempty"` //ipv4 or ipv6
	ClientFingerprints []string `json:"client_fingerprints,omitempty"`
}

type Custom struct {
	//more can be added to this as needed
	Platforms []string `json:"platforms,omitempty"` //where does this come from? do you want it called partnerName instead?
}

type Options struct {
	ResolveGeoIp     bool `json:"resolve_geoip"`      //default true
	MergeCustomData  bool `json:"merge_custom_data"`  //default false https://docs.unit21.ai/docs/how-data-merges-on-updates#custom-data-merge-strategy
	UpsertOnConflict bool `json:"upsert_on_conflict"` //default true BUT should we change to false? don't need update endpoint if allowed
}

type EntityResponse struct {
	Ignored           bool   `json:"ignored,omitempty"`
	EntityId          string `json:"entity_id"`
	PreviouslyExisted bool   `json:"previusly_existed"`
	Unit21Id          string `json:"unit21_id"`
}
