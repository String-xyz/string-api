package data

var AuhorizationApprovedJSON = `{
  "id": "evt_az5sblvku4ge3dwpztvyizgcau",
  "type": "authorization_approved",
  "version": "2.0.0",
  "created_on": "2018-04-10T08:13:14Z",
  "data": {
    "card_id": "crd_fa6psq242dcd6fdn5gifcq1491",
    "transaction_id": "trx_y3oqhf46pyzuxjbcn2giaqnb44",
    "transaction_type": "purchase",
    "transmission_date_time": "2018-04-10T08:12:14Z",
    "local_transaction_date_time": "2018-04-10T08:12:14",
    "authorization_type": "final_authorization",
    "transaction_amount": 1000,
    "transaction_currency": "GBP",
    "billing_amount": 900,
    "billing_currency": "EUR",
    "billing_conversion_rate": 0.925643,
    "ecb_conversion_rate": 0.9,
    "requested_exemption_type": "merchant_initiated_transaction",
    "digital_card_information": {
      "digital_card_id": "dcr_fa6psq242dcd6fdn5gifcq1491",
      "device_type": "phone",
      "wallet_type": "apple_pay"
    },
    "merchant": {
      "merchant_id": "59889",
      "name": "Carrefour",
      "city": "Paris",
      "state": "",
      "country": "FRA",
      "category_code": "5021"
    },
    "transit_information": {
        "transaction_type": "authorized_aggregated_split_clearing",
        "transportation_mode": "urban_bus"
    },
    "point_of_sale_transaction_date": "2023-04-25"
  }
}`

var AuhorizationDeclinedJSON = `
{
  "id": "evt_az5sblvku4ge3dwpztvyizgcau",
  "type": "authorization_declined",
  "version": "2.0.0",
  "created_on": "2018-04-10T08:13:14Z",
  "data": {
    "card_id": "crd_fa6psq242dcd6fdn5gifcq1491",
    "transaction_id": "trx_y3oqhf46pyzuxjbcn2giaqnb44",
    "transaction_type": "purchase",
    "transmission_date_time": "2018-04-10T08:12:14Z",
    "local_transaction_date_time": "2018-04-10T08:12:14",
    "authorization_type": "final_authorization",
    "transaction_amount": 1000,
    "transaction_currency": "GBP",
    "billing_amount": 900,
    "billing_currency": "EUR",
    "billing_conversion_rate": 0.925643,
    "ecb_conversion_rate": 0.9,
    "requested_exemption_type": "merchant_initiated_transaction",
    "digital_card_information": {
      "digital_card_id": "dcr_fa6psq242dcd6fdn5gifcq1491",
      "device_type": "phone",
      "wallet_type": "apple_pay"
    },
    "merchant": {
      "merchant_id": "59889",
      "name": "Carrefour",
      "city": "Paris",
      "state": "",
      "country": "FRA",
      "category_code": "5021"
    },
    "decline_reason": "insufficient_funds",
    "transit_information": {
        "transaction_type": "authorized_aggregated_split_clearing",
        "transportation_mode": "urban_bus"
    },
    "point_of_sale_transaction_date": "2023-04-25"
  }
}
`
var PaymentApprovedJSON = `
{
  "id": "evt_jyvkfwne6gnenlfm5yuyosi4wa",
  "type": "payment_approved",
  "version": "1.0.20",
  "created_on": "2022-11-14T12:18:54.3211949Z",
  "data": {
    "id": "pay_m3ncbzghvosu7b2j77hqln4agu",
    "action_id": "act_dglxv4ixum3ezijibwc4zaz7ce",
    "reference": "test payment",
    "amount": 2,
    "auth_code": "FPB9A0",
    "currency": "GBP",
    "payment_type": "Regular",
    "processed_on": "2022-11-14T12:18:52.4877348Z",
    "processing": {
      "acquirer_transaction_id": "123438482318443321234",
      "retrieval_reference_number": "123412701234"
    },
    "response_code": "10000",
    "response_summary": "Approved",
    "risk": {
      "flagged": false
    },
    "3ds": {
      "version": "2.2.0",
      "challenged": true,
      "challenge_indicator": "no_challenge_requested",
      "exemption": "none",
      "eci": "05",
      "cavv": "ABEBASVDkQBBASDCgmMYdQAAAAA=",
      "xid": "e6473b22-08ac-492d-a587-2aba9b80df6f",
      "downgraded": false,
      "enrolled": "Y",
      "authentication_response": "Y",
      "flow_type": "challenged"
    },
    "scheme_id": "482318443327456",
    "source": {
      "id": "src_ps6ukw5ifuietcr6wthijpzsvy",
      "type": "card",
      "billing_address": {},
      "expiry_month": 1,
      "expiry_year": 2029,
      "scheme": "VISA",
      "last_4": "1234",
      "fingerprint": "71580b426f1d190d29087ff265d8f48df1ad34ede41c27cbff9d23c1a14d1776",
      "bin": "123456",
      "card_type": "DEBIT",
      "card_category": "CONSUMER",
      "issuer": "Test Bank",
      "issuer_country": "GB",
      "product_id": "I",
      "product_type": "Visa Infinite",
      "avs_check": "I"
    },
    "balances": {
      "total_authorized": 2,
      "total_voided": 0,
      "available_to_void": 2,
      "total_captured": 0,
      "available_to_capture": 2,
      "total_refunded": 0,
      "available_to_refund": 0
    },
    "event_links": {
      "payment": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu",
      "payment_actions": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu/actions",
      "capture": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu/captures",
      "void": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu/voids"
    }
  },
  "_links": {
    "self": {
      "href": "https://api.checkout.com/workflows/events/evt_jyvkfwne6gnenlfm5yuyosi4wa"
    },
    "subject": {
      "href": "https://api.checkout.com/workflows/events/subject/pay_m3ncbzghvosu7b2j77hqln4agu"
    },
    "payment": {
      "href": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu"
    },
    "payment_actions": {
      "href": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu/actions"
    },
    "capture": {
      "href": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu/captures"
    },
    "void": {
      "href": "https://api.checkout.com/payments/pay_m3ncbzghvosu7b2j77hqln4agu/voids"
    }
  }
}`

