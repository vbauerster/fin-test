# github.com/vbauerster/fin-test

fin-test REST API.

## Routes

<details>
<summary>`/accounts/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/accounts/***
	- **/**
		- _GET_
			- [(*server).listAccounts-fm](/app/handlers.go#L15)
		- _POST_
			- [(*server).createAccount-fm](/app/handlers.go#L31)

</details>
<details>
<summary>`/accounts/*/{accountID}/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/accounts/***
	- **/{accountID}/***
		- [(*server).accountCtx-fm](/app/middleware.go#L16)
		- **/**
			- _PUT_
				- [(*server).updateAccount-fm](/app/handlers.go#L71)
			- _DELETE_
				- [(*server).deleteAccount-fm](/app/handlers.go#L91)
			- _GET_
				- [(*server).getAccount-fm](/app/handlers.go#L61)

</details>
<details>
<summary>`/accounts/*/{accountID}/*/deposit`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/accounts/***
	- **/{accountID}/***
		- [(*server).accountCtx-fm](/app/middleware.go#L16)
		- **/deposit**
			- _POST_
				- [(*server).doDeposit-fm](/app/handlers.go#L103)

</details>
<details>
<summary>`/accounts/*/{accountID}/*/transfer`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/accounts/***
	- **/{accountID}/***
		- [(*server).accountCtx-fm](/app/middleware.go#L16)
		- **/transfer**
			- _POST_
				- [(*server).doTransfer-fm](/app/handlers.go#L212)

</details>
<details>
<summary>`/accounts/*/{accountID}/*/withdraw`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/accounts/***
	- **/{accountID}/***
		- [(*server).accountCtx-fm](/app/middleware.go#L16)
		- **/withdraw**
			- _POST_
				- [(*server).doWithdraw-fm](/app/handlers.go#L154)

</details>
<details>
<summary>`/payments/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/payments/***
	- **/**
		- _GET_
			- [(*server).listPayments-fm](/app/handlers.go#L22)

</details>

Total # of routes: 6
