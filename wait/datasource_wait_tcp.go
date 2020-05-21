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
			"ip": {
				Type:        schema.TypeString,
				Description: "IP address the connection was made to.",
				Computed:    true,
			},
		},
		Read: dataSourceTcpRead,
	}
}

func dataSourceTcpRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	target := d.Get("host").(string) +
		":" +
		strconv.Itoa(d.Get("port").(int))
	var remoteIp net.Addr
	for {
		dialer := net.Dialer{Timeout: time.Second}

		log.Println(fmt.Sprintf("attempting to connect %s", target))
		conn, err := dialer.DialContext(ctx, "tcp", target)
		if err == nil && conn != nil {
			log.Println(fmt.Sprintf("connection to %s successful", target))
			remoteIp = conn.RemoteAddr()
			_ = conn.Close()
			break
		} else {
			log.Println(fmt.Sprintf("connection to %s failed, retrying", target))
			log.Println(err)
		}
		if ctx.Err() != nil {
			log.Println(ctx.Err())
			return ctx.Err()
		}
		time.Sleep(10 * time.Second)
	}
	d.SetId(remoteIp.String() + ":" + strconv.Itoa(d.Get("port").(int)))

	return d.Set("ip", remoteIp.String())
}