var PaymentCapturedJSON = `{
  "id": "evt_6aznipgxbuaure3qen5qbzyswy",
  "type": "payment_captured",
  "version": "1.0.1",
  "created_on": "2019-06-07T08:25:22Z",
  "data": {
    "action_id": "act_gse7gcrhleuedmzhq25n3mhweq",
    "response_code": "10000",
    "response_summary": "Approved",
    "amount": 10000,
    "balances": {
      "total_authorized": 10000,
      "total_voided": 0,
      "available_to_void": 0,
      "total_captured": 10000,
      "available_to_capture": 0,
      "total_refunded": 0,
      "available_to_refund": 10000
    },
    "metadata": {
      "coupon_code": "NY2018",
      "partner_id": 123989
    },
    "processing": {
      "acquirer_transaction_id": "8137549557",
      "acquirer_reference_number": "000220552364"
    },
    "id": "pay_waji5li3mqtetnaor77xmow4bq",
    "currency": "EUR",
    "processed_on": "2019-06-07T08:25:22Z",
    "reference": "ORD-5023-4E89"
  },
  "_links": {
    "self": {
      "href": "https://api.checkout.com/events/evt_6aznipgxbuaure3qen5qbzyswy"
    },
    "subject": {
      "href": "https://api.checkout.com/workflows/events/subject/pay_jlfj2ful7z3u5lbykhy5lzezvm"
    }
  }
}`

