# ica-gdpr-parse

Tool for making sense of GDPR exports from ICA.

## Usage example

```
$ ls example-input
Butik kvitto.xml      Butik kvittorader.xml
```
```
$ head example-input/Butik\ kvitto.xml
<businessObjectToFileArea><resObject>
  <TransactionHeader>
    <responseType>POS_Receipt_Items</responseType>
    <customerId>123</customerId>
    <transactions>
      <receiptType>Selfcheckout</receiptType>
      <totalDiscount>0.0</totalDiscount>
      <transactionValue>39.9</transactionValue>
      <transactionId>C12345678.12345678</transactionId>
      <transactionTimestamp>2023-11-22 10:00:00</transactionTimestamp>
```
```
$ head example-input/Butik\ kvittorader.xml
<businessObjectToFileArea><resObject>
  <LineItems>
    <responseType>POS_Line_Items</responseType>
    <customerId>123</customerId>
    <transactions>
      <quantity>1.0</quantity>
      <price>17.95</price>
      <personalOfferId>0</personalOfferId>
      <voucherValue>0.0</voucherValue>
      <discountQuantity>0.0</discountQuantity>
```
```
$ ica-gdpr-parse example-input | jq .
{
  "from": "2023-01-01T00:00:00+01:00",
  "to": "2024-01-01T00:00:00+01:00",
  "items": [
    {
      "item": "Annas Pepparkakor Kardemumma 150g",
      "total_quantity": 1,
      "total_price": 21.95,
      "total_discount_value": 0,
      "total_discounted_price": 21.95
    },
    {
      "item": "Banan.",
      "total_quantity": 0.47,
      "total_price": 11.83,
      "total_discount_value": 0,
      "total_discounted_price": 11.83
    },
    {
      "item": "Färsk mellanmjölk 1,5%",
      "total_quantity": 1,
      "total_price": 13.8,
      "total_discount_value": 0,
      "total_discounted_price": 13.8
    },
    {
      "item": "Tomat Piccolini ICA Sel 200g..",
      "total_quantity": 1,
      "total_price": 29.95,
      "total_discount_value": 0,
      "total_discounted_price": 29.95
    },
    {
      "item": "Yoghurt naturell 3%",
      "total_quantity": 3,
      "total_price": 55.849999999999994,
      "total_discount_value": 0,
      "total_discounted_price": 55.849999999999994
    }
  ]
}
```
