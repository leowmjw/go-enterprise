# README

## Objective

Show how to make simple go-script become invincible!

The scripts does some file operations with g-script and with uv python
The uv python using httpie for some scenarios

## To run test

```
$ gotest ./...
```

Run the traditional payment collection script; show multi run will not be idempotent ..

```bash
$ ./internal/batch/scripts/traditional_payment_collection.sh

Starting batch processing of 10 OrderIDs
======================================================

Processing OrderID: 7307 (1/10)
ERROR: OrderID 7307 failed with exit code 1 in 0s
  Starting payment processing for OrderID: 7307
  Starting processing step 1...
  Step 1 failed: FAILED: Processing Step 1 for OrderID 7307
  Cleaning up resources...
  ERROR: Script terminated with exit code: 1 - Step 1 failed: FAILED: Processing Step 1 for OrderID 7307
------------------------------------------------------

Processing OrderID: 5493 (2/10)
SUCCESS: OrderID 5493 processed successfully in 4s
  Starting payment processing for OrderID: 5493
  Starting processing step 1...
  Step 1 completed successfully: Step1 5493
  Starting processing step 2...
  Step 2 completed successfully: Step2 5493
  Payment processing completed successfully for OrderID: 5493
  Cleaning up resources...
------------------------------------------------------

Processing OrderID: 7387 (3/10)
ERROR: OrderID 7387 failed with exit code 1 in 0s
  Starting payment processing for OrderID: 7387
  Starting processing step 1...
  Step 1 failed: FAILED: Processing Step 1 for OrderID 7387
  Cleaning up resources...
  ERROR: Script terminated with exit code: 1 - Step 1 failed: FAILED: Processing Step 1 for OrderID 7387
------------------------------------------------------

Processing OrderID: 2614 (4/10)
SUCCESS: OrderID 2614 processed successfully in 4s
  Starting payment processing for OrderID: 2614
  Starting processing step 1...
  Step 1 completed successfully: Step1 2614
  Starting processing step 2...
  Step 2 completed successfully: Step2 2614
  Payment processing completed successfully for OrderID: 2614
  Cleaning up resources...
------------------------------------------------------

Processing OrderID: 5999 (5/10)
ERROR: OrderID 5999 failed with exit code 1 in 0s
  Starting payment processing for OrderID: 5999
  Starting processing step 1...
  Step 1 failed: FAILED: Processing Step 1 for OrderID 5999
  Cleaning up resources...
  ERROR: Script terminated with exit code: 1 - Step 1 failed: FAILED: Processing Step 1 for OrderID 5999
------------------------------------------------------

Processing OrderID: 3078 (6/10)
SUCCESS: OrderID 3078 processed successfully in 5s
  Starting payment processing for OrderID: 3078
  Starting processing step 1...
  Step 1 completed successfully: Step1 3078
  Starting processing step 2...
  Step 2 completed successfully: Step2 3078
  Payment processing completed successfully for OrderID: 3078
  Cleaning up resources...
------------------------------------------------------

Processing OrderID: 8577 (7/10)
ERROR: OrderID 8577 failed with exit code 1 in 0s
  Starting payment processing for OrderID: 8577
  Starting processing step 1...
  Step 1 failed: FAILED: Processing Step 1 for OrderID 8577
  Cleaning up resources...
  ERROR: Script terminated with exit code: 1 - Step 1 failed: FAILED: Processing Step 1 for OrderID 8577
------------------------------------------------------

Processing OrderID: 5479 (8/10)
ERROR: OrderID 5479 failed with exit code 1 in 0s
  Starting payment processing for OrderID: 5479
  Starting processing step 1...
  Step 1 failed: FAILED: Processing Step 1 for OrderID 5479
  Cleaning up resources...
  ERROR: Script terminated with exit code: 1 - Step 1 failed: FAILED: Processing Step 1 for OrderID 5479
------------------------------------------------------

Processing OrderID: 6606 (9/10)
ERROR: OrderID 6606 failed with exit code 3 in 2s
  Starting payment processing for OrderID: 6606
  Starting processing step 1...
  Step 1 completed successfully: Step1 6606
  Starting processing step 2...
  tr: Illegal byte sequence
  Step 2 failed:  6606 ERROR!
  Cleaning up resources...
  ERROR: Script terminated with exit code: 3 - Step 2 failed:  6606 ERROR!
------------------------------------------------------

Processing OrderID: 8448 (10/10)
ERROR: OrderID 8448 failed with exit code 1 in 0s
  Starting payment processing for OrderID: 8448
  Starting processing step 1...
  Step 1 failed: FAILED: Processing Step 1 for OrderID 8448
  Cleaning up resources...
  ERROR: Script terminated with exit code: 1 - Step 1 failed: FAILED: Processing Step 1 for OrderID 8448
------------------------------------------------------

Batch Processing Summary:
======================================================
Total OrderIDs processed: 10
Successful: 3
Failed: 7
Success rate: 30%
======================================================
```