var PersonAccountJSON = `
{
  "type": "event",
  "id": "evt_Ej9ZbvZGjn11CyudXiHqV7tN",
  "attributes": {
    "name": "account.created",
    "created-at": "2023-07-25T19:59:14.691Z",
    "redacted-at": null,
    "payload": {
      "data": {
        "type": "account",
        "id": "act_ze7eJgwEHbkx1iiUiu42G2ey",
        "attributes": {
          "reference-id": null,
          "created-at": "2023-07-25T19:59:14.000Z",
          "updated-at": "2023-07-25T19:59:14.000Z",
          "redacted-at": null,
          "fields": {
            "name": {
              "type": "hash",
              "value": {
                "first": {
                  "type": "string",
                  "value": "Mister"
                },
                "middle": {
                  "type": "string",
                  "value": null
                },
                "last": {
                  "type": "string",
                  "value": "Tester"
                }
              }
            },
            "address": {
              "type": "hash",
              "value": {
                "street-1": {
                  "type": "string",
                  "value": null
                },
                "street-2": {
                  "type": "string",
                  "value": null
                },
                "subdivision": {
                  "type": "string",
                  "value": null
                },
                "city": {
                  "type": "string",
                  "value": null
                },
                "postal-code": {
                  "type": "string",
                  "value": null
                },
                "country-code": {
                  "type": "string",
                  "value": null
                }
              }
            },
            "identification-numbers": {
              "type": "array",
              "value": []
            },
            "birthdate": {
              "type": "date",
              "value": null
            },
            "phone-number": {
              "type": "string",
              "value": null
            },
            "email-address": {
              "type": "string",
              "value": null
            },
            "selfie-photo": {
              "type": "file",
              "value": null
            }
          },
          "name-first": "Mister",
          "name-middle": null,
          "name-last": "Tester",
          "phone-number": null,
          "email-address": null,
          "address-street-1": null,
          "address-street-2": null,
          "address-city": null,
          "address-subdivision": null,
          "address-postal-code": null,
          "country-code": null,
          "birthdate": null,
          "social-security-number": null,
          "tags": [],
          "identification-numbers": {}
        }
      }
    }
  }
}
`
var PersonInquiryJSON = `
{
  "type": "event",
  "id": "evt_gEmTS7n2t3hHe4m2UYykLAym",
  "attributes": {
    "name": "inquiry.created",
    "created-at": "2023-07-25T19:32:25.582Z",
    "redacted-at": null,
    "payload": {
      "data": {
        "type": "inquiry",
        "id": "inq_bo5P7Ea1grrZc68Rg3mpFi2K",
        "attributes": {
          "status": "created",
          "reference-id": null,
          "note": null,
          "behaviors": {
            "request-spoof-attempts": null,
            "user-agent-spoof-attempts": null,
            "distraction-events": null,
            "hesitation-baseline": null,
            "hesitation-count": null,
            "hesitation-time": null,
            "shortcut-copies": null,
            "shortcut-pastes": null,
            "autofill-cancels": null,
            "autofill-starts": null,
            "devtools-open": null,
            "completion-time": null,
            "hesitation-percentage": null,
            "behavior-threat-level": null
          },
          "tags": [],
          "creator": "API",
          "reviewer-comment": null,
          "created-at": "2023-07-25T19:32:25.000Z",
          "started-at": null,
          "completed-at": null,
          "failed-at": null,
          "marked-for-review-at": null,
          "decisioned-at": null,
          "expired-at": null,
          "redacted-at": null,
          "previous-step-name": null,
          "next-step-name": "start_biometric_80e902_start",
          "name-first": null,
          "name-middle": null,
          "name-last": null,
          "birthdate": null,
          "address-street-1": null,
          "address-street-2": null,
          "address-city": null,
          "address-subdivision": null,
          "address-subdivision-abbr": null,
          "address-postal-code": null,
          "address-postal-code-abbr": null,
          "social-security-number": null,
          "identification-number": null,
          "email-address": null,
          "phone-number": null,
          "fields": {
            "phone-number": {
              "type": "string",
              "value": null
            },
            "selected-country-code": {
              "type": "string",
              "value": "US"
            },
            "current-government-id": {
              "type": "government_id",
              "value": null
            },
            "selected-id-class": {
              "type": "string",
              "value": null
            },
            "address-street-1": {
              "type": "string",
              "value": null
            },
            "address-street-2": {
              "type": "string",
              "value": null
            },
            "address-city": {
              "type": "string",
              "value": null
            },
            "address-subdivision": {
              "type": "string",
              "value": null
            },
            "address-postal-code": {
              "type": "string",
              "value": null
            },
            "address-country-code": {
              "type": "string",
              "value": null
            },
            "birthdate": {
              "type": "date",
              "value": null
            },
            "email-address": {
              "type": "string",
              "value": null
            },
            "identification-class": {
              "type": "string",
              "value": null
            },
            "identification-number": {
              "type": "string",
              "value": null
            },
            "name-first": {
              "type": "string",
              "value": null
            },
            "name-middle": {
              "type": "string",
              "value": null
            },
            "name-last": {
              "type": "string",
              "value": null
            },
            "current-selfie": {
              "type": "selfie",
              "value": null
            }
          }
        },
        "relationships": {
          "account": {
            "data": {
              "type": "account",
              "id": "act_ndJNqdhWNi44S4Twf4bqzod1"
            }
          },
          "template": {
            "data": null
          },
          "inquiry-template": {
            "data": {
              "type": "inquiry-template",
              "id": "itmpl_z2so7W2bCFHELp2dhxqqQjGy"
            }
          },
          "inquiry-template-version": {
            "data": {
              "type": "inquiry-template-version",
              "id": "itmplv_CQbgqFNNgwrGe6jAB496VXpd"
            }
          },
          "reviewer": {
            "data": null
          },
          "reports": {
            "data": []
          },
          "verifications": {
            "data": []
          },
          "sessions": {
            "data": []
          },
          "documents": {
            "data": []
          },
          "selfies": {
            "data": []
          }
        }
      }
    }
  }
}
`
