package handler

// This file contains only Swagger annotations for proxy routes.
// The gateway forwards these requests upstream — no business logic lives here.

// --- Shared response schemas ---

// ErrorResponse is the standard gateway error body.
type ErrorResponse struct {
	Error string `json:"error" example:"unauthorized"`
}

// --- Auth Service ---

// swaggerAuthRegister documents POST /api/v1/auth/register
//
//	@Summary		Register a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{email=string,password=string,name=string}	true	"Registration payload"
//	@Success		201		{object}	object{id=string,email=string}
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse	"Email already taken"
//	@Router			/api/v1/auth/register [post]
func swaggerAuthRegister() {}

// swaggerAuthLogin documents POST /api/v1/auth/login
//
//	@Summary		Login and obtain tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{email=string,password=string}	true	"Login credentials"
//	@Success		200		{object}	object{access_token=string,refresh_token=string}
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/v1/auth/login [post]
func swaggerAuthLogin() {}

// swaggerAuthRefresh documents POST /api/v1/auth/refresh
//
//	@Summary		Refresh access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{refresh_token=string}	true	"Refresh token"
//	@Success		200		{object}	object{access_token=string,refresh_token=string}
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/v1/auth/refresh [post]
func swaggerAuthRefresh() {}

// swaggerAuthLogout documents POST /api/v1/auth/logout
//
//	@Summary		Logout (revoke refresh token)
//	@Tags			auth
//	@Security		BearerAuth
//	@Success		204
//	@Failure		401	{object}	ErrorResponse
//	@Router			/api/v1/auth/logout [post]
func swaggerAuthLogout() {}

// swaggerAuthMe documents GET /api/v1/auth/me
//
//	@Summary		Get current user profile
//	@Tags			auth
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	object{id=string,email=string,role=string}
//	@Failure		401	{object}	ErrorResponse
//	@Router			/api/v1/auth/me [get]
func swaggerAuthMe() {}

// swaggerAuthDeleteMe documents DELETE /api/v1/auth/me
//
//	@Summary		Delete current user account
//	@Tags			auth
//	@Security		BearerAuth
//	@Success		204
//	@Failure		401	{object}	ErrorResponse
//	@Router			/api/v1/auth/me [delete]
func swaggerAuthDeleteMe() {}

// --- Catalog Service (Products) ---

// swaggerProductsList documents GET /api/v1/products
//
//	@Summary	List products
//	@Tags		products
//	@Produce	json
//	@Param		page		query		int		false	"Page number"		default(1)
//	@Param		limit		query		int		false	"Items per page"	default(20)
//	@Param		category_id	query		string	false	"Filter by category"
//	@Success	200			{object}	object{items=[]object,total=int}
//	@Router		/api/v1/products [get]
func swaggerProductsList() {}

// swaggerProductsGet documents GET /api/v1/products/:id
//
//	@Summary	Get product by ID
//	@Tags		products
//	@Produce	json
//	@Param		id	path		string	true	"Product ID"
//	@Success	200	{object}	object{id=string,name=string,price=number,category_id=string}
//	@Failure	404	{object}	ErrorResponse
//	@Router		/api/v1/products/{id} [get]
func swaggerProductsGet() {}

// swaggerProductsCreate documents POST /api/v1/products
//
//	@Summary	Create product
//	@Tags		products
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object{name=string,price=number,category_id=string}	true	"Product data"
//	@Success	201		{object}	object{id=string}
//	@Failure	401		{object}	ErrorResponse
//	@Router		/api/v1/products [post]
func swaggerProductsCreate() {}

// swaggerProductsUpdate documents PUT /api/v1/products/:id
//
//	@Summary	Update product
//	@Tags		products
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string										true	"Product ID"
//	@Param		body	body		object{name=string,price=number}	true	"Updated fields"
//	@Success	200		{object}	object{id=string}
//	@Failure	401		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Router		/api/v1/products/{id} [put]
func swaggerProductsUpdate() {}

// swaggerProductsDelete documents DELETE /api/v1/products/:id
//
//	@Summary	Delete product
//	@Tags		products
//	@Security	BearerAuth
//	@Param		id	path	string	true	"Product ID"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Router		/api/v1/products/{id} [delete]
func swaggerProductsDelete() {}

// --- Catalog Service (Categories) ---

