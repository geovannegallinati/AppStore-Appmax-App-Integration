package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
)

// TestWebhookEnvelope_AllRealEventNames_ParseWithoutError verifies that every
// real event envelope received from Appmax can be unmarshaled into WebhookEnvelopeRequest.
func TestWebhookEnvelope_AllRealEventNames_ParseWithoutError(t *testing.T) {
	envelopes := []string{
		`{"event":"CustomerCreated","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Laura","lastname":"Montenegro","email":"maximo33@gmail.com","telephone":"3526001145","postcode":"97120572","address_street":"Avenida Santiago","address_street_number":"55","address_street_complement":"Bloco C","address_street_district":"qui","address_city":"Santa Daniel","address_state":"SP","document_number":"57135108914","created_at":"2026-03-17 09:19:46","visited_url":"","uf":"SP","fullname":"Laura Montenegro","interested_bundle":[]}}`,
		`{"event":"CustomerInterested","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Josefina","lastname":"Fonseca","email":"thiago.beltrao@vale.com","telephone":"75997495154","postcode":"63502560","address_street":"Av. Felipe","address_street_number":"04","address_street_complement":"Apto 9","address_street_district":"non","address_city":"Porto Noel","address_state":"SP","document_number":"96948998003","created_at":"2026-03-17 09:24:25","visited_url":"","uf":"SP","fullname":"Josefina Fonseca","interested_bundle":[]}}`,
		`{"event":"CustomerContacted","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Pedro","lastname":"Godoi","email":"fernando.avila@hotmail.com","telephone":"22995701715","postcode":"63839964","address_street":"Travessa Paulina Vale","address_street_number":"7839","address_street_complement":"75 Andar","address_street_district":"inventore","address_city":"Santa Isabella do Leste","address_state":"RR","document_number":"24200165058","created_at":"2026-03-17 09:25:08","visited_url":"","uf":"RR","fullname":"Pedro Godoi","interested_bundle":[]}}`,
		`{"event":"OrderApproved","event_type":"","data":{"id":1,"customer_id":1,"total_products":265,"status":"aprovado","freight_value":2.48,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:13","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":267.48,"total_refunded":0,"refund_amount":0,"customer":{"id":1,"site_id":1470,"firstname":"Pamela","lastname":"Furtado","email":"daniel15@yahoo.com","telephone":"3728675986","postcode":"78448688","address_street":"Rua Verdugo","address_street_number":"3","address_street_complement":"11 Andar","address_street_district":"laboriosam","address_city":"Sao Emanuel do Sul","address_state":"TO","document_number":"82132210799","created_at":"2026-03-17 09:25:13","visited_url":"http://example.com/test-product","uf":"TO","fullname":"Pamela Furtado"},"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"267.48","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderAuthorized","event_type":"","data":{"id":1,"customer_id":1,"total_products":330,"status":"autorizado","freight_value":82.74,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:18","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":412.74,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"412.74","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderBilletCreated","event_type":"","data":{"id":1,"customer_id":1,"total_products":177,"status":"pendente","freight_value":130.15,"freight_type":"PAC","payment_type":"Boleto","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:24","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":307.15,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"307.15","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderPaid","event_type":"","data":{"id":1,"customer_id":1,"total_products":265,"status":"aprovado","freight_value":66.86,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:29","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":331.86,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"331.86","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderPendingIntegration","event_type":"","data":{"id":1,"customer_id":1,"total_products":285,"status":"pendente_integracao","freight_value":124.54,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:33","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":409.54,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"409.54","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderRefund","event_type":"","data":{"id":1,"customer_id":1,"total_products":493,"status":"estornado","freight_value":89.67,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:39","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":582.67,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"582.67","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderPartialRefund","event_type":"","data":{"id":1,"customer_id":1,"total_products":203,"status":"aprovado","freight_value":1.99,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:43","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":204.99,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"204.99","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderUpSold","event_type":"","data":{"id":1,"customer_id":1,"total_products":389,"status":"aprovado","freight_value":13.4,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:48","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":402.4,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"402.40","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderPixCreated","event_type":"","data":{"id":1,"customer_id":1,"total_products":362,"status":"pendente","freight_value":72.01,"freight_type":"PAC","payment_type":"Pix","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:25:54","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":434.01,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"434.01","pix_payment_link":"https://breakingcode.sandboxappmax.com.br/show-pix/1","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderPaidByPix","event_type":"","data":{"id":1,"customer_id":1,"total_products":129,"status":"aprovado","freight_value":83.22,"freight_type":"PAC","payment_type":"Pix","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:05","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":212.22,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"212.22","pix_payment_link":"https://breakingcode.sandboxappmax.com.br/show-pix/1","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderPixExpired","event_type":"","data":{"id":1,"customer_id":1,"total_products":335,"status":"aprovado","freight_value":133.47,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:00","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":468.47,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"468.47","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderIntegrated","event_type":"","data":{"id":1,"customer_id":1,"total_products":190,"status":"integrado","freight_value":147.42,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:10","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":337.42,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"337.42","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderBilletOverdue","event_type":"","data":{"id":1,"customer_id":1,"total_products":439,"status":"cancelado","freight_value":8.77,"freight_type":"PAC","payment_type":"Boleto","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:15","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":447.77,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"447.77","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderChargeBackInTreatment","event_type":"","data":{"id":1,"customer_id":1,"total_products":121,"status":"aprovado","freight_value":32.18,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:25","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":153.18,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"153.18","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"OrderChargeBackGain","event_type":"","data":{"id":1,"customer_id":1,"total_products":316,"status":"chargeback_vencido","freight_value":104.02,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:31","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":420.02,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"420.02","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"ChargeFailed","event_type":"","data":{"id":1,"customer_id":1,"total_products":427,"status":"aprovado","freight_value":70.2,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:46","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":497.2,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"497.20","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"ChargeSuccess","event_type":"","data":{"id":1,"customer_id":1,"total_products":213,"status":"aprovado","freight_value":41.56,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:49","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":254.56,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"254.56","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"CreatedSubscription","event_type":"","data":{"id":1,"customer_id":1,"total_products":434,"status":"aprovado","freight_value":74.8,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:26:54","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":508.8,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"508.80","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"CanceledSubscription","event_type":"","data":{"id":1,"customer_id":1,"total_products":428,"status":"aprovado","freight_value":38.4,"freight_type":"PAC","payment_type":"CreditCard","card_brand":null,"card_details":null,"pix_creation_date":null,"pix_expiration_date":null,"pix_emv":null,"pix_ref":null,"pix_qrcode":null,"pix_end_to_end_id":null,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null,"order_billet_payment_code":null,"installments":null,"paid_at":null,"refunded_at":null,"integrated_at":null,"created_at":"2026-03-17 09:27:16","discount":null,"interest":null,"upsell_order_id":null,"origin":"Site","seller_name":null,"total":466.4,"total_refunded":0,"refund_amount":0,"bundles":[],"products":[],"visit":[],"company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com","co_production_commission":null,"affiliate_comission":[],"traffic_description":null,"full_payment_amount":"466.40","pix_payment_link":"","payment_link_id":null,"invoice_id":null,"issuer_message":null}}`,
		`{"event":"SubscriptionCancellationEvent","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Noeli","lastname":"Guerra","email":"felipe22@galhardo.biz","telephone":"5191913915","postcode":"91642854","address_street":"Travessa Mendes","address_street_number":"39332","address_street_complement":"4 Andar","address_street_district":"voluptatum","address_city":"Noel do Norte","address_state":"SE","document_number":"75798407829","created_at":"2026-03-17 09:26:35","visited_url":"","uf":"SE","fullname":"Noeli Guerra","interested_bundle":[],"subscription":{"id":null,"installments":null,"total":null,"created_at":"2026-03-17 09:26:35","cancelled_at":null,"bundle":{"id":1,"name":"My product 1","description":"","production_cost":"R$ 0,00","identifier":null,"products":[{"id":1,"sku":"9000010","name":"Livro de receitas","description":"","price":"62.00","quantity":1,"image":"https://gateway-boleto-homolog-appmax.s3-sa-east-1.amazonaws.com/","external_id":null}]}}}}`,
		`{"event":"SubscriptionDelayedEvent","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Ricardo","lastname":"Faro","email":"laura90@avila.org","telephone":"3121504676","postcode":"63209123","address_street":"Avenida Andres Vale","address_street_number":"7","address_street_complement":"Apto 7","address_street_district":"soluta","address_city":"Maia do Norte","address_state":"SE","document_number":"35432870444","created_at":"2026-03-17 09:26:41","visited_url":"","uf":"SE","fullname":"Ricardo Faro","interested_bundle":[],"subscription":{"id":null,"installments":null,"total":null,"created_at":"2026-03-17 09:26:41","cancelled_at":null,"bundle":{"id":1,"name":"My product 1","description":"","production_cost":"R$ 0,00","identifier":null,"products":[{"id":1,"sku":"9000010","name":"Livro de receitas","description":"","price":"62.00","quantity":1,"image":"https://gateway-boleto-homolog-appmax.s3-sa-east-1.amazonaws.com/","external_id":null}]}}}}`,
		`{"event":"order_paid","event_type":"order","data":{"order_id":12844,"foo":"bar"}}`,
		`{"event":"OrderApproved","event_type":"","data":{"order_id":1,"order_total_products":369,"order_status":"aprovado","order_freight_value":24.26,"order_freight_type":"PAC","order_payment_type":"CreditCard","order_total":393.26,"order_total_refunded":0,"order_refund_amount":0,"order_created_at":"2026-03-17 09:28:17","customer_id":1,"customer_firstname":"Leandro","customer_lastname":"Soto","customer_email":"toledo.sergio@gmail.com","customer_uf":"AC","customer_fullname":"Leandro Soto","company_name":"DemoAppGeovanne","company_cnpj":"66535306046","company_email":"geovanne.gallinati@teste.com"}}`,
		`{"event":"OrderApproved","event_type":"","data":{"id":1,"customer_id":1,"total_products":465,"status":"aprovado","freight_value":42.22,"freight_type":"PAC","payment_type":"CreditCard","total":507.22,"full_payment_amount":"507.22","pix_payment_link":"","company_name":"DemoAppGeovanne","bundles":[],"products":[],"visit":[],"meta":[]}}`,
	}

	for _, raw := range envelopes {
		raw := raw
		var env requests.WebhookEnvelopeRequest
		err := json.Unmarshal([]byte(raw), &env)
		require.NoError(t, err, "failed to parse envelope: %s", raw[:50])
		assert.NotEmpty(t, env.Event)
		assert.NotNil(t, env.Data)
	}
}

