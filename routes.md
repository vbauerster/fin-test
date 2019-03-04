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
			- [(*server).createAccount-fm](/app/handlers.go#L32)

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
			- _DELETE_
				- [(*server).deleteAccount-fm](/app/handlers.go#L102)
			- _GET_
				- [(*server).getAccount-fm](/app/handlers.go#L71)
			- _PUT_
				- [(*server).updateAccount-fm](/app/handlers.go#L82)

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
				- [(*server).doDeposit-fm](/app/handlers.go#L114)

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
				- [(*server).doTransfer-fm](/app/handlers.go#L223)

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
				- [(*server).doWithdraw-fm](/app/handlers.go#L165)

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
