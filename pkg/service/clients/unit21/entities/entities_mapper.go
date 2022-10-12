package service

type StringData struct { //temporary until the input data is established
	id           string
	tags         map[string]string
	partnerName  string
	status       string
	createdAt    int
	firstName    string
	middleName   string
	lastName     string
	emails       []string
	phones       []string
	ipAddresses  []string
	fingerprints []string
}

func MapStringDataToEntity(data StringData) *NewEntity {
	var userTagArr []string
	var userTags map[string]string //assumes the values are strings in the jsonb - to confirm
	if data.tags != nil {
		for key, value := range data.tags {
			userTagArr = append(userTagArr, key+":"+value)
		}
	}
	userTagArr = append(userTagArr, "platform:"+data.partnerName)

	for key, value := range userTags {
		userTagArr = append(userTagArr, key+":"+value)
	}

	// https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go
	jsonBody := &NewEntity{
		GeneralData: &Entity{
			EntityId:     data.id,
			EntityType:   "user",
			Status:       data.status,
			RegisteredAt: data.createdAt, //probably need to convert to seconds
			Tags:         userTagArr,     // convert from jsonb into array of key:value string pairs
		},
		UserData: &User{
			FirstName:  data.firstName,
			MiddleName: data.middleName,
			LastName:   data.lastName,
		},
		CommunicationData: &Communication{
			Emails: data.emails, //might need to be converted to []string
			Phones: data.phones, //might need to be converted to []string
		},
		DigitalData: &DigitalInfo{
			IpAddresses:        data.ipAddresses,  //might need to be converted to []string
			ClientFingerprints: data.fingerprints, //schema doesn't have a fingerprint, might need to be convered to []string
		},
		CustomData: &Custom{
			Platform: data.partnerName, //partnerName is not in the schema
		},
		// add WorkflowOptions if not default?
	}

	return jsonBody
}
