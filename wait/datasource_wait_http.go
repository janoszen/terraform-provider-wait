package wait

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
				//Will not work with HTTPS due to bug https://github.com/golang/go/issues/16736
				Optional: true,
			},
			"skip_certificate_verify": {
				Type:        schema.TypeBool,
				Description: "Skip TLS certificate verification",
				Optional:    true,
				Default:     false,
			},
			"expect_content": {
				Type:        schema.TypeString,
				Description: "Only pass if the returned content matches the specified string",
				Optional:    true,
			},
			"expect_status": {
				Type:        schema.TypeInt,
				Description: "Expected HTTP status code in response",
				Optional:    true,
				Default:     200,
			},
			"response_status": {
				Type:        schema.TypeInt,
				Description: "Returned HTTP status code in response",
				Computed:    true,
			},
			"response_body": {
				Type:        schema.TypeString,
				Description: "Returned HTTP response body",
				Computed:    true,
			},
		},
		Read: dataSourceHttpRead,
	}
}

func dataSourceHttpRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	var rootCAs *x509.CertPool
	if d.Get("ca_certificate_pem") != nil {
		rootCAs = x509.NewCertPool()
		rootCAs.AppendCertsFromPEM([]byte(d.Get("ca_certificate_pem").(string)))
	} else {
		var err error
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			return err
		}
	}

	expectedStatus := d.Get("expect_status")
	expectedContent := d.Get("expect_content")
	url := d.Get("url").(string)
	var responseStatus int
	var responseBody []byte

	for {
		log.Println(fmt.Sprintf("attempting to get %s", url))
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: d.Get("skip_certificate_verify").(bool),
				RootCAs:            rootCAs,
			},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(url)
		if err == nil {
			log.Println(fmt.Sprintf("http query to %s successful", url))
			var readErr error
			responseStatus = resp.StatusCode
			responseBody, readErr = ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if expectedStatus != 0 {
				if responseStatus != expectedStatus.(int) {
					log.Println(fmt.Sprintf("expected http status code %d does not match %d", expectedStatus, responseStatus))
					time.Sleep(10 * time.Second)
					continue
				}
			}
			if expectedContent != "" {
				if readErr != nil {
					log.Println(fmt.Sprintf("http body read error %s", readErr))
					time.Sleep(10 * time.Second)
					continue
				}
				if strings.TrimSpace(string(responseBody)) != strings.TrimSpace(expectedContent.(string)) {
					log.Println(fmt.Sprintf("expected response body %s does not match %s", expectedContent, string(responseBody)))
					time.Sleep(10 * time.Second)
					continue
				}
			}
			break
		} else {
			log.Println(fmt.Sprintf("http get error %s", err))
		}
		if ctx.Err() != nil {
			log.Println(fmt.Sprintf("context error %s", ctx.Err()))
			break
		}
		log.Println("retrying in 10 seconds")
		time.Sleep(10 * time.Second)
	}

	err := d.Set("response_status", responseStatus)
	if err != nil {
		return err
	}
	return d.Set("response_body", string(responseBody))
}
