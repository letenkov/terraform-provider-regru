package regru

import "fmt"

const successResult = "success"

// APIResponse represents the structure of the response received from the API.
// It contains the result status, answer data, and error information if applicable.
type APIResponse struct {
	Result string `json:"result"`

	Answer *Answer `json:"answer,omitempty"`

	ErrorCode string `json:"error_code,omitempty"`
	ErrorText string `json:"error_text,omitempty"`
}

// Error implements the error interface for APIResponse.
// It returns a formatted error string that includes the result status, error code, and error text.
func (a APIResponse) Error() string {
	return fmt.Sprintf("API %s: %s: %s", a.Result, a.ErrorCode, a.ErrorText)
}

func (a APIResponse) HasError() error {
	if a.Result != successResult {
		return a
	}

	if a.Answer != nil {
		for _, domResp := range a.Answer.Domains {
			if domResp.Result != successResult {
				return domResp
			}
		}
	}

	return nil
}

type Answer struct {
	Domains []DomainResponse `json:"domains,omitempty"`
}

type DomainResponse struct {
	Result string `json:"result"`

	DName string `json:"dname"`

	ErrorCode string `json:"error_code,omitempty"`
	ErrorText string `json:"error_text,omitempty"`
}

func (d DomainResponse) Error() string {
	return fmt.Sprintf("API %s: %s: %s", d.Result, d.ErrorCode, d.ErrorText)
}

type CreateRecordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Domains           []Domain `json:"domains,omitempty"`
	SubDomain         string   `json:"subdomain,omitempty"`
	OutputContentType string   `json:"output_content_type,omitempty"`
}

type CreateARecordRequest struct {
	CreateRecordRequest
	IPAddr string `json:"ipaddr,omitempty"`
}

type CreateAAAARecordRequest struct {
	CreateRecordRequest
	IPAddr string `json:"ipaddr,omitempty"`
}

// CreateCnameRecordRequest represents a request to create a CNAME record for a domain.
// It extends the basic CreateRecordRequest with a canonical_name field that specifies the target hostname.
type CreateCnameRecordRequest struct {
	CreateRecordRequest
	CanonicalName string `json:"canonical_name,omitempty"`
}

// CreateTxtRecordRequest represents a request to create a TXT record for a domain.
// It extends the basic CreateRecordRequest with a text field that contains the record value.
type CreateTxtRecordRequest struct {
	CreateRecordRequest
	Text string `json:"text,omitempty"`
}

// CreateMxRecordRequest represents a request to create an MX (Mail Exchange) record for a domain.
// It extends the basic CreateRecordRequest with mail server information and priority settings.
type CreateMxRecordRequest struct {
	CreateRecordRequest
	MailServer string `json:"mail_server,omitempty"`
	Priority   string `json:"priority,omitempty"`
}

// DeleteRecordRequest represents a request to delete a DNS record from the specified domain.
// It requires authentication credentials, domain information, and details about the record to delete.
type DeleteRecordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Domains           []Domain `json:"domains,omitempty"`
	SubDomain         string   `json:"subdomain,omitempty"`
	Content           string   `json:"content,omitempty"`
	RecordType        string   `json:"record_type,omitempty"`
	OutputContentType string   `json:"output_content_type,omitempty"`
}

// Domain represents a domain name structure used in API requests.
// It contains the domain name that will be used for DNS operations.
type Domain struct {
	DName string `json:"dname"`
}