// swaggerCategoriesList documents GET /api/v1/categories
//
//	@Summary	List categories
//	@Tags		categories
//	@Produce	json
//	@Success	200	{object}	object{items=[]object}
//	@Router		/api/v1/categories [get]
func swaggerCategoriesList() {}

// swaggerCategoriesGet documents GET /api/v1/categories/:id
//
//	@Summary	Get category by ID
//	@Tags		categories
//	@Produce	json
//	@Param		id	path		string	true	"Category ID"
//	@Success	200	{object}	object{id=string,name=string}
//	@Failure	404	{object}	ErrorResponse
//	@Router		/api/v1/categories/{id} [get]
func swaggerCategoriesGet() {}

// swaggerCategoriesCreate documents POST /api/v1/categories
//
//	@Summary	Create category
//	@Tags		categories
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object{name=string}	true	"Category name"
//	@Success	201		{object}	object{id=string}
//	@Failure	401		{object}	ErrorResponse
//	@Router		/api/v1/categories [post]
func swaggerCategoriesCreate() {}

// swaggerCategoriesDelete documents DELETE /api/v1/categories/:id
//
//	@Summary	Delete category
//	@Tags		categories
//	@Security	BearerAuth
//	@Param		id	path	string	true	"Category ID"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Router		/api/v1/categories/{id} [delete]
func swaggerCategoriesDelete() {}

// --- Catalog Service (Search) ---

// swaggerSearch documents GET /api/v1/search
//
//	@Summary	Full-text product search
//	@Tags		search
//	@Produce	json
//	@Param		q		query		string	true	"Search query"
//	@Param		page	query		int		false	"Page number"	default(1)
//	@Param		limit	query		int		false	"Page size"		default(20)
//	@Success	200		{object}	object{items=[]object,total=int}
//	@Router		/api/v1/search [get]
func swaggerSearch() {}

// swaggerSearchAutocomplete documents GET /api/v1/search/autocomplete
//
//	@Summary	Search autocomplete suggestions
//	@Tags		search
//	@Produce	json
//	@Param		q		query		string	true	"Partial query"
//	@Success	200		{object}	object{suggestions=[]string}
//	@Router		/api/v1/search/autocomplete [get]
func swaggerSearchAutocomplete() {}

// --- Catalog Service (Favorites) ---

// swaggerFavoritesList documents GET /api/v1/favorites
//
//	@Summary	Get current user's favorites
//	@Tags		favorites
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	object{items=[]object}
//	@Failure	401	{object}	ErrorResponse
//	@Router		/api/v1/favorites [get]
func swaggerFavoritesList() {}

// swaggerFavoritesAdd documents POST /api/v1/favorites
//
//	@Summary	Add product to favorites
//	@Tags		favorites
//	@Security	BearerAuth
//	@Accept		json
//	@Param		body	body	object{product_id=string}	true	"Product to favorite"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Router		/api/v1/favorites [post]
func swaggerFavoritesAdd() {}

// swaggerFavoritesRemove documents DELETE /api/v1/favorites
//
//	@Summary	Remove product from favorites
//	@Tags		favorites
//	@Security	BearerAuth
//	@Accept		json
//	@Param		body	body	object{product_id=string}	true	"Product to remove"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Router		/api/v1/favorites [delete]
func swaggerFavoritesRemove() {}

// --- Cart Service ---

// swaggerCartGet documents GET /api/v1/cart/:userID
//
//	@Summary	Get cart contents
//	@Tags		cart
//	@Security	BearerAuth
//	@Produce	json
//	@Param		userID	path		string	true	"User ID (must match JWT sub)"
//	@Success	200		{object}	object{items=[]object,total=number}
//	@Failure	401		{object}	ErrorResponse
//	@Failure	403		{object}	ErrorResponse	"Accessing another user's cart"
//	@Router		/api/v1/cart/{userID} [get]
func swaggerCartGet() {}

// swaggerCartClear documents DELETE /api/v1/cart/:userID
//
//	@Summary	Clear cart
//	@Tags		cart
//	@Security	BearerAuth
//	@Param		userID	path	string	true	"User ID (must match JWT sub)"
//	@Success	204
//	@Failure	401		{object}	ErrorResponse
//	@Failure	403		{object}	ErrorResponse
//	@Router		/api/v1/cart/{userID} [delete]
func swaggerCartClear() {}

