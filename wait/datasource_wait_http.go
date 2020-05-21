package wait

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func dataSourceHttp() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "URL to connect to",
				Required:    true,
			},
			"ca_certificate_pem": {
				Type:        schema.TypeString,
				Description: "CA certificate in PEM format",
				Required:    false,
			},
			"skip_certificate_verify": {
				Type:        schema.TypeBool,
				Description: "Skip TLS certificate verification",
				Required:    false,
				Default:     false,
			},
			"expect_content": {
				Type:        schema.TypeString,
				Description: "Only pass if the returned content matches the specified string",
				Required:    false,
			},
		},
		Read: dataSourceHttpRead,
	}
}

func dataSourceHttpRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	rootCAs, _ := x509.SystemCertPool()
	if d.Get("ca_certificate_pem") != nil {
		rootCAs = x509.NewCertPool()
		rootCAs.AppendCertsFromPEM([]byte(d.Get("ca_certificate_pem").(string)))
	}

	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: d.Get("skip_certificate_verify").(bool),
				RootCAs:            rootCAs,
			},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(d.Get("url").(string))
		if err == nil {
			if d.Get("expect_content") != nil {
				responseBody, readErr := ioutil.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if readErr == nil && strings.TrimSpace(string(responseBody)) == strings.TrimSpace(d.Get("expected_content").(string)) {
					break
				}
			} else {
				_ = resp.Body.Close()
				break
			}
		} else {
			log.Println(err)
		}
		if ctx.Err() != nil {
			log.Println(ctx.Err())
			break
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}