// TestWebhookEnvelope_OrderApproved_EnvelopeFields verifies all top-level fields
// of a real Appmax order event envelope.
func TestWebhookEnvelope_OrderApproved_EnvelopeFields(t *testing.T) {
	raw := `{"event":"OrderApproved","event_type":"","data":{"id":1,"customer_id":1,"total_products":265,"status":"aprovado","freight_value":2.48,"freight_type":"PAC","payment_type":"CreditCard","total":267.48,"full_payment_amount":"267.48","pix_payment_link":"","company_name":"DemoAppGeovanne","bundles":[],"products":[],"visit":[]}}`

	var env requests.WebhookEnvelopeRequest
	err := json.Unmarshal([]byte(raw), &env)

	require.NoError(t, err)
	assert.Equal(t, "OrderApproved", env.Event)
	assert.Equal(t, "", env.EventType)
	assert.NotNil(t, env.Data)
}

// TestWebhookEnvelope_CustomerCreated_EnvelopeFields verifies all top-level fields
// of a real Appmax customer event envelope.
func TestWebhookEnvelope_CustomerCreated_EnvelopeFields(t *testing.T) {
	raw := `{"event":"CustomerCreated","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Laura","lastname":"Montenegro","email":"maximo33@gmail.com","telephone":"3526001145","postcode":"97120572","address_street":"Avenida Santiago","address_street_number":"55","address_street_complement":"Bloco C","address_street_district":"qui","address_city":"Santa Daniel","address_state":"SP","document_number":"57135108914","created_at":"2026-03-17 09:19:46","visited_url":"","uf":"SP","fullname":"Laura Montenegro","interested_bundle":[]}}`

	var env requests.WebhookEnvelopeRequest
	err := json.Unmarshal([]byte(raw), &env)

	require.NoError(t, err)
	assert.Equal(t, "CustomerCreated", env.Event)
	assert.Equal(t, "", env.EventType)
	assert.NotNil(t, env.Data)
}

