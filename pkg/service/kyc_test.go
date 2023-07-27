package service

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	env "github.com/String-xyz/go-lib/v2/config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/stretchr/testify/assert"
)

// define a test case struct
type kycCase struct {
	CaseName  string
	User      model.User
	Identity  model.Identity
	AssetType string
	Cost      float64
	Met       bool
}

func getTestCases() (cases []kycCase) {
	assetTypes := []string{"NFT", "TOKEN", "NFT_AND_TOKEN"}
	assetCosts := []float64{0, 999.99, 1000.00, 4999.99, 5000.00, 100000.00}
	kycLevels := []int{0, 1, 2, 3}
	now := time.Now()
	verifications := [][]*time.Time{
		{nil, nil, nil, nil},
		{&now, nil, nil, nil},
		{&now, &now, nil, nil},
		{&now, &now, &now, nil},
		{&now, &now, &now, &now},
	}
	baseUser := model.User{
		Id:         uuid.NewString(),
		CreatedAt:  now,
		UpdatedAt:  now,
		Type:       "User",
		Status:     "Onboarded",
		Tags:       nil,
		FirstName:  "Test",
		MiddleName: "A",
		LastName:   "User",
		Email:      "FakeUser123@nomail.com",
	}

	for _, assetType := range assetTypes {
		for _, assetCost := range assetCosts {
			for _, kycLevel := range kycLevels {
				for i, verification := range verifications {
					trueLevel := 0
					if i == 4 {
						trueLevel = 2
					} else if i >= 1 {
						trueLevel = 1
					}
					met := (assetType == "NFT" && assetCost < 1000.00 && trueLevel >= 1) || (assetCost < 5000.00 && trueLevel >= 2)
					testCase := kycCase{
						CaseName: fmt.Sprintf("AssetType: %s, AssetCost: %f, KYCLevel: %d, TrueLevel: %v", assetType, assetCost, kycLevel, trueLevel),
						User:     baseUser,
						Identity: model.Identity{
							Id:               uuid.NewString(),
							Level:            kycLevel,
							AccountId:        "",
							UserId:           "",
							CreatedAt:        now,
							UpdatedAt:        now,
							DeletedAt:        nil,
							EmailVerified:    verification[0],
							PhoneVerified:    verification[1],
							SelfieVerified:   verification[2],
							DocumentVerified: verification[3],
						},
						AssetType: assetType,
						Cost:      assetCost,
						Met:       met,
					}
					cases = append(cases, testCase)
				}
			}
		}
	}

	return cases
}

func setup(t *testing.T) (kyc KYC, ctx context.Context, mock sqlmock.Sqlmock, db *sql.DB) {
	env.LoadEnv(&config.Var, "../../.env")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	repos := repository.Repositories{
		Auth:        repository.NewAuth(nil, sqlxDB),
		Apikey:      repository.NewApikey(sqlxDB),
		User:        repository.NewUser(sqlxDB),
		Contact:     repository.NewContact(sqlxDB),
		Contract:    repository.NewContract(sqlxDB),
		Instrument:  repository.NewInstrument(sqlxDB),
		Device:      repository.NewDevice(sqlxDB),
		Asset:       repository.NewAsset(sqlxDB),
		Network:     repository.NewNetwork(sqlxDB),
		Transaction: repository.NewTransaction(sqlxDB),
		TxLeg:       repository.NewTxLeg(sqlxDB),
		Location:    repository.NewLocation(sqlxDB),
		Platform:    repository.NewPlatform(sqlxDB),
		Identity:    repository.NewIdentity(sqlxDB),
	}
	kyc = NewKYC(repos)
	ctx = context.Background()
	return kyc, ctx, mock, db
}

func TestMeetsRequirements(t *testing.T) {
	kyc, ctx, mock, db := setup(t)
	defer db.Close()
	testCases := getTestCases()
	for _, tc := range testCases {
		mockedIdentityRow := sqlmock.NewRows([]string{"id", "level", "account_id", "user_id", "email_verified", "phone_verified", "selfie_verified", "document_verified"}).
			AddRow(tc.Identity.Id, tc.Identity.Level, tc.Identity.AccountId, tc.Identity.UserId, tc.Identity.EmailVerified, tc.Identity.PhoneVerified, tc.Identity.SelfieVerified, tc.Identity.DocumentVerified)
		mock.ExpectQuery("SELECT * FROM identity WHERE user_id=$1").WithArgs(tc.User.Id).WillReturnRows(mockedIdentityRow)

		fmt.Printf("Running test for: %v\n", tc.CaseName)
		met, err := kyc.MeetsRequirements(ctx, tc.User.Id, tc.AssetType, tc.Cost)
		assert.NoError(t, err)
		assert.Equal(t, tc.Met, met)
	}
}
