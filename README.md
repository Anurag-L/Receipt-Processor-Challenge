
1. After cloning, run "go build".
2. Run "go run main.go"
   The above are the steps to get the backend up


   
3. From another power shell run this to generate id for a receipt. Replace the data values with any value you want to test the first API. This should generate a unique Id.

$jsonData = '{
>>   "retailer": "Target",
>>   "purchaseDate": "2022-01-01",
>>   "purchaseTime": "13:01",
>>   "items": [
PS C:\Users\abbas\receipt-processor> $jsonData = '{
>>   "retailer": "M&M Corner Market",
>>   "purchaseDate": "2022-03-20",
>>   "purchaseTime": "14:33",
>>   "items": [
>>     {
>>       "shortDescription": "Gatorade",
>>       "price": "2.25"
>>     },{
>>       "shortDescription": "Gatorade",
>>       "price": "2.25"
>>     },{
>>       "shortDescription": "Gatorade",
>>       "price": "2.25"
>>     },{
>>       "shortDescription": "Gatorade",
>>       "price": "2.25"
>>     }
>>   ],
>>   "total": "9.00"
>> }'
>>
>> Invoke-RestMethod -Uri "http://localhost:8080/receipts/process" -Method Post -ContentType "application/json" -Body $jsonData

4.To calculate the total score of the receipt using the unique id generated above.
  
   Invoke-RestMethod -Uri "http://localhost:8080/receipts/{id you got above}/points" -Method Get .

