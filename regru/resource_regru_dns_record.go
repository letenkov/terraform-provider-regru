package regru

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceRegruDNSRecord returns a schema.Resource that defines the DNS record resource
// for the Reg.ru provider. It defines the necessary CRUD operations and schemas for the resource.
// This resource allows you to create, read, and delete DNS records in Reg.ru domains.
func resourceRegruDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegruDNSRecordCreate,
		Read:   resourceRegruDNSRecordRead,
		Delete: resourceRegruDNSRecordDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"record": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

// Create a new DNS record in Reg.ru domain
func resourceRegruDNSRecordCreate(d *schema.ResourceData, m interface{}) error {
	record_type := d.Get("type").(string)
	record_name := d.Get("name").(string)
	value := d.Get("record").(string)
	zone := d.Get("zone").(string)

	c := m.(*Client)

	// Determine action based on record type
	action := ""
	switch strings.ToUpper(record_type) {
	case "A":
		action = "add_alias"
	case "AAAA":
		action = "add_aaaa"
	case "CNAME":
		action = "add_cname"
	case "MX":
		action = "add_mx"
	case "TXT":
		action = "add_txt"
	default:
		return fmt.Errorf("invalid record type '%s'", record_type)
	}

	// Form request parameters depending on record type
	params := make(map[string]interface{})

	// Common parameters for all types
	params["domain_name"] = zone // Use domain_name instead of zone!
	params["subdomain"] = record_name
	params["output_content_type"] = "plain"

	// Specific parameters for each type
	switch strings.ToUpper(record_type) {
	case "A":
		params["ipaddr"] = value
	case "AAAA":
		params["ipaddr"] = value
	case "CNAME":
		params["canonical_name"] = value
	case "MX":
		fields := strings.Fields(value)
		if len(fields) != 2 {
			return fmt.Errorf("invalid MX record format, expected 'priority mailserver'")
		}
		params["mail_server"] = fields[1]
		params["priority"] = fields[0]
	case "TXT":
		params["text"] = value
	}

	// Add parameters in the format expected by the API
	// For some APIs, a domains structure may be required
	if action == "add_mx" || action == "add_cname" || action == "add_txt" {
		// Use domains field for some record types (format required by API)
		params["domains"] = []map[string]string{{"dname": zone}}
		delete(params, "domain_name") // Remove previously added domain_name
	}

	// Log parameters for debugging
	fmt.Printf("Sending params to %s: %+v\n", action, params)

	// Execute the request
	resp, err := c.doRequest(params, "zone", action)
	if err != nil {
		return err
	}
	if resp.HasError() != nil {
		return resp.HasError()
	}

	d.SetId(strings.Join([]string{record_name, zone}, "."))
	return nil
}

func resourceRegruDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Delete an existing DNS record from Reg.ru domain
func resourceRegruDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	record_type := d.Get("type").(string)
	record_name := d.Get("name").(string)
	value := d.Get("record").(string)
	zone := d.Get("zone").(string)

	c := m.(*Client)

	request := DeleteRecordRequest{
		Username:          c.username,
		Password:          c.password,
		Domains:           []Domain{{DName: zone}},
		SubDomain:         record_name,
		Content:           value,
		RecordType:        strings.ToUpper(record_type),
		OutputContentType: "plain",
	}

	resp, err := c.doRequest(request, "zone", "remove_record")
	if err != nil {
		return err
	}

	return resp.HasError()
}