// TestWebhookEnvelope_LegacyOrderPaid_HasEventType verifies that the legacy
// snake_case format includes a non-empty event_type field.
func TestWebhookEnvelope_LegacyOrderPaid_HasEventType(t *testing.T) {
	raw := `{"event":"order_paid","event_type":"order","data":{"order_id":12844,"foo":"bar"}}`

	var env requests.WebhookEnvelopeRequest
	err := json.Unmarshal([]byte(raw), &env)

	require.NoError(t, err)
	assert.Equal(t, "order_paid", env.Event)
	assert.Equal(t, "order", env.EventType)
}

// TestWebhookEnvelope_SubscriptionEvent_NestedSubscriptionObjectParsed verifies
// that the nested subscription object in the data field is preserved as raw JSON.
func TestWebhookEnvelope_SubscriptionEvent_NestedSubscriptionObjectParsed(t *testing.T) {
	raw := `{"event":"SubscriptionCancellationEvent","event_type":"","data":{"id":1,"site_id":1470,"firstname":"Noeli","lastname":"Guerra","email":"felipe22@galhardo.biz","telephone":"5191913915","postcode":"91642854","address_street":"Travessa Mendes","address_street_number":"39332","address_street_complement":"4 Andar","address_street_district":"voluptatum","address_city":"Noel do Norte","address_state":"SE","document_number":"75798407829","created_at":"2026-03-17 09:26:35","visited_url":"","uf":"SE","fullname":"Noeli Guerra","interested_bundle":[],"subscription":{"id":null,"installments":null,"total":null,"created_at":"2026-03-17 09:26:35","cancelled_at":null,"bundle":{"id":1,"name":"My product 1","description":"","production_cost":"R$ 0,00","identifier":null,"products":[{"id":1,"sku":"9000010","name":"Livro de receitas","description":"","price":"62.00","quantity":1,"image":"https://gateway-boleto-homolog-appmax.s3-sa-east-1.amazonaws.com/","external_id":null}]}}}}`

	var env requests.WebhookEnvelopeRequest
	err := json.Unmarshal([]byte(raw), &env)

	require.NoError(t, err)
	assert.Equal(t, "SubscriptionCancellationEvent", env.Event)

	var data map[string]any
	require.NoError(t, json.Unmarshal(env.Data, &data))

	subscription, ok := data["subscription"].(map[string]any)
	require.True(t, ok, "expected subscription to be a nested object")

	bundle, ok := subscription["bundle"].(map[string]any)
	require.True(t, ok, "expected bundle to be a nested object")
	assert.Equal(t, "My product 1", bundle["name"])

	products, ok := bundle["products"].([]any)
	require.True(t, ok)
	require.Len(t, products, 1)
	assert.Equal(t, "9000010", products[0].(map[string]any)["sku"])
}

