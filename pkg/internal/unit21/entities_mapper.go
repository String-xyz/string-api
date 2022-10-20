package unit21

import "github.com/String-xyz/string-api/pkg/model"

func MapUserToEntity(user model.User, communication communication) *u21entity {
	var userTagArr []string
	if user.Tags != nil {
		for key, value := range user.Tags {
			userTagArr = append(userTagArr, key+":"+value)
		}
	}

	jsonBody := &u21entity{
		GeneralData: &general{
			EntityId:     user.ID,
			EntityType:   "user",
			Status:       user.Status,
			RegisteredAt: int(user.CreatedAt.Unix()),
			Tags:         userTagArr, // convert from jsonb into array of key:value string pairs
		},
		UserData: &personal{
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			LastName:   user.LastName,
		},
		CommunicationData: &communication,
		DigitalData: &digitalInfo{
			IpAddresses:        nil, //data.ipAddresses,  //might need to be converted to []string
			ClientFingerprints: nil, //data.fingerprints, //schema doesn't have a fingerprint, might need to be convered to []string
		},
		CustomData: nil, //&custom{
		// 	Platform: nil,//data.partnerName,
		// },
		// add WorkflowOptions if not default?
	}

	return jsonBody
}