// swaggerCartAddItem documents POST /api/v1/cart/:userID/items
//
//	@Summary	Add item to cart
//	@Tags		cart
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		userID	path		string								true	"User ID"
//	@Param		body	body		object{product_id=string,qty=int}	true	"Item to add"
//	@Success	201		{object}	object{item_id=string}
//	@Failure	401		{object}	ErrorResponse
//	@Failure	403		{object}	ErrorResponse
//	@Router		/api/v1/cart/{userID}/items [post]
func swaggerCartAddItem() {}

// swaggerCartRemoveItem documents DELETE /api/v1/cart/:userID/items/:productID
//
//	@Summary	Remove item from cart
//	@Tags		cart
//	@Security	BearerAuth
//	@Param		userID		path	string	true	"User ID"
//	@Param		productID	path	string	true	"Product ID"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	403	{object}	ErrorResponse
//	@Router		/api/v1/cart/{userID}/items/{productID} [delete]
func swaggerCartRemoveItem() {}

// swaggerCartUpdateItem documents PATCH /api/v1/cart/:userID/items/:productID
//
//	@Summary	Update item quantity in cart
//	@Tags		cart
//	@Security	BearerAuth
//	@Accept		json
//	@Param		userID		path	string				true	"User ID"
//	@Param		productID	path	string				true	"Product ID"
//	@Param		body		body	object{qty=int}		true	"New quantity"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	403	{object}	ErrorResponse
//	@Router		/api/v1/cart/{userID}/items/{productID} [patch]
func swaggerCartUpdateItem() {}

// --- Order Service ---

// swaggerOrdersList documents GET /api/v1/orders
//
//	@Summary	List current user's orders
//	@Tags		orders
//	@Security	BearerAuth
//	@Produce	json
//	@Param		page	query		int	false	"Page number"	default(1)
//	@Param		limit	query		int	false	"Page size"		default(20)
//	@Success	200		{object}	object{items=[]object,total=int}
//	@Failure	401		{object}	ErrorResponse
//	@Router		/api/v1/orders [get]
func swaggerOrdersList() {}

// swaggerOrdersCreate documents POST /api/v1/orders
//
//	@Summary	Place a new order
//	@Tags		orders
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object{items=[]object,address=string}	true	"Order payload"
//	@Success	201		{object}	object{id=string,status=string}
//	@Failure	401		{object}	ErrorResponse
//	@Failure	422		{object}	ErrorResponse	"Validation error"
//	@Router		/api/v1/orders [post]
func swaggerOrdersCreate() {}

// swaggerOrdersGet documents GET /api/v1/orders/:id
//
//	@Summary	Get order by ID
//	@Tags		orders
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	object{id=string,status=string,items=[]object}
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Router		/api/v1/orders/{id} [get]
func swaggerOrdersGet() {}

// swaggerOrdersUpdateStatus documents PATCH /api/v1/orders/:id/status
//
//	@Summary	Update order status
//	@Tags		orders
//	@Security	BearerAuth
//	@Accept		json
//	@Param		id		path	string					true	"Order ID"
//	@Param		body	body	object{status=string}	true	"New status"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Router		/api/v1/orders/{id}/status [patch]
func swaggerOrdersUpdateStatus() {}

// swaggerOrdersCancel documents POST /api/v1/orders/:id/cancel
//
//	@Summary	Cancel an order
//	@Tags		orders
//	@Security	BearerAuth
//	@Param		id	path	string	true	"Order ID"
//	@Success	204
//	@Failure	401	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	409	{object}	ErrorResponse	"Order cannot be cancelled"
//	@Router		/api/v1/orders/{id}/cancel [post]
func swaggerOrdersCancel() {}

// --- Gateway system endpoints ---

// swaggerHealth documents GET /health
//
//	@Summary	Liveness probe
//	@Tags		system
//	@Produce	json
//	@Success	200
//	@Router		/health [get]
func swaggerHealth() {}

// swaggerReady documents GET /ready
//
//	@Summary	Readiness probe — checks upstream connectivity
//	@Tags		system
//	@Produce	json
//	@Success	200
//	@Failure	503	{object}	object{status=string,details=object}
//	@Router		/ready [get]
func swaggerReady() {}
