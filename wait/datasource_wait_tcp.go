package wait

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net"
	"strconv"
	"time"
)

func dataSourceTcp() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Description: "Host name to connect to",
				Required:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Port number to connect to",
				Required:    true,
			},
		},
		Read: dataSourceTcpRead,
	}
}

func dataSourceTcpRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	for {
		dialer := net.Dialer{Timeout: time.Second}
		target := d.Get("host").(string) +
			":" +
			strconv.Itoa(d.Get("port").(int))
		log.Println(fmt.Sprintf("attempting to connect %s...", target))
		conn, err := dialer.DialContext(ctx, "tcp", target)
		log.Println(conn)
		if err == nil && conn != nil {
			_ = conn.Close()
			break
		} else {
			log.Println(err)
		}
		if ctx.Err() != nil {
			log.Println(ctx.Err())
			break
		}
		time.Sleep(time.Second)
	}

	return nil
}