// TestWebhookDataRequest_LegacyPayload_OrderIdExtracted verifies that the legacy
// Appmax format using "order_id" field is correctly parsed.
func TestWebhookDataRequest_LegacyPayload_OrderIdExtracted(t *testing.T) {
	legacyPayload := `{"order_id":12844,"foo":"bar"}`

	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(legacyPayload), &data)

	require.NoError(t, err)
	require.NotNil(t, data.OrderID.Ptr())
	assert.Equal(t, 12844, *data.OrderID.Ptr())
}

// TestExtractOrderID_PixStandardPayload_ReturnsDataId verifies that Pix event payloads
// in Standard model (data.id + data.customer_id) correctly extract the order ID from data.id.
func TestExtractOrderID_PixStandardPayload_ReturnsDataId(t *testing.T) {
	pixPayload := `{"id":1,"customer_id":1,"total_products":362,"status":"pendente","payment_type":"Pix","total":434.01,"pix_payment_link":"https://breakingcode.sandboxappmax.com.br/show-pix/1"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(pixPayload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 1, *orderID)
}

// TestExtractOrderID_BoletoStandardPayload_ReturnsDataId verifies that Boleto event payloads
// in Standard model (data.id + data.customer_id) correctly extract the order ID from data.id.
func TestExtractOrderID_BoletoStandardPayload_ReturnsDataId(t *testing.T) {
	boletoPayload := `{"id":1,"customer_id":1,"total_products":177,"status":"pendente","payment_type":"Boleto","total":307.15,"billet_date_overdue":"","billet_url":null,"billet_digitable_line":null}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(boletoPayload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 1, *orderID)
}
