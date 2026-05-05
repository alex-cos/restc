package restc

// Standard HTTP header constants.
const (
	// Accept is the media types that are acceptable.
	Accept = "Accept"
	// AcceptCharset is the character sets that are acceptable.
	AcceptCharset = "Accept-Charset"
	// AcceptEncoding is the encodings that are acceptable.
	AcceptEncoding = "Accept-Encoding"
	// AcceptLanguage is the natural languages that are acceptable.
	AcceptLanguage = "Accept-Language"
	// AcceptDatetime is the acceptable datetime.
	AcceptDatetime = "Accept-Datetime"
	// Authorization is the credentials for authenticating the user agent.
	Authorization = "Authorization"
	// CacheControl is the cache control directives.
	CacheControl = "Cache-Control"
	// Connection controls whether the network connection stays open.
	Connection = "Connection"
	// ContentLength is the size of the request body in bytes.
	ContentLength = "Content-Length"
	// ContentMD5 is the MD5 checksum of the request body.
	ContentMD5 = "Content-MD5"
	// ContentType is the media type of the request body.
	ContentType = "Content-Type"
	// Cookie contains stored cookies.
	Cookie = "Cookie"
	// Date is the date and time at which the message was sent.
	Date = "Date"
	// DoNotTrack is the user's tracking preference.
	DoNotTrack = "DNT"
	// Expect indicates expectations that need to be fulfilled.
	Expect = "Expect"
	// Forwarded provides information about the original request.
	Forwarded = "Forwarded"
	// From is the email address of the user.
	From = "From"
	// Host specifies the host and port of the server.
	Host = "Host"
	// IfMatch is the entity tag.
	IfMatch = "If-Match"
	// IfModifiedSince is the modification date constraint.
	IfModifiedSince = "If-Modified-Since"
	// IfNoneMatch is the entity tag condition.
	IfNoneMatch = "If-None-Match"
	// IfRange is the byte range request.
	IfRange = "If-Range"
	// IfUnmodifiedSince is the modification date constraint.
	IfUnmodifiedSince = "If-Unmodified-Since"
	// MaxForwards limits the number of proxies or gateways.
	MaxForwards = "Max-Forwards"
	// Origin indicates the origin of the request.
	Origin = "Origin"
	// Pragma is the implementation-specific directives.
	Pragma = "Pragma"
	// ProxyAuthorization contains credentials for the proxy.
	ProxyAuthorization = "Proxy-Authorization"
	// Range indicates the byte range to be sent.
	Range = "Range"
	// Referer is the address of the previous web page.
	Referer = "Referer"
	// TE indicates the transfer encodings the client is willing to accept.
	TE = "TE"
	// Trailer indicates that the header provides additional fields.
	Trailer = "Trailer"
	// TransferEncoding is the form of encoding used to transfer the payload.
	TransferEncoding = "Transfer-Encoding"
	// UserAgent contains information about the user agent.
	UserAgent = "User-Agent"
	// Via indicates intermediate proxies or gateways.
	Via = "Via"
	// Warning contains warnings about the request.
	Warning = "Warning"

	// AcceptPatch is the patch content types the server supports.
	AcceptPatch = "Accept-Patch"
	// AcceptRanges indicates unit types the server supports.
	AcceptRanges = "Accept-Ranges"
	// Age is the time the object was in a proxy cache.
	Age = "Age"
	// Allow lists the methods allowed for the requested URL.
	Allow = "Allow"
	// AltSvc is the alternative services.
	AltSvc = "Alt-Svc"
	// ContentDisposition indicates if content is inline or attachment.
	ContentDisposition = "Content-Disposition"
	// ContentEncoding is the encoding transformations applied.
	ContentEncoding = "Content-Encoding"
	// ContentLanguage describes the language(s) of the payload.
	ContentLanguage = "Content-Language"
	// ContentLocation indicates the location of the resource.
	ContentLocation = "Content-Location"
	// ContentRange indicates where in the full payload the partial content applies.
	ContentRange = "Content-Range"
	// ETag is the identifier for a specific version of a resource.
	ETag = "ETag"
	// Expires is the date/time after which the response is stale.
	Expires = "Expires"
	// LastModified is the date/time the resource was last modified.
	LastModified = "Last-Modified"
	// Link provides references to other resources.
	Link = "Link"
	// Location is used in redirection or when a resource has been created.
	Location = "Location"
	// P3P is the platform for privacy preferences.
	P3P = "P3P"
	// ProxyAuthenticate is the authentication method for the proxy.
	ProxyAuthenticate = "Proxy-Authenticate"
	// Refresh is automatically refresh the page.
	Refresh = "Refresh"
	// RetryAfter indicates the client should wait before retrying the request.
	RetryAfter = "Retry-After"
	// Server provides information about the server software.
	Server = "Server"
	// SetCookie contains cookies to be sent back to the server.
	SetCookie = "Set-Cookie"
	// StrictTransportSecurity forces secure communication.
	StrictTransportSecurity = "Strict-Transport-Security"
	// Upgrade is the protocol to switch to.
	Upgrade = "Upgrade"
	// Vary indicates the fields that vary in the response.
	Vary = "Vary"
	// WWWAuthenticate is the authentication method for the resource.
	WWWAuthenticate = "WWW-Authenticate"

	// AccessControlAllowCredentials indicates if credentials can be exposed.
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// AccessControlAllowHeaders indicates which headers can be used.
	AccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// AccessControlAllowMethods indicates which methods can be used.
	AccessControlAllowMethods = "Access-Control-Allow-Methods"
	// AccessControlAllowOrigin indicates if the resource can be accessed.
	AccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// AccessControlExposeHeaders indicates which headers can be exposed.
	AccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// AccessControlMaxAge indicates how long the preflight can be cached.
	AccessControlMaxAge = "Access-Control-Max-Age"
	// AccessControlRequestHeaders indicates which headers can be used.
	AccessControlRequestHeaders = "Access-Control-Request-Headers"
	// AccessControlRequestMethod indicates which method can be used.
	AccessControlRequestMethod = "Access-Control-Request-Method"
	// ClearSiteData is the clearance of browser data.
	ClearSiteData = "Clear-Site-Data"
	// ContentSecurityPolicy restricts where resources can be loaded from.
	ContentSecurityPolicy = "Content-Security-Policy"
	// ContentSecurityPolicyReportOnly restricts where resources can be loaded from.
	ContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	// CrossOriginEmbedderPolicy controls embedding of cross-origin resources.
	CrossOriginEmbedderPolicy = "Cross-Origin-Embedder-Policy"
	// CrossOriginOpenerPolicy controls access to the browsing context.
	CrossOriginOpenerPolicy = "Cross-Origin-Opener-Policy"
	// CrossOriginResourcePolicy controls access to the resource.
	CrossOriginResourcePolicy = "Cross-Origin-Resource-Policy"
	// NEL is the network error logging.
	NEL = "NEL"
	// PermissionsPolicy controls which features can be used.
	PermissionsPolicy = "Permissions-Policy"
	// ReferrerPolicy indicates which referrer to send.
	ReferrerPolicy = "Referrer-Policy"
	// ReportTo is the reporting endpoints.
	ReportTo = "Report-To"
	// TimingAllowOrigin indicates which origins can see timing data.
	TimingAllowOrigin = "Timing-Allow-Origin"

	// SecCHUA is the client hints user agent architecture.
	SecCHUA = "Sec-CH-UA"
	// SecCHUAMobile is the client hints user agent mobile.
	SecCHUAMobile = "Sec-CH-UA-Mobile"
	// SecCHUAPlatform is the client hints user agent platform.
	SecCHUAPlatform = "Sec-CH-UA-Platform"
	// SecFetchDest is the fetch destination.
	SecFetchDest = "Sec-Fetch-Dest"
	// SecFetchMode is the fetch mode.
	SecFetchMode = "Sec-Fetch-Mode"
	// SecFetchSite is the fetch site.
	SecFetchSite = "Sec-Fetch-Site"
	// SecFetchUser is the fetch user preference.
	SecFetchUser = "Sec-Fetch-User"
	// SecGPC is the global privacy control flag.
	SecGPC = "Sec-GPC"
	// SaveData is the save data preference.
	SaveData = "Save-Data"
	// UpgradeInsecureRequests is the upgrade insecure requests preference.
	UpgradeInsecureRequests = "Upgrade-Insecure-Requests"

	// XAspNetVersion is the ASP.NET version.
	XAspNetVersion = "X-AspNet-Version"
	// XAspNetMvcVersion is the ASP.NET MVC version.
	XAspNetMvcVersion = "X-AspNetMvc-Version"
	// XContentSecurityPolicy is the content security policy.
	XContentSecurityPolicy = "X-Content-Security-Policy"
	// XContentTypeOptions prevents sniffing the content type.
	XContentTypeOptions = "X-Content-Type-Options"
	// XCSRFToken is the CSRF token.
	XCSRFToken = "X-CSRF-Token" //nolint:gosec
	// XCorrelationId is the correlation ID.
	XCorrelationId = "X-Correlation-Id"
	// XDNSPrefetchControl controls DNS prefetching.
	XDNSPrefetchControl = "X-DNS-Prefetch-Control"
	// XDownloadOptions indicates how to handle the download.
	XDownloadOptions = "X-Download-Options"
	// XForwardedFor is the original client IP address.
	XForwardedFor = "X-Forwarded-For"
	// XForwardedProto is the original protocol.
	XForwardedProto = "X-Forwarded-Proto"
	// XFrameOptions protects against clickjacking.
	XFrameOptions = "X-Frame-Options"
	// XHTTPMethodOverride overrides the HTTP method.
	XHTTPMethodOverride = "X-HTTP-Method-Override"
	// XPermittedCrossDomainPolicies controls cross-domain requests.
	XPermittedCrossDomainPolicies = "X-Permitted-Cross-Domain-Policies"
	// XPoweredBy indicates the technology supporting the server.
	XPoweredBy = "X-Powered-By"
	// XRealIP is the real client IP address.
	XRealIP = "X-Real-IP"
	// XRequestedWith is the XMLHttpRequest header.
	XRequestedWith = "X-Requested-With"
	// XRequestId is the request ID.
	XRequestId = "X-Request-Id"
	// XRobotsTag controls indexing by robots.
	XRobotsTag = "X-Robots-Tag"
	// XRatelimitLimit is the rate limit maximum requests.
	XRatelimitLimit = "X-Ratelimit-Limit"
	// XRatelimitRemaining is the rate limit remaining requests.
	XRatelimitRemaining = "X-Ratelimit-Remaining"
	// XRatelimitReset is the rate limit reset time.
	XRatelimitReset = "X-Ratelimit-Reset"
	// XTraceId is the trace ID.
	XTraceId = "X-Trace-Id"
	// XUACompatible indicates the user agent compatibility.
	XUACompatible = "X-UA-Compatible"
	// XXSSProtection is the XSS filter.
	XXSSProtection = "X-XSS-Protection"
)
