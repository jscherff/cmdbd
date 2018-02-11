General Enhancements
--------------------
- [X] Investigate web frameworks for code simplicity and maintainability.
	* [Beego](https://beego.me/docs/intro/)
	* [Echo](https://echo.labstack.com/)
	* [Gin](https://gin-gonic.github.io/gin/)
	* [Revel](https://revel.github.io/)
- [X] Add code to check if generated serial number is already taken.

Logging
-------
- [ ] Add database logging for error logs and access logs.

API Version Control
-------------------

- [X] Add /v1/prefix to all routes so future changes (v2) don't break things.

Authentication and Authorization
--------------------------------

- [X] Implement API security with [JSON web tokens](http://jwt.io/). This will prevent things like vulnerability scans and testing with production URLs from exhausing the supply of serial numbers. 
- [ ] Obtain a valid 24 Hour Fitness certificate and perform authentication over TLS.

Model Enhancements
------------------

- [X] Use enhanced database package like 'sqlx' or ORM package like 'gorp' to simplify database interactions.
	* [sqlx](https://github.com/jmoiron/sqlx)
	* [gorp](https://github.com/go-gorp/gorp)
- [X] Add code to record individual changes in database (cmdb_changes vs cmdb_audits)
- [ ] Use some kind of caching solution (sqlite memcache?) for metadata.

Scaling Improvements
--------------------

- [ ] Investigate [eTags](http://en.wikipedia.org/wiki/HTTP_ETag).

Error Handling Improvements
---------------------------
- [X] Implement gorilla recovery handler or other custom approaches. See:
	* [StackOverflow Article](https://stackoverflow.com/questions/33904503/go-gorilla-panic-handler-to-respond-with-custom-status)
	* [GitHub Article](https://elithrar.github.io/article/http-handler-error-handling-revisited/) 
	* [Dave Cheney Blog](https://dave.cheney.net/2014/12/24/inspecting-errors)

API Improvements
----------------

- [ ] Align response codes with generally-accepted standards:
	* 200 (OK) — successful HTTP request.
	* 201 (CREATED) — request that resulted in an item being successfully created.
	* 204 (NO CONTENT) —successful HTTP request where nothing is being returned in the response body.
	* 400 (BAD REQUEST) — bad request syntax, excessive size, or another client error.
	* 403 (FORBIDDEN) — client does not have permission to access this resource.
	* 404 (NOT FOUND) — resource could not be found at this time — possibly deleted or does not exist yet.
	* 500 (INTERNAL SERVER ERROR) — unexpected failure if there is no more specific information available.
- [ ] Use standard responses for different HTTP verbs:
	* GET — return 200 (OK)
	* POST — return 201 (CREATED)
	* PUT — return 200 (OK)
	* DELETE — return 204 (NO CONTENT) If the operation fails, return the most specific status code possible corresponding to the problem that was encountered.
- [ ] Remove 'cmdb' from endpoint path
