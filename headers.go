package grest

import (
	"fmt"
)

const (
	// HeaderKeyContentType The MIME type of the body of the request (used with POST and PUT requests). -> Content-Type: application/x-www-form-urlencoded
	HeaderKeyContentType = "Content-Type"
	// HeaderKeyAuthorization Authentication credentials for HTTP authentication. -> Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
	HeaderKeyAuthorization = "Authorization"
	// HeaderKeyContentLength The length of the response body in octets (8-bit bytes) -> Content-Length: 348
	HeaderKeyContentLength = "Content-Length"
	// HeaderKeyCookie An HTTP cookie previously sent by the server with Set-Cookie (below). -> Cookie: $Version=1; Skin=new;
	HeaderKeyCookie = "Cookie"
	// HeaderKeyContentMD5 A Base64-encoded binary MD5 sum of the content of the response. -> Content-MD5: Q2hlY2sgSW50ZWdyaXR5IQ==
	HeaderKeyContentMD5 = "Content-MD5"
	// HeaderKeyContentEncoding The type of encoding used on the data. See HTTP compression. -> Content-Encoding: gzip
	HeaderKeyContentEncoding = "Content-Encoding"
	// HeaderKeySetCookie An HTTP cookie -> Set-Cookie: UserID=JohnDoe; Max-Age=3600; Version=1
	HeaderKeySetCookie = "Set-Cookie"
)

const (
	// ContentTypeJSON "application/json"
	ContentTypeJSON = "application/json"

	// ContentTypeXML "application/xml"
	ContentTypeXML = "application/xml"

	// ContentTypeText "text/plain"
	ContentTypeText = "text/plain"

	// ContentTypeHTML "text/html"
	ContentTypeHTML = "text/html"
)

// ContentType sets the "Content-Type" header
func ContentType(contentType string) WebPart {
	return SetHeader(HeaderKeyContentType, contentType)
}

// ContentType (composing) sets the "Content-Type" header
func (w WebPart) ContentType(contentType string) WebPart {
	return Compose(w, ContentType(contentType))
}

// Authorization sets the "HeaderAuthorization" header
func Authorization(authorizationType, credentials string) WebPart {
	return SetHeader(HeaderKeyAuthorization, fmt.Sprintf("%s %s", authorizationType, credentials))
}

// Authorization (composing) sets the "HeaderAuthorization" header
func (w WebPart) Authorization(authorizationType, credentials string) WebPart {
	return Compose(w, Authorization(authorizationType, credentials))
}

// SetHeader sets a http header (key-value pair) for the response
func SetHeader(key, value string) WebPart {
	return func(u WebUnit) *WebUnit {
		u.Writer.Header().Set(key, value)
		return &u
	}
}

// SetHeaders sets the given http headers (key-value pairs) for the response
func SetHeaders(headers map[string]string) WebPart {
	return func(u WebUnit) *WebUnit {
		for k, v := range headers {
			u.Writer.Header().Set(k, v)
		}
		return &u
	}
}

// SetHeaders sets the given http headers (key-value pairs) for the response
func (w WebPart) SetHeaders(headers map[string]string) WebPart {
	return Compose(w, SetHeaders(headers))
}

// SetHeader (composing) sets a http header (key-value pair) for the response
func (w WebPart) SetHeader(key, value string) WebPart {
	return Compose(w, SetHeader(key, value))
}
