package persona

import (
	"net/http"

	"github.com/String-xyz/go-lib/v2/common"
)

/*

To verify a set of inputs from an individual, a Verification object is created.
A Verification enables an Organization to answer “Is this person who they claim to be?” with a focus on verifying digital transactions.
The verification process is accomplished through the generation and processing of Checks against the information provided by the individual.

The collective process of mixing and matching verifications to achieve sufficient assurance that
an individual is indeed who they claim to be is often called identity proofing.
The goal of identity proofing is often tying digital identities to physical identities.

An Inquiry contains one or more verifications. The attributes available for any given verification depends on its type.
Each inquiry’s relationships field lists the IDs of all associated verifications.
To authenticate when fetching photo URLs, pass the same Authorization header.

Verifications change statuses as the individual progresses through the flow.
Check for the following statuses to monitor progress and find completed results.

* Initiated			- Verification has started, claimed information can now sent and saved to the server for verification
* Confirmed			- Verification has been confirmed. This is a status specific to PhoneNumber verifications where they have verified a confirmation code that was entered.
* Submitted			- Verification has been submitted, the claimed information is frozen and the server will process the verification
* Passed			- Verification has passed. The required checks have passed and the information is verified
* Requires Retry	- Verification requires a resubmission. The checks could not be fully processed due to issues with the submitted information
* Failed			- Verification has failed. Some or all of the required checks have failed and verification has failed

*/

type Verification struct {
	Id            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    VerificationAttributes `json:"attributes"`
	Relationships Relationships          `json:"relationships"`
}

func (v Verification) GetType() string {
	return v.Type
}

type VerificationResponse struct {
	Data Verification `json:"data"`
}

func (c *PersonaClient) GetVerificationById(id string) (*VerificationResponse, error) {
	verification := &VerificationResponse{}
	err := c.doRequest(http.MethodGet, "/v1/verifications/"+id, nil, verification)
	if err != nil {
		return nil, common.StringError(err, "failed to get verification by id")
	}

	return verification, nil
}
