# Authorization

There are various IAM solutions available for Authorization.
Two versions we will look at is:
- Permit
- OpenFGA

## OpenFGA

- Get the StoreID
```
$ openfga-cli store list
{
  "continuation_token":"",
  "stores": [
    {
      "created_at":"2024-06-29T03:56:25.857187Z",
      "id":"01J1H266C1JPHAB8X1FV1M43D7",
      "name":"test",
      "updated_at":"2024-06-29T03:56:25.857187Z"
    },
    {
      "created_at":"2024-06-29T06:14:18.695953Z",
      "id":"01J1HA2NA7YKG13G8MJAGN3TDY",
      "name":"abc",
      "updated_at":"2024-06-29T06:14:18.695953Z"
    },
    {
      "created_at":"2024-06-29T07:08:45.63616Z",
      "id":"01J1HD6BP4MRQE75CQ5HGBXBK1",
      "name":"FGA Demo",
      "updated_at":"2024-06-29T07:08:45.63616Z"
    },
    {
      "created_at":"2024-06-29T07:09:11.683243Z",
      "id":"01J1HD75438ZM343BDR1ZDC48Q",
      "name":"FGA Demo",
      "updated_at":"2024-06-29T07:09:11.683243Z"
    }
  ]
}
```

- Observe the models returned, can be written via file (or in UI); the latest model will appear in the top .. versioned
```shell

```

- Run the binary 
```
openfga-cli model list --store-id=01J1H266C1JPHAB8X1FV1M43D7
{
  "authorization_models": [
    {
      "id":"01J1H5GHX46A6S8GXQW0F0YN42",
      "created_at":"2024-06-29T04:54:31Z"
    },
    {
      "id":"01J1H2AW63J889MMRTY9QNEKZ1",
      "created_at":"2024-06-29T03:58:59Z"
    },
    {
      "id":"01J1H29YD3TENRK8B6GQBAC9TT",
      "created_at":"2024-06-29T03:58:28Z"
    }
  ]
}
```

Once the example demo is saved ... you can double check the authorization model ..
```shell

```

- Run the example authorization
```

```
- Show the basic allow / deny
- Show the more sophisticated Google Drive scenario 
