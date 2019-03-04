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
		- _POST_
			- [(*server).createAccount-fm](/app/handlers.go#L40)
		- _GET_
			- [(*server).listAccounts-fm](/app/handlers.go#L30)

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
			- _GET_
				- [(*server).getAccount-fm](/app/handlers.go#L79)
			- _PUT_
				- [(*server).updateAccount-fm](/app/handlers.go#L91)
			- _DELETE_
				- [(*server).deleteAccount-fm](/app/handlers.go#L109)

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
				- [(*server).doDeposit-fm](/app/handlers.go#L121)

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
				- [(*server).doTransfer-fm](/app/handlers.go#L228)

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
				- [(*server).doWithdraw-fm](/app/handlers.go#L171)

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
			- [(*server).listPayments-fm](/app/handlers.go#L14)

</details>
<details>
<summary>`/payments/*/{paymentID}/*`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- [SetContentType.func1]()
- **/payments/***
	- **/{paymentID}/***
		- [(*server).paymentCtx-fm](/app/middleware.go#L40)
		- **/**
			- _GET_
				- [(*server).getPayment-fm](/app/handlers.go#L21)

</details>

Total # of routes: 7
