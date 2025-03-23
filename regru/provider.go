package regru

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const defaultApiEndpoint = "https://api.reg.ru/api/regru2/"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_API_USERNAME", nil),
				Description: "API username for reg.ru",
			},
			"api_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_API_PASSWORD", nil),
				Description: "API password for reg.ru",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultApiEndpoint,
				Description: "reg.ru API endpoint",
			},
			"cert_file": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_CERT_FILE", nil),
				Description: "Path to the client SSL certificate file",
			},
			"key_file": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_KEY_FILE", nil),
				Description: "Path to the client SSL key file",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"regru_dns_record": resourceRegruDNSRecord(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("api_username").(string)
	password := d.Get("api_password").(string)
	endpoint := d.Get("api_endpoint").(string)
	certFile := d.Get("cert_file").(string)
	keyFile := d.Get("key_file").(string)

	// Проверяем обязательные параметры
	if username == "" {
		return nil, errors.New("missing required parameter: api_username")
	}
	if password == "" {
		return nil, errors.New("missing required parameter: api_password")
	}
	if certFile == "" {
		return nil, errors.New("missing required parameter: cert_file")
	}
	if keyFile == "" {
		return nil, errors.New("missing required parameter: key_file")
	}

	// Устанавливаем значение по умолчанию для endpoint, если он не указан
	if endpoint == "" {
		endpoint = defaultApiEndpoint
	}

	// Создаем клиент с сертификатом
	c, err := NewClient(username, password, endpoint, certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return c, nil
}
