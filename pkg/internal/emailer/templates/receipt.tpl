<a href='https://www.string.xyz'>
<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a> 
<header>Your {{.params.ReceiptType}} Details</header> 
<br>Dear {{.params.CustomerName}},
<br>Thank you for using String.  Here is your transaction receipt: 
<br>Transaction Date: {{.params.TransactionDate}}
<br>String Payment ID: {{.params.StringPaymentId}}

<br>Transaction ID: <a href={{.params.TransactionExplorer}}>{{.params.TransactionId}}</a>
<br>Destination Wallet: <a href='{{.params.DestinationExplorer}}'>{{.params.DestinationAddress}}</a>
<br>Payment Descriptor: {{.params.PaymentDescriptor}}
<br>Payment Method: {{.params.PaymentMethod}}
<br>Platform: {{.params.Platform}}
<br>Item Ordered: {{.params.ItemOrdered}}
<br>Token ID: {{.params.TokenId}}
<br>Subtotal: {{.params.Subtotal}}
<br>Network Fee: {{.params.NetworkFee}}
<br>Processing Fee: {{.params.ProcessingFee}}
<br>Total Charge: {{.params.Total}}

<br>The transaction will appear on your card statement as {{.params.PaymentDescriptor}}
<br>All sales are final.  Please see our <a href='https://www.string.xyz/terms-of-service'>Terms of Service</a> 
<br>Please reference your String Payment ID {{.params.StringPaymentId}}
<br><br>Service powered by String 
<br>String XYZ LLC | 490 43rd St, #86, Oakland CA 94609. | NMLS ID: 2400614 
<br>Please visit us at string.xyz.  Should you need to reach us, please contact us at <a href='mailto:support@string.xyz'>support@string.xyz</a>. 
<br><br>Consumer Fraud Warning 
<br>If you feel you have been the victim of a scam you can contact the FTC at 1-877-FTC-HELP (382-4357) 
<br>or online at <a href='www.ftc.gov'>www.ftc.gov</a> (link is external); or the Consumer Financial Protection Bureau (CFPB) at 1-855-411-CFPB (2372) 
<br>or online at <a href='www.consumerfinance.gov'>www.consumerfinance.gov</a